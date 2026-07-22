 env) {
  const email = normalizeEmail(request.headers.get("oai-authenticated-user-email"));
  const fullName = String(request.headers.get("oai-authenticated-user-full-name") || "").trim();
  const maintainers = new Set(parseCsv(env.MAINTAINER_EMAILS).map(normalizeEmail));
  return {
    signedIn: Boolean(email),
    email,
    fullName: fullName || email,
    isMaintainer: Boolean(email && maintainers.has(email)),
  };
}

function requireIdentity(request, env) {
  const user = identityFrom(request, env);
  if (!user.signedIn) return { error: fail("Inicia sesión con ChatGPT para continuar.", 401, "AUTH_REQUIRED") };
  return { user };
}

function randomId(prefix) {
  return `${prefix}-${crypto.randomUUID()}`;
}

async function ensureDatabase(env) {
  if (!env.DB) throw new Error("D1 binding DB is not configured");
  if (!databaseReady) {
    databaseReady = (async () => {
      await env.DB.exec(SCHEMA);
      const now = new Date().toISOString();
      const statements = [];
      for (const unit of SEED_UNITS) {
        statements.push(
          env.DB.prepare(
            `INSERT OR IGNORE INTO units
            (id, slug, title, summary, content_md, group_name, locale, created_by, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, 'es', 'seed', ?, ?)`,
          ).bind(unit.id, unit.slug, unit.title, unit.summary, unit.content, unit.group, now, now),
        );
      }
      for (const [unitId, prerequisiteId] of SEED_DEPENDENCIES) {
        statements.push(
          env.DB.prepare(
            `INSERT OR IGNORE INTO dependencies
            (unit_id, prerequisite_id, created_by, created_at) VALUES (?, ?, 'seed', ?)`,
          ).bind(unitId, prerequisiteId, now),
        );
      }
      statements.push(env.DB.prepare("INSERT OR IGNORE INTO app_meta (key, value) VALUES ('schema_version', '1')"));
      await env.DB.batch(statements);
    })().catch((error) => {
      databaseReady = undefined;
      throw error;
    });
  }
  return databaseReady;
}

async function readRequestData(request) {
  const contentType = request.headers.get("content-type") || "";
  if (contentType.includes("multipart/form-data")) {
    const form = await request.formData();
    return { form, data: Object.fromEntries([...form.entries()].filter(([, value]) => typeof value === "string")) };
  }
  if (contentType.includes("application/json")) {
    return { form: null, data: await request.json() };
  }
  throw new Error("UNSUPPORTED_CONTENT_TYPE");
}

function validateUnitPayload(data, knownUnitIds, targetUnitId) {
  const title = String(data.title || "").trim().slice(0, 160);
  const summary = String(data.summary || "").trim().slice(0, 700);
  const content = String(data.content || data.content_md || "").trim().slice(0, 100000);
  const groupName = String(data.groupName || data.group_name || "General").trim().slice(0, 100) || "General";
  const locale = String(data.locale || "es").trim().toLowerCase().slice(0, 8) || "es";
  const rawPrerequisites = Array.isArray(data.prerequisiteIds)
    ? data.prerequisiteIds
    : String(data.prerequisiteIds || data.prerequisite_ids || "").split(",");
  const prerequisiteIds = uniqueStrings(rawPrerequisites).filter((id) => id !== targetUnitId);
  if (!title) throw new Error("El título es obligatorio.");
  if (!summary) throw new Error("El resumen es obligatorio.");
  for (const id of prerequisiteIds) {
    if (!knownUnitIds.has(id)) throw new Error("Una de las dependencias ya no existe.");
  }
  return { title, summary, content, groupName, locale, prerequisiteIds };
}

