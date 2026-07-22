 return fail("Envía JSON o un formulario multipart válido.", 415, "UNSUPPORTED_CONTENT_TYPE");
  }
  const { data, form } = parsed;
  const kind = data.kind === "update_unit" ? "update_unit" : "create_unit";
  const existingUnits = await env.DB.prepare("SELECT id, slug FROM units").all();
  const knownUnitIds = new Set((existingUnits.results || []).map((unit) => unit.id));
  let targetUnitId = String(data.targetUnitId || data.target_unit_id || "").trim();
  if (kind === "update_unit" && !knownUnitIds.has(targetUnitId)) {
    return fail("La unidad que quieres modificar no existe.", 404, "UNIT_NOT_FOUND");
  }
  if (kind === "create_unit") {
    targetUnitId = `unit-${slugify(data.title)}-${crypto.randomUUID().slice(0, 8)}`;
  }
  let payload;
  try {
    payload = validateUnitPayload(data, knownUnitIds, targetUnitId);
  } catch (error) {
    return fail(error.message || "Datos de unidad no válidos.");
  }
  const rationale = String(data.rationale || "").trim().slice(0, 4000);
  if (!rationale) return fail("Explica por qué mejora el currículo.");
  const proposalId = randomId("proposal");
  const now = new Date().toISOString();
  let storedObjectKey = "";
  let fileStatement = null;

  const file = form?.get("file");
  if (file instanceof File && file.size > 0) {
    if (!env.FILES) return fail("El almacenamiento R2 FILES no está configurado.", 503, "R2_NOT_CONFIGURED");
    if (file.size > MAX_UPLOAD_BYTES) return fail("El archivo supera el límite de 15 MB.", 413, "FILE_TOO_LARGE");
    if (file.type && !ALLOWED_UPLOAD_TYPES.has(file.type)) {
      return fail("Tipo de archivo no admitido en este prototipo.", 415, "FILE_TYPE_NOT_ALLOWED");
    }
    const bytes = new Uint8Array(await file.arrayBuffer());
    const digest = await sha256(bytes);
    const filename = safeFileName(file.name);
    const fileId = randomId("file");
    storedObjectKey = `proposals/${proposalId}/${fileId}/${filename}`;
    await env.FILES.put(storedObjectKey, bytes, {
      httpMetadata: { contentType: file.type || "application/octet-stream" },
      customMetadata: { proposalId, targetUnitId, originalName: file.name, uploader: user.email },
    });
    fileStatement = env.DB.prepare(
      `INSERT INTO proposal_files
       (id, proposal_id, target_unit_id, title, filename, content_type, size_bytes, sha256, r2_key, approved, created_by, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
    ).bind(
      fileId,
      proposalId,
      targetUnitId,
      String(data.fileTitle || file.name).trim().slice(0, 160),
      filename,
      file.type || "application/octet-stream",
      file.size,
      digest,
      storedObjectKey,
      user.email,
      now,
    );
  }

  const statements = [
    env.DB.prepare(
      `INSERT INTO proposals
       (id, kind, target_unit_id, title, rationale, payload_json, status, author_email, author_name, created_at)
       VALUES (?, ?, ?, ?, ?, ?, 'open', ?, ?, ?)`,
    ).bind(
      proposalId,
      kind,
      targetUnitId,
      kind === "create_unit" ? `Añadir: ${payload.title}` : `Actualizar: ${payload.title}`,
      rationale,
      JSON.stringify(payload),
      user.email,
      user.fullName,
      now,
    ),
  ];
  if (fileStatement) statements.push(fileStatement);
  try {
    await env.DB.batch(statements);
  } catch (error) {
    if (storedObjectKey && env.FILES) await env.FILES.delete(storedObjectKey);
    console.error("create proposal failed", error);
    return fail("No se pudo guardar la propuesta.", 500, "DATABASE_ERROR");
  }
  await audit(env.DB, user.email, "proposal.created", "proposal", proposalId, { kind, targetUnitId });
  return json({ ok: true, proposalId }, 201);
}

async function voteProposal(request, env, proposalId) {
  const auth = requireIdentity(request, env);
  if (auth.error) return auth.error;
  const proposal = await env.DB.prepare("SELECT status FROM proposals WHERE id = ?").bind(proposalId).first();
  if (!proposal) return fail("La propuesta no existe.", 404, "NOT_FOUND");
  if (proposal.status !== "open") return fail("La votación ya está cerrada.", 409, "POLL_CLOSED");
  const { data } = await readRequestData(request);
  const value = Number(data.value);
  if (![1, 0, -1].includes(value)) return fail("El voto debe ser sí, abstención o no.");
  const now = new Date().toISOString();
  await env.DB.prepare(
    `INSERT INTO votes (proposal_id, voter_email, value, created_at, updated_at)
     VALUES (?, ?, ?, ?, ?)
     ON CONFLICT (proposal_id, voter_email)
     DO UPDATE SET value = excluded.value, updated_at = excluded.updated_at`,
  ).bind(proposalId, auth.user.email, value, now, now).run();
  await audit(env.DB, auth.user.email, "proposal.voted", "proposal", proposalId, { value });
  return json({ ok: true });
}

async function commentProposal(request, env, proposalId) {
  const auth = requireIdentity(request, env);
  if (auth.error) return auth.error;
  const proposal = await env.DB.prepare("SELECT id FROM proposals WHERE id = ?").bind(proposalId).first();
  if (!proposal) return fail("La propuesta no existe.", 404, "NOT_FOUND");
  const { data } = await readRequestData(request);
  const body = String(data.body || "").trim().slice(0, 4000);
  if (!body) return fail("El comentario está vacío.");
  const id = randomId("comment");
  const now = new Date().toISOString();
  await env.DB.prepare(
    `INSERT INTO proposal_comments (id, proposal_id, author_email, author_name, body, created_at)
     VALUES (?, ?, ?, ?, ?, ?)`,
  ).bind(id, proposalId, auth.user.email, auth.user.fullName, body, now).run();
  await audit(env.DB, auth.user.email, "proposal.commented", "proposal", proposalId);
  return json({ ok: true, id }, 201);
}

async function resolveProposal(request, env, proposalId) {
  const auth = requireIdentity(request, env);
  if (auth.error) return auth.error;
  if (!auth.user.isMaintainer) return fail("Solo una persona mantenedora puede resolver propuestas.", 403, "FORBIDDEN");
  const { data } = await readRequestData(request);
  const resolution = data.resolution === "accepted" ? "accepted" : data.resolution === "rejected" ? "rejected" : "";
  if (!resolution) return fail("La resolución debe ser accepted o rejected.");
  const proposal = await env.DB.prepare("SELECT * FROM proposals WHERE id = ?").bind(proposalId).first();
  if (!proposal) return fail("La propuesta no existe.", 404, "NOT_FOUND");
  if (proposal.status !== "open") return fail("La propuesta ya fue resuelta.", 409, "ALREADY_RESOLVED");
  const now = new Date().toISOString();

  if (resolution === "rejected") {
    await env.DB.prepare(
      "UPDATE proposals SET status = 'rejected', resolved_at = ?, resolved_by = ? WHERE id = ? AND status = 'open'",
    ).bind(now, auth.user.email, proposalId).run();
    await audit(env.DB, auth.user.email, "proposal.rejected", "proposal", proposalId);
    return json({ ok: true });
  }

  const payload = JSON.parse(proposal.payload_json || "{}");
  const targetUnitId = proposal.target_unit_id;
  const allUnits = await env.DB.prepare("SELECT id FROM units").all();
  const knownUnitIds = new Set((allUnits.results || []).map((unit) => unit.id));
  if (proposal.kind === "create_unit" && knownUnitIds.has(targetUnitId)) {
    return fail("El identificador de la unidad ya está ocupado.", 409, "UNIT_ID_CONFLICT");
  }
  if (proposal.kind === "update_unit" && !knownUnitIds.has(targetUnitId)) {
    return fail("La unidad ya no existe.", 409, "UNIT_MISSING");
  }
  const edgesResult = await env.DB.prepare("SELECT unit_id, prerequisite_id FROM dependencies").all();
  const edges = (edgesResult.results || []).filter((edge) => edge.unit_id !== targetUnitId);
  for (const prerequisiteId of payload.prerequisiteIds || []) {
    if (!knownUnitIds.has(prerequisiteId)) return fail("Una dependencia ya no existe.", 409, "DEPENDENCY_MISSING");
    if (wouldCreateCycle(edges, targetUnitId, prerequisiteId)) {
      return fail("La propuesta introduciría un ciclo en el currículo.", 409, 