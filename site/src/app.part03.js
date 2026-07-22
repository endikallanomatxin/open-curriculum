rce-list">${files.map((file) => `<a class="resource-link" href="/files/${encodeURIComponent(file.id)}" target="_blank" rel="noopener"><span>${escapeHtml(file.title || file.filename)}</span><small>${humanBytes(file.size_bytes)}</small></a>`).join("")}</div>` : ""}
    <h4>Votación</h4>
    <div class="vote-bar">
      <button class="vote-button ${ownVote === 1 ? "selected" : ""}" data-vote="1" ${proposal.status !== "open" ? "disabled" : ""}>Sí · ${votes.yes_votes}</button>
      <button class="vote-button ${ownVote === 0 ? "selected" : ""}" data-vote="0" ${proposal.status !== "open" ? "disabled" : ""}>Abstención · ${votes.abstain_votes}</button>
      <button class="vote-button ${ownVote === -1 ? "selected" : ""}" data-vote="-1" ${proposal.status !== "open" ? "disabled" : ""}>No · ${votes.no_votes}</button>
    </div>
    ${!user.signedIn && proposal.status === "open" ? `<p><a href="/signin-with-chatgpt">Inicia sesión con ChatGPT</a> para votar y comentar.</p>` : ""}
    <h4>Discusión</h4>
    <div class="comments">${comments.length ? comments.map((comment) => `<article class="comment"><header><strong>${escapeHtml(comment.author_name || comment.author_email)}</strong><time>${formatDate(comment.created_at)}</time></header><p>${escapeHtml(comment.body)}</p></article>`).join("") : `<p class="empty-state">Aún no hay comentarios.</p>`}</div>
    ${user.signedIn ? `<form class="comment-form" id="comment-form"><textarea name="body" rows="3" maxlength="4000" placeholder="Añade contexto, una objeción o una mejora concreta" required></textarea><button class="button secondary small" type="submit">Comentar</button></form>` : ""}
    ${user.isMaintainer && proposal.status === "open" ? `<div class="maintainer-actions"><button class="button primary small" data-resolution="accepted">Aceptar y aplicar</button><button class="button danger small" data-resolution="rejected">Rechazar</button></div>` : ""}`;

  proposalDetail.querySelectorAll("[data-vote]").forEach((button) => button.addEventListener("click", async () => {
    if (!user.signedIn) { location.href = "/signin-with-chatgpt"; return; }
    try {
      await api(`/api/proposals/${encodeURIComponent(proposal.id)}/vote`, { method: "POST", body: JSON.stringify({ value: Number(button.dataset.vote) }) });
      await loadState();
      await selectProposal(proposal.id, false);
    } catch (error) { toast(error.message); }
  }));
  el("comment-form")?.addEventListener("submit", async (event) => {
    event.preventDefault();
    const body = new FormData(event.currentTarget).get("body");
    try {
      await api(`/api/proposals/${encodeURIComponent(proposal.id)}/comments`, { method: "POST", body: JSON.stringify({ body }) });
      await loadState();
      await selectProposal(proposal.id, false);
    } catch (error) { toast(error.message); }
  });
  proposalDetail.querySelectorAll("[data-resolution]").forEach((button) => button.addEventListener("click", async () => {
    const resolution = button.dataset.resolution;
    const verb = resolution === "accepted" ? "aceptar y aplicar" : "rechazar";
    if (!confirm(`¿Confirmas que quieres ${verb} esta propuesta?`)) return;
    try {
      await api(`/api/proposals/${encodeURIComponent(proposal.id)}/resolve`, { method: "POST", body: JSON.stringify({ resolution }) });
      await loadState();
      await selectProposal(proposal.id, false);
      toast(resolution === "accepted" ? "Propuesta aplicada" : "Propuesta rechazada");
    } catch (error) { toast(error.message); }
  }));
}

function openProposalForm(unit = null) {
  if (!state.user.signedIn) {
    location.href = "/signin-with-chatgpt";
    return;
  }
  proposalForm.reset();
  const kind = unit ? "update_unit" : "create_unit";
  el("proposal-kind").value = kind;
  el("proposal-target-id").value = unit?.id || "";
  el("proposal-form-title").textContent = unit ? `Proponer cambios en ${unit.title}` : "Proponer una nueva unidad";
  el("proposal-title").value = unit?.title || "";
  el("proposal-group").value = unit?.group_name || "General";
  el("proposal-locale").value = unit?.locale || "es";
  el("proposal-summary").value = unit?.summary || "";
  el("proposal-content").value = unit?.content_md || "";
  el("proposal-rationale").value = "";
  el("proposal-form-message").textContent = "";
  const selected = new Set(unit ? state.dependencies.filter((edge) => edge.unit_id === unit.id).map((edge) => edge.prerequisite_id) : []);
  el("prerequisite-options").innerHTML = state.units
    .filter((candidate) => candidate.id !== unit?.id)
    .map((candidate) => `<label class="prerequisite-option"><input type="checkbox" value="${escapeHtml(candidate.id)}" ${selected.has(candidate.id) ? "checked" : ""}/><span><strong>${escapeHtml(candidate.title)}</strong><br><small>${escapeHtml(candidate.group_name)}</small></span></label>`).join("");
  updatePrerequisiteInput();
  el("prerequisite-options").querySelectorAll("input").forEach((input) => input.addEventListener("change", updatePrerequisiteInput));
  proposalDialog.showModal();
}

function updatePrerequisiteInput() {
  el("proposal-prerequisites").value = [...el("prerequisite-options").querySelectorAll("input:checked")].map((input) => input.value).join(",");
}

proposalForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  updatePrerequisiteInput();
  const message = el("proposal-form-message");
  const submit = proposalForm.querySelector("button[type=submit]");
  message.textContent = "";
  submit.disabled = true;
  submit.textContent = "Guardando…";
  try {
    const data = await api("/api/proposals", { method: "POST", body: new FormData(proposalForm) });
    proposalDialog.close();
    state.selectedProposalId = data.proposalId;
    state.proposalFilter = "open";
    updateProposalFilterButtons();
    await loadState({ preserveProposal: false });
    await selectProposal(data.proposalId);
    location.hash = "proposals";
    toast("Propuesta creada");
  } catch (error) {
    message.textContent = error.message;
  } finally {
    submit.disabled = false;
    submit.textContent = "Enviar propuesta";
  }
});

function updateProposalFilterButtons() {
  document.querySelectorAll("[data-proposal-filter]").forEach((button) => button.classList.toggle("selected", button.dataset.proposalFilter === state.proposalFilter));
}

el("search-input").addEventListener("input", (event) => { state.search = event.target.value; renderGraph(); });
el("group-filter").addEventListener("change", (event) => { state.group = event.target.value; renderGraph(); });
el("target-select").addEventListener("change", (event) => { state.target = event.target.value; renderGraph(); });
el("new-proposal-button").addEventListener("click", () => openProposalForm());
el("new-proposal-button-secondary").addEventListener("click", () => openProposalForm());
el("close-proposal-dialog").addEventListener("click", () => proposalDialog.close());
el("cancel-proposal").addEventListener("click", () => proposalDialog.close());
document.querySelectorAll("[data-proposal-filter]").forEach((button) => button.addEventListener("click", () => {
  state.proposalFilter = button.dataset.proposalFilter;
  updateProposalFilterButtons();
  renderProposalList();
}));

window.addEventListener("resize", () => requestAnimationFrame(drawGraphLines));
const observer = new IntersectionObserver((entries) => {
  const visible = entries.filter((entry) => entry.isIntersecting).sort((a, b) => b.intersectionRatio - a.intersectionRatio)[0];
  if (!visible) return;
  document.querySelectorAll("[data-nav]").forEach((link) => link.classList.toggle("active", link.dataset.nav === visible.target.id));
}, { threshold: [0.15, 0.45] });
document.querySelectorAll("main > section[id]").forEach((section) => observer.observe(section));

loadState();
