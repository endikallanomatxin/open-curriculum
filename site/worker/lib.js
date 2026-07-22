export function normalizeEmail(value) {
  return String(value || "").trim().toLowerCase();
}

export function slugify(value) {
  return String(value || "")
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 72) || "unidad";
}

export function safeFileName(value) {
  const normalized = String(value || "archivo")
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/[^a-zA-Z0-9._-]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .slice(0, 120);
  return normalized || "archivo";
}

export function parseCsv(value) {
  return String(value || "")
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

export function uniqueStrings(values) {
  return [...new Set((values || []).map(String).map((value) => value.trim()).filter(Boolean))];
}

export function wouldCreateCycle(edges, unitId, prerequisiteId) {
  if (!unitId || !prerequisiteId || unitId === prerequisiteId) return true;
  const adjacency = new Map();
  for (const edge of edges || []) {
    const from = String(edge.unit_id ?? edge.unitId ?? "");
    const to = String(edge.prerequisite_id ?? edge.prerequisiteId ?? "");
    if (!from || !to) continue;
    if (!adjacency.has(from)) adjacency.set(from, []);
    adjacency.get(from).push(to);
  }
  if (!adjacency.has(unitId)) adjacency.set(unitId, []);
  adjacency.get(unitId).push(prerequisiteId);

  const visiting = new Set();
  const visited = new Set();
  function visit(node) {
    if (visiting.has(node)) return true;
    if (visited.has(node)) return false;
    visiting.add(node);
    for (const next of adjacency.get(node) || []) {
      if (visit(next)) return true;
    }
    visiting.delete(node);
    visited.add(node);
    return false;
  }
  return visit(unitId);
}

export function topologicalLevels(units, edges) {
  const ids = (units || []).map((unit) => String(unit.id));
  const prerequisites = new Map(ids.map((id) => [id, new Set()]));
  for (const edge of edges || []) {
    const unitId = String(edge.unit_id ?? edge.unitId ?? "");
    const prerequisiteId = String(edge.prerequisite_id ?? edge.prerequisiteId ?? "");
    if (prerequisites.has(unitId) && prerequisites.has(prerequisiteId)) {
      prerequisites.get(unitId).add(prerequisiteId);
    }
  }
  const assigned = new Map();
  let changed = true;
  while (changed) {
    changed = false;
    for (const id of ids) {
      if (assigned.has(id)) continue;
      const deps = [...prerequisites.get(id)];
      if (deps.every((dep) => assigned.has(dep))) {
        const level = deps.length ? Math.max(...deps.map((dep) => assigned.get(dep))) + 1 : 0;
        assigned.set(id, level);
        changed = true;
      }
    }
  }
  for (const id of ids) if (!assigned.has(id)) assigned.set(id, 0);
  return Object.fromEntries(assigned);
}

export function escapeDispositionFilename(value) {
  return safeFileName(value).replace(/["\\]/g, "_");
}