async function sha256(bytes) {
  const digest = await crypto.subtle.digest("SHA-256", bytes);
  return [...new Uint8Array(digest)].map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

async function audit(db, actor, action, entityType, entityId, details = {}) {
  await db.prepare(
    `INSERT INTO audit_events
    (id, actor_email, action, entity_type, entity_id, details_json, created_at)
    VALUES (?, ?, ?, ?, ?, ?, ?)`,
  ).bind(randomId("audit"), actor || "system", action, entityType, entityId, JSON.stringify(details), new Date().toISOString()).run();
}

async function getState(request, env) {
  const user = identityFrom(request, env);
  const [units, dependencies, files, proposals, progress] = await Promise.all([
    env.DB.prepare(
      `SELECT id, slug, title, summary, content_md, group_name, locale, created_by, created_at, updated_at
       FROM units ORDER BY group_name, title`,
    ).all(),
    env.DB.prepare("SELECT unit_id, prerequisite_id FROM dependencies ORDER BY unit_id, prerequisite_id").all(),
    env.DB.prepare(
      `SELECT id, target_unit_id AS unit_id, title, filename, content_type, size_bytes, created_at
       FROM proposal_files WHERE approved = 1 ORDER BY created_at DESC`,
    ).all(),
    env.DB.prepare(
      `SELECT p.id, p.kind, p.target_unit_id, p.title, p.rationale, p.payload_json, p.status,
              p.author_email, p.author_name, p.created_at, p.resolved_at, p.resolved_by,
              (SELECT COUNT(*) FROM votes v WHERE v.proposal_id = p.id AND v.value = 1) AS yes_votes,
              (SELECT COUNT(*) FROM votes v WHERE v.proposal_id = p.id AND v.value = -1) AS no_votes,
              (SELECT COUNT(*) FROM votes v WHERE v.proposal_id = p.id AND v.value = 0) AS abstain_votes,
              (SELECT COUNT(*) FROM proposal_comments c WHERE c.proposal_id = p.id) AS comment_count
       FROM proposals p
       ORDER BY CASE p.status WHEN 'open' THEN 0 ELSE 1 END, p.created_at DESC`,
    ).all(),
    user.signedIn
      ? env.DB.prepare("SELECT unit_id, status, updated_at FROM user_progress WHERE user_email = ?").bind(user.email).all()
      : Promise.resolve({ results: [] }),
  ]);

  return json({
    ok: true,
    user,
    storage: { database: Boolean(env.DB), files: Boolean(env.FILES) },
    units: units.results || [],
    dependencies: dependencies.results || [],
    resources: files.results || [],
    proposals: (proposals.results || []).map((proposal) => ({
      ...proposal,
      payload: JSON.parse(proposal.payload_json || "{}"),
      payload_json: undefined,
    })),
    progress: progress.results || [],
  });
}

async function getProposal(request, env, proposalId) {
  const user = identityFrom(request, env);
  const proposal = await env.DB.prepare("SELECT * FROM proposals WHERE id = ?").bind(proposalId).first();
  if (!proposal) return fail("La propuesta no existe.", 404, "NOT_FOUND");
  const [votes, comments, files] = await Promise.all([
    env.DB.prepare(
      `SELECT
        COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) AS yes_votes,
        COALESCE(SUM(CASE WHEN value = -1 THEN 1 ELSE 0 END), 0) AS no_votes,
        COALESCE(SUM(CASE WHEN value = 0 THEN 1 ELSE 0 END), 0) AS abstain_votes
       FROM votes WHERE proposal_id = ?`,
    ).bind(proposalId).first(),
    env.DB.prepare(
      `SELECT id, author_email, author_name, body, created_at
       FROM proposal_comments WHERE proposal_id = ? ORDER BY created_at`,
    ).bind(proposalId).all(),
    env.DB.prepare(
      `SELECT id, title, filename, content_type, size_bytes, approved, created_at
       FROM proposal_files WHERE proposal_id = ? ORDER BY created_at`,
    ).bind(proposalId).all(),
  ]);
  const ownVote = user.signedIn
    ? await env.DB.prepare("SELECT value FROM votes WHERE proposal_id = ? AND voter_email = ?").bind(proposalId, user.email).first()
    : null;
  return json({
    ok: true,
    proposal: { ...proposal, payload: JSON.parse(proposal.payload_json || "{}"), payload_json: undefined },
    votes: votes || { yes_votes: 0, no_votes: 0, abstain_votes: 0 },
    ownVote: ownVote?.value ?? null,
    comments: comments.results || [],
    files: files.results || [],
    user,
  });
}

async function createProposal(request, env) {
  const auth = requireIdentity(request, env);
  if (auth.error) return auth.error;
  const { user } = auth;
  let parsed;
  try {
    parsed = await readRequestData(request);
  } catch {
    