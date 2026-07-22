import {
  escapeDispositionFilename,
  normalizeEmail,
  parseCsv,
  safeFileName,
  slugify,
  uniqueStrings,
  wouldCreateCycle,
} from "./lib.js";

const MAX_UPLOAD_BYTES = 15 * 1024 * 1024;
const ALLOWED_UPLOAD_TYPES = new Set([
  "application/pdf",
  "text/plain",
  "text/markdown",
  "image/png",
  "image/jpeg",
  "image/webp",
  "image/svg+xml",
  "audio/mpeg",
  "audio/ogg",
  "video/mp4",
]);

const SCHEMA = `
PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS app_meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS units (
  id TEXT PRIMARY KEY, slug TEXT NOT NULL UNIQUE, title TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '', content_md TEXT NOT NULL DEFAULT '',
  group_name TEXT NOT NULL DEFAULT 'General', locale TEXT NOT NULL DEFAULT 'es',
  created_by TEXT NOT NULL, created_at TEXT NOT NULL, updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS dependencies (
  unit_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  prerequisite_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  created_by TEXT NOT NULL, created_at TEXT NOT NULL,
  PRIMARY KEY (unit_id, prerequisite_id), CHECK (unit_id <> prerequisite_id)
);
CREATE TABLE IF NOT EXISTS proposals (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL CHECK (kind IN ('create_unit', 'update_unit')),
  target_unit_id TEXT NOT NULL, title TEXT NOT NULL, rationale TEXT NOT NULL DEFAULT '',
  payload_json TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'accepted', 'rejected')),
  author_email TEXT NOT NULL, author_name TEXT NOT NULL,
  created_at TEXT NOT NULL, resolved_at TEXT, resolved_by TEXT
);
CREATE TABLE IF NOT EXISTS proposal_files (
  id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  target_unit_id TEXT NOT NULL, title TEXT NOT NULL, filename TEXT NOT NULL,
  content_type TEXT NOT NULL, size_bytes INTEGER NOT NULL, sha256 TEXT NOT NULL,
  r2_key TEXT NOT NULL UNIQUE, approved INTEGER NOT NULL DEFAULT 0 CHECK (approved IN (0, 1)),
  created_by TEXT NOT NULL, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS votes (
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  voter_email TEXT NOT NULL, value INTEGER NOT NULL CHECK (value IN (-1, 0, 1)),
  created_at TEXT NOT NULL, updated_at TEXT NOT NULL,
  PRIMARY KEY (proposal_id, voter_email)
);
CREATE TABLE IF NOT EXISTS proposal_comments (
  id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  author_email TEXT NOT NULL, author_name TEXT NOT NULL, body TEXT NOT NULL, created_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS user_progress (
  user_email TEXT NOT NULL, unit_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('not_started', 'in_progress', 'completed')),
  updated_at TEXT NOT NULL, PRIMARY KEY (user_email, unit_id)
);
CREATE TABLE IF NOT EXISTS audit_events (
  id TEXT PRIMARY KEY, actor_email TEXT NOT NULL, action TEXT NOT NULL,
  entity_type TEXT NOT NULL, entity_id TEXT NOT NULL, details_json TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_proposals_status_created ON proposals(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_proposal_created ON proposal_comments(proposal_id, created_at);
CREATE INDEX IF NOT EXISTS idx_files_unit_approved ON proposal_files(target_unit_id, approved);
`;

const SEED_UNITS = [
  {
    id: "unit-mathematical-thinking",
    slug: "pensamiento-matematico",
    title: "Pensamiento matemático",
    summary: "Lenguaje, lógica y estrategias para construir y comprobar argumentos.",
    content: "# Pensamiento matemático\n\nAprende a formular problemas, distinguir hipótesis de conclusiones y justificar cada paso de un razonamiento.\n\n## Resultados\n\n- Interpretar definiciones con precisión.\n- Construir contraejemplos.\n- Escribir demostraciones breves y verificables.",
    group: "Fundamentos",
  },
  {
    id: "unit-programming-foundations",
    slug: "fundamentos-programacion",
    title: "Fundamentos de programación",
    summary: "Algoritmos, datos, control de flujo y descomposición de problemas.",
    content: "# Fundamentos de programación\n\nUna introducción práctica al pensamiento algorítmico, las estructuras de datos básicas y la creación de programas comprobables.",
    group: "Computación",
  },
  {
    id: "unit-linear-algebra",
    slug: "algebra-lineal",
    title: "Álgebra lineal",
    summary: "Vectores, transformaciones lineales, matrices y espacios vectoriales.",
    content: "# Álgebra lineal\n\nEstudia las estructuras que permiten representar y transformar información multidimensional.",
    group: "Matemáticas",
  },
  {
    id: "unit-calculus",
    slug: "calculo",
    title: "Cálculo",
    summary: "Cambio, aproximación, derivadas, integrales y modelos continuos.",
    content: "# Cálculo\n\nComprende cómo describir el cambio y la acumulación mediante límites, derivadas e integrales.",
    group: "Matemáticas",
  },
  {
    id: "unit-probability",
    slug: "probabilidad",
    title: "Probabilidad",
    summary: "Modelado de incertidumbre, variables aleatorias e inferencia básica.",
    content: "# Probabilidad\n\nConstruye modelos explícitos para razonar bajo incertidumbre y evaluar evidencia.",
    group: "Matemáticas",
  },
  {
    id: "unit-machine-learning",
    slug: "aprendizaje-automatico",
    title: "Aprendizaje automático",
    summary: "Modelos que aprenden patrones a partir de datos y cómo evaluarlos.",
    content: "# Aprendizaje automático\n\nDesde la preparación de datos hasta la validación, estudia cómo entrenar modelos útiles sin confundir ajuste con conocimiento.",
    group: "Inteligencia artificial",
  },
  {
    id: "unit-responsible-ai",
    slug: "ia-responsable",
    title: "IA responsable",
    summary: "Impacto, sesgos, seguridad, transparencia y gobernanza de sistemas de IA.",
    content: "# IA responsable\n\nAnaliza riesgos técnicos y sociales y aprende a documentar límites, decisiones y mecanismos de control.",
    group: "Inteligencia artificial",
  },
];

const SEED_DEPENDENCIES = [
  ["unit-linear-algebra", "unit-mathematical-thinking"],
  ["unit-calculus", "unit-mathematical-thinking"],
  ["unit-probability", "unit-mathematical-thinking"],
  ["unit-machine-learning", "unit-programming-foundations"],
  ["unit-machine-learning", "unit-linear-algebra"],
  ["unit-machine-learning", "unit-calculus"],
  ["unit-machine-learning", "unit-probability"],
  ["unit-responsible-ai", "unit-machine-learning"],
];

let databaseReady;

function securityHeaders(headers = new Headers()) {
  headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  headers.set("X-Content-Type-Options", "nosniff");
  headers.set("X-Frame-Options", "SAMEORIGIN");
  headers.set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()");
  headers.set("Cross-Origin-Opener-Policy", "same-origin");
  headers.set(
    "Content-Security-Policy",
    "default-src 'self'; img-src 'self' data:; media-src 'self'; style-src 'self'; script-src 'self'; connect-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'self'; form-action 'self'",
  );
  return headers;
}

function withSecurityHeaders(response) {
  const headers = securityHeaders(new Headers(response.headers));
  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers,
  });
}

function json(data, status = 200, extraHeaders = {}) {
  const headers = securityHeaders(new Headers(extraHeaders));
  headers.set("Content-Type", "application/json; charset=utf-8");
  headers.set("Cache-Control", "no-store");
  return new Response(JSON.stringify(data), { status, headers });
}

function fail(message, status = 400, code = "BAD_REQUEST") {
  return json({ ok: false, error: { code, message } }, status);
}

function identityFrom(request, 