>Nivel ${index + 1}</p>${units.map((unit) => {
      const onPath = path.has(unit.id);
      const dimmed = state.target && !onPath;
      const status = progress.get(unit.id) || "";
      return `<button class="unit-card ${onPath ? "path" : ""} ${dimmed ? "dimmed" : ""} ${status}" data-unit-id="${escapeHtml(unit.id)}">
        <span class="group-label">${escapeHtml(unit.group_name)}</span>
        <h3>${escapeHtml(unit.title)}</h3>
        <p>${escapeHtml(unit.summary)}</p>
      </button>`;
    }).join("")}</div>`;
  }).join("");
  graph.querySelectorAll("[data-unit-id]").forEach((button) => button.addEventListener("click", () => openUnit(button.dataset.unitId)));
  requestAnimationFrame(drawGraphLines);
}

function drawGraphLines() {
  const shellRect = graphShell.getBoundingClientRect();
  const width = graphShell.scrollWidth;
  const height = graphShell.scrollHeight;
  graphLines.setAttribute("viewBox", `0 0 ${width} ${height}`);
  graphLines.style.width = `${width}px`;
  graphLines.style.height = `${height}px`;
  const pathSet = prerequisiteClosure(state.target);
  const paths = [];
  for (const edge of state.dependencies) {
    const from = graph.querySelector(`[data-unit-id="${CSS.escape(edge.prerequisite_id)}"]`);
    const to = graph.querySelector(`[data-unit-id="${CSS.escape(edge.unit_id)}"]`);
    if (!from || !to) continue;
    const a = from.getBoundingClientRect();
    const b = to.getBoundingClientRect();
    const x1 = a.right - shellRect.left + graphShell.scrollLeft;
    const y1 = a.top + a.height / 2 - shellRect.top + graphShell.scrollTop;
    const x2 = b.left - shellRect.left + graphShell.scrollLeft;
    const y2 = b.top + b.height / 2 - shellRect.top + graphShell.scrollTop;
    const bend = Math.max(34, (x2 - x1) * 0.45);
    const highlighted = state.target && pathSet.has(edge.unit_id) && pathSet.has(edge.prerequisite_id);
    paths.push(`<path class="${highlighted ? "path" : ""}" d="M ${x1} ${y1} C ${x1 + bend} ${y1}, ${x2 - bend} ${y2}, ${x2} ${y2}" />`);
  }
  graphLines.innerHTML = paths.join("");
}

function unitById(id) {
  return state.units.find((unit) => unit.id === id);
}

function openUnit(id) {
  const unit = unitById(id);
  if (!unit) return;
  const prerequisites = state.dependencies
    .filter((edge) => edge.unit_id === id)
    .map((edge) => unitById(edge.prerequisite_id))
    .filter(Boolean);
  const dependents = state.dependencies
    .filter((edge) => edge.prerequisite_id === id)
    .map((edge) => unitById(edge.unit_id))
    .filter(Boolean);
  const resources = state.resources.filter((resource) => resource.unit_id === id);
  const status = progressMap().get(id) || "not_started";
  el("unit-dialog-content").innerHTML = `
    <p class="eyebrow">${escapeHtml(unit.group_name)}</p>
    <h2>${escapeHtml(unit.title)}</h2>
    <p class="unit-summary">${escapeHtml(unit.summary)}</p>
    <div class="change-preview"><dl>
      <dt>Prerrequisitos</dt><dd>${prerequisites.length ? prerequisites.map((item) => escapeHtml(item.title)).join(", ") : "Ninguno"}</dd>
      <dt>Abre camino a</dt><dd>${dependents.length ? dependents.map((item) => escapeHtml(item.title)).join(", ") : "—"}</dd>
      <dt>Actualizada</dt><dd>${formatDate(unit.updated_at)}</dd>
    </dl></div>
    <div class="unit-content">${renderMarkdown(unit.content_md)}</div>
    ${resources.length ? `<h3>Recursos</h3><div class="resource-list">${resources.map((resource) => `<a class="resource-link" href="/files/${encodeURIComponent(resource.id)}" target="_blank" rel="noopener"><span>${escapeHtml(resource.title || resource.filename)}</span><small>${humanBytes(resource.size_bytes)}</small></a>`).join("")}</div>` : ""}
    <div class="unit-actions">
      ${state.user.signedIn ? `
        <button class="button small ${status === "in_progress" ? "primary" : "secondary"}" data-progress="in_progress">En progreso</button>
        <button class="button small ${status === "completed" ? "primary" : "secondary"}" data-progress="completed">Completada</button>
        <button class="button secondary small" id="edit-unit-button">Proponer edición</button>` : `<a class="button secondary small" href="/signin-with-chatgpt">Inicia sesión para guardar progreso</a>`}
    </div>`;
  el("unit-dialog-content").querySelectorAll("[data-progress]").forEach((button) => button.addEventListener("click", async () => {
    try {
      await api(`/api/progress/${encodeURIComponent(id)}`, { method: "PUT", body: JSON.stringify({ status: button.dataset.progress }) });
      await loadState();
      openUnit(id);
      toast("Progreso actualizado");
    } catch (error) { toast(error.message); }
  }));
  el("edit-unit-button")?.addEventListener("click", () => {
    unitDialog.close();
    openProposalForm(unit);
  });
  unitDialog.showModal();
}

function filteredProposals() {
  if (state.proposalFilter === "all") return state.proposals;
  return state.proposals.filter((proposal) => proposal.status === state.proposalFilter);
}

function statusLabel(status) {
  return ({ open: "Abierta", accepted: "Aceptada", rejected: "Rechazada" })[status] || status;
}

function renderProposalList() {
  const proposals = filteredProposals();
  if (!proposals.length) {
    proposalList.innerHTML = `<div class="empty-state"><p>No hay propuestas en esta vista.</p></div>`;
    return;
  }
  proposalList.innerHTML = proposals.map((proposal) => `
    <button class="proposal-card ${state.selectedProposalId === proposal.id ? "selected" : ""}" data-proposal-id="${escapeHtml(proposal.id)}">
      <header><h3>${escapeHtml(proposal.title)}</h3><span class="status-badge ${proposal.status}">${statusLabel(proposal.status)}</span></header>
      <p>${escapeHtml(proposal.rationale)}</p>
      <div class="proposal-meta"><span>↑ ${proposal.yes_votes}</span><span>↓ ${proposal.no_votes}</span><span>${proposal.comment_count} comentarios</span><span>${formatDate(proposal.created_at)}</span></div>
    </button>`).join("");
  proposalList.querySelectorAll("[data-proposal-id]").forEach((button) => button.addEventListener("click", () => selectProposal(button.dataset.proposalId)));
}

async function selectProposal(id, scroll = true) {
  state.selectedProposalId = id;
  renderProposalList();
  proposalDetail.className = "proposal-detail";
  proposalDetail.innerHTML = `<div class="loading">Cargando propuesta…</div>`;
  try {
    const data = await api(`/api/proposals/${encodeURIComponent(id)}`);
    renderProposalDetail(data);
    if (scroll && window.innerWidth < 980) proposalDetail.scrollIntoView({ behavior: "smooth", block: "start" });
  } catch (error) {
    proposalDetail.innerHTML = `<div class="error-panel">${escapeHtml(error.message)}</div>`;
  }
}

function renderProposalDetail(data) {
  const { proposal, votes, ownVote, comments, files, user } = data;
  const payload = proposal.payload || {};
  const prerequisiteNames = (payload.prerequisiteIds || []).map((id) => unitById(id)?.title || id);
  proposalDetail.innerHTML = `
    <p class="eyebrow">${proposal.kind === "create_unit" ? "Nueva unidad" : "Edición de unidad"}</p>
    <h3>${escapeHtml(proposal.title)}</h3>
    <p>${escapeHtml(proposal.rationale)}</p>
    <div class="proposal-meta"><span>Por ${escapeHtml(proposal.author_name || proposal.author_email)}</span><span>${formatDate(proposal.created_at)}</span><span class="status-badge ${proposal.status}">${statusLabel(proposal.status)}</span></div>
    <div class="change-preview"><dl>
      <dt>Título</dt><dd>${escapeHtml(payload.title)}</dd>
      <dt>Grupo</dt><dd>${escapeHtml(payload.groupName)}</dd>
      <dt>Resumen</dt><dd>${escapeHtml(payload.summary)}</dd>
      <dt>Prerrequisitos</dt><dd>${prerequisiteNames.length ? prerequisiteNames.map(escapeHtml).join(", ") : "Ninguno"}</dd>
    </dl>${payload.content ? `<details><summary>Ver contenido propuesto</summary><div class="unit-content">${renderMarkdown(payload.content)}</div></details>` : ""}</div>
    ${files.length ? `<h4>Recursos adjuntos</h4><div class="resou