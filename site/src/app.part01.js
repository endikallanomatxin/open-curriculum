const state = {
  user: { signedIn: false, email: "", fullName: "", isMaintainer: false },
  storage: { database: false, files: false },
  units: [],
  dependencies: [],
  resources: [],
  proposals: [],
  progress: [],
  search: "",
  group: "",
  target: "",
  proposalFilter: "open",
  selectedProposalId: "",
};

const el = (id) => document.getElementById(id);
const graph = el("graph");
const graphShell = el("graph-shell");
const graphLines = el("graph-lines");
const proposalList = el("proposal-list");
const proposalDetail = el("proposal-detail");
const unitDialog = el("unit-dialog");
const proposalDialog = el("proposal-dialog");
const proposalForm = el("proposal-form");

function escapeHtml(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

function renderMarkdown(markdown) {
  const escaped = escapeHtml(markdown || "");
  const lines = escaped.split(/\r?\n/);
  const html = [];
  let inList = false;
  const closeList = () => {
    if (inList) html.push("</ul>");
    inList = false;
  };
  for (const rawLine of lines) {
    const line = rawLine.trimEnd();
    if (/^###\s+/.test(line)) {
      closeList(); html.push(`<h3>${inlineMarkdown(line.replace(/^###\s+/, ""))}</h3>`);
    } else if (/^##\s+/.test(line)) {
      closeList(); html.push(`<h2>${inlineMarkdown(line.replace(/^##\s+/, ""))}</h2>`);
    } else if (/^#\s+/.test(line)) {
      closeList(); html.push(`<h1>${inlineMarkdown(line.replace(/^#\s+/, ""))}</h1>`);
    } else if (/^-\s+/.test(line)) {
      if (!inList) { html.push("<ul>"); inList = true; }
      html.push(`<li>${inlineMarkdown(line.replace(/^-\s+/, ""))}</li>`);
    } else if (!line.trim()) {
      closeList();
    } else {
      closeList(); html.push(`<p>${inlineMarkdown(line)}</p>`);
    }
  }
  closeList();
  return html.join("");
}

function inlineMarkdown(value) {
  return value
    .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
    .replace(/`(.+?)`/g, "<code>$1</code>");
}

function formatDate(value) {
  if (!value) return "";
  return new Intl.DateTimeFormat("es-ES", { dateStyle: "medium" }).format(new Date(value));
}

function humanBytes(value) {
  const bytes = Number(value || 0);
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 ** 2).toFixed(1)} MB`;
}

function toast(message) {
  const node = el("toast");
  node.textContent = message;
  node.classList.add("visible");
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => node.classList.remove("visible"), 2800);
}

async function api(path, options = {}) {
  const response = await fetch(path, {
    ...options,
    headers: options.body instanceof FormData
      ? options.headers
      : { "Content-Type": "application/json", ...(options.headers || {}) },
  });
  const data = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(data?.error?.message || `Error ${response.status}`);
  return data;
}

async function loadState({ preserveProposal = true } = {}) {
  try {
    const data = await api("/api/state");
    Object.assign(state, {
      user: data.user,
      storage: data.storage,
      units: data.units,
      dependencies: data.dependencies,
      resources: data.resources,
      proposals: data.proposals,
      progress: data.progress,
    });
    renderAll();
    if (preserveProposal && state.selectedProposalId) await selectProposal(state.selectedProposalId, false);
  } catch (error) {
    graph.innerHTML = `<div class="error-panel">${escapeHtml(error.message)}</div>`;
    proposalList.innerHTML = `<div class="error-panel">${escapeHtml(error.message)}</div>`;
  }
}

function renderAll() {
  renderIdentity();
  renderStats();
  renderFilters();
  renderGraph();
  renderProposalList();
}

function renderIdentity() {
  const identity = el("identity");
  if (!state.user.signedIn) {
    identity.innerHTML = `<a class="button secondary small" href="/signin-with-chatgpt">Iniciar sesión</a>`;
    return;
  }
  identity.innerHTML = `
    <span class="identity-copy"><strong>${escapeHtml(state.user.fullName || state.user.email)}</strong><small>${state.user.isMaintainer ? "Mantenimiento" : "Participante"}</small></span>
    <a class="button secondary small" href="/signout-with-chatgpt">Salir</a>`;
}

function renderStats() {
  const openProposals = state.proposals.filter((proposal) => proposal.status === "open").length;
  const groups = new Set(state.units.map((unit) => unit.group_name)).size;
  el("stats").innerHTML = [
    ["Unidades", state.units.length],
    ["Relaciones", state.dependencies.length],
    ["Áreas", groups],
    ["Propuestas abiertas", openProposals],
  ].map(([label, value]) => `<div><dt>${escapeHtml(label)}</dt><dd>${value}</dd></div>`).join("");
  el("storage-status").innerHTML = `
    <div class="status-line"><span>Base de datos D1</span><span class="status-pill ${state.storage.database ? "" : "off"}">${state.storage.database ? "Activa" : "No configurada"}</span></div>
    <div class="status-line"><span>Object storage R2</span><span class="status-pill ${state.storage.files ? "" : "off"}">${state.storage.files ? "Activo" : "No configurado"}</span></div>`;
}

function renderFilters() {
  const groups = [...new Set(state.units.map((unit) => unit.group_name))].sort((a, b) => a.localeCompare(b, "es"));
  el("group-filter").innerHTML = `<option value="">Todos</option>${groups.map((group) => `<option value="${escapeHtml(group)}" ${state.group === group ? "selected" : ""}>${escapeHtml(group)}</option>`).join("")}`;
  el("target-select").innerHTML = `<option value="">Sin objetivo</option>${state.units.map((unit) => `<option value="${escapeHtml(unit.id)}" ${state.target === unit.id ? "selected" : ""}>${escapeHtml(unit.title)}</option>`).join("")}`;
}

function levelsForGraph() {
  const remaining = new Set(state.units.map((unit) => unit.id));
  const levelMap = new Map();
  let level = 0;
  while (remaining.size) {
    const ready = [...remaining].filter((id) => state.dependencies
      .filter((edge) => edge.unit_id === id)
      .every((edge) => levelMap.has(edge.prerequisite_id)));
    if (!ready.length) {
      for (const id of remaining) levelMap.set(id, level);
      break;
    }
    for (const id of ready) { levelMap.set(id, level); remaining.delete(id); }
    level += 1;
  }
  return levelMap;
}

function prerequisiteClosure(targetId) {
  if (!targetId) return new Set();
  const path = new Set([targetId]);
  const queue = [targetId];
  while (queue.length) {
    const current = queue.shift();
    for (const edge of state.dependencies.filter((item) => item.unit_id === current)) {
      if (!path.has(edge.prerequisite_id)) {
        path.add(edge.prerequisite_id);
        queue.push(edge.prerequisite_id);
      }
    }
  }
  return path;
}

function progressMap() {
  return new Map(state.progress.map((item) => [item.unit_id, item.status]));
}

function renderGraph() {
  const levels = levelsForGraph();
  const path = prerequisiteClosure(state.target);
  const progress = progressMap();
  const query = state.search.trim().toLowerCase();
  const visible = state.units.filter((unit) => {
    const matchesSearch = !query || `${unit.title} ${unit.summary} ${unit.group_name}`.toLowerCase().includes(query);
    const matchesGroup = !state.group || unit.group_name === state.group;
    return matchesSearch && matchesGroup;
  });
  if (!visible.length) {
    graph.innerHTML = `<div class="empty-state"><p>No hay unidades que coincidan con los filtros.</p></div>`;
    graphLines.innerHTML = "";
    return;
  }
  const maxLevel = Math.max(...visible.map((unit) => levels.get(unit.id) || 0));
  graph.innerHTML = Array.from({ length: maxLevel + 1 }, (_, index) => {
    const units = visible.filter((unit) => (levels.get(unit.id) || 0) === index);
    return `<div class="graph-column"><p class="graph-column-title"