 "CYCLE_DETECTED");
    }
    edges.push({ unit_id: targetUnitId, prerequisite_id: prerequisiteId });
  }

  const statements = [];
  if (proposal.kind === "create_unit") {
    statements.push(
      env.DB.prepare(
        `INSERT INTO units
         (id, slug, title, summary, content_md, group_name, locale, created_by, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      ).bind(
        targetUnitId,
        `${slugify(payload.title)}-${targetUnitId.slice(-6)}`,
        payload.title,
        payload.summary,
        payload.content,
        payload.groupName,
        payload.locale,
        proposal.author_email,
        now,
        now,
      ),
    );
  } else {
    statements.push(
      env.DB.prepare(
        `UPDATE units SET title = ?, summary = ?, content_md = ?, group_name = ?, locale = ?, updated_at = ?
         WHERE id = ?`,
      ).bind(payload.title, payload.summary, payload.content, payload.groupName, payload.locale, now, targetUnitId),
    );
    statements.push(env.DB.prepare("DELETE FROM dependencies WHERE unit_id = ?").bind(targetUnitId));
  }
  for (const prerequisiteId of payload.prerequisiteIds || []) {
    statements.push(
      env.DB.prepare(
        `INSERT INTO dependencies (unit_id, prerequisite_id, created_by, created_at) VALUES (?, ?, ?, ?)`,
      ).bind(targetUnitId, prerequisiteId, proposal.author_email, now),
    );
  }
  statements.push(
    env.DB.prepare("UPDATE proposal_files SET approved = 1 WHERE proposal_id = ?").bind(proposalId),
    env.DB.prepare(
      "UPDATE proposals SET status = 'accepted', resolved_at = ?, resolved_by = ? WHERE id = ? AND status = 'open'",
    ).bind(now, auth.user.email, proposalId),
  );
  try {
    await env.DB.batch(statements);
  } catch (error) {
    console.error("resolve proposal failed", error);
    return fail("No se pudo aplicar la propuesta de forma atómica.", 500, "APPLY_FAILED");
  }
  await audit(env.DB, auth.user.email, "proposal.accepted", "proposal", proposalId, { targetUnitId });
  return json({ ok: true, targetUnitId });
}

async function updateProgress(request, env, unitId) {
  const auth = requireIdentity(request, env);
  if (auth.error) return auth.error;
  const unit = await env.DB.prepare("SELECT id FROM units WHERE id = ?").bind(unitId).first();
  if (!unit) return fail("La unidad no existe.", 404, "NOT_FOUND");
  const { data } = await readRequestData(request);
  const status = String(data.status || "");
  if (!["not_started", "in_progress", "completed"].includes(status)) return fail("Estado de progreso no válido.");
  const now = new Date().toISOString();
  if (status === "not_started") {
    await env.DB.prepare("DELETE FROM user_progress WHERE user_email = ? AND unit_id = ?").bind(auth.user.email, unitId).run();
  } else {
    await env.DB.prepare(
      `INSERT INTO user_progress (user_email, unit_id, status, updated_at)
       VALUES (?, ?, ?, ?)
       ON CONFLICT (user_email, unit_id)
       DO UPDATE SET status = excluded.status, updated_at = excluded.updated_at`,
    ).bind(auth.user.email, unitId, status, now).run();
  }
  return json({ ok: true });
}

async function serveFile(request, env, fileId) {
  const metadata = await env.DB.prepare(
    `SELECT pf.*, p.author_email
     FROM proposal_files pf JOIN proposals p ON p.id = pf.proposal_id WHERE pf.id = ?`,
  ).bind(fileId).first();
  if (!metadata) return fail("El archivo no existe.", 404, "NOT_FOUND");
  const user = identityFrom(request, env);
  if (!metadata.approved && !(user.signedIn && (user.email === metadata.author_email || user.isMaintainer))) {
    return fail("No tienes acceso a este archivo.", 403, "FORBIDDEN");
  }
  if (!env.FILES) return fail("R2 FILES no está configurado.", 503, "R2_NOT_CONFIGURED");
  const object = await env.FILES.get(metadata.r2_key);
  if (!object) return fail("El objeto no está disponible.", 404, "OBJECT_MISSING");
  const headers = securityHeaders(new Headers());
  object.writeHttpMetadata?.(headers);
  headers.set("Content-Type", metadata.content_type || headers.get("Content-Type") || "application/octet-stream");
  headers.set("Content-Length", String(metadata.size_bytes));
  headers.set("Content-Disposition", `inline; filename="${escapeDispositionFilename(metadata.filename)}"`);
  headers.set("Cache-Control", metadata.approved ? "public, max-age=3600" : "private, no-store");
  return new Response(object.body, { headers });
}

async function handleApi(request, env, url) {
  await ensureDatabase(env);
  if (request.method === "GET" && url.pathname === "/api/health") {
    return json({ ok: true, database: Boolean(env.DB), files: Boolean(env.FILES), schemaVersion: 1 });
  }
  if (request.method === "GET" && url.pathname === "/api/state") return getState(request, env);
  if (request.method === "POST" && url.pathname === "/api/proposals") return createProposal(request, env);

  let match = url.pathname.match(/^\/api\/proposals\/([^/]+)$/);
  if (match && request.method === "GET") return getProposal(request, env, decodeURIComponent(match[1]));
  match = url.pathname.match(/^\/api\/proposals\/([^/]+)\/vote$/);
  if (match && request.method === "POST") return voteProposal(request, env, decodeURIComponent(match[1]));
  match = url.pathname.match(/^\/api\/proposals\/([^/]+)\/comments$/);
  if (match && request.method === "POST") return commentProposal(request, env, decodeURIComponent(match[1]));
  match = url.pathname.match(/^\/api\/proposals\/([^/]+)\/resolve$/);
  if (match && request.method === "POST") return resolveProposal(request, env, decodeURIComponent(match[1]));
  match = url.pathname.match(/^\/api\/progress\/([^/]+)$/);
  if (match && request.method === "PUT") return updateProgress(request, env, decodeURIComponent(match[1]));
  match = url.pathname.match(/^\/files\/([^/]+)$/);
  if (match && request.method === "GET") return serveFile(request, env, decodeURIComponent(match[1]));
  return fail("Ruta no encontrada.", 404, "NOT_FOUND");
}

const worker = {
  async fetch(request, env) {
    const url = new URL(request.url);
    try {
      if (url.pathname.startsWith("/api/") || url.pathname.startsWith("/files/")) {
        return await handleApi(request, env, url);
      }
      const response = await env.ASSETS.fetch(request);
      if (response.status !== 404) return withSecurityHeaders(response);
      if (!["GET", "HEAD"].includes(request.method) || !/\btext\/html\b/i.test(request.headers.get("accept") || "")) {
        return withSecurityHeaders(response);
      }
      const indexRequest = new Request(new URL("/index.html", request.url), {
        method: request.method,
        headers: request.headers,
      });
      return withSecurityHeaders(await env.ASSETS.fetch(indexRequest));
    } catch (error) {
      console.error("unhandled request error", error);
      return fail("Se ha producido un error interno.", 500, "INTERNAL_ERROR");
    }
  },
};

export default worker;
