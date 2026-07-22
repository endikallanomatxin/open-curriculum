PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS app_meta (
  key TEXT PRIMARY KEY,
  value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS units (
  id TEXT PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  content_md TEXT NOT NULL DEFAULT '',
  group_name TEXT NOT NULL DEFAULT 'General',
  locale TEXT NOT NULL DEFAULT 'es',
  created_by TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS dependencies (
  unit_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  prerequisite_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  created_by TEXT NOT NULL,
  created_at TEXT NOT NULL,
  PRIMARY KEY (unit_id, prerequisite_id),
  CHECK (unit_id <> prerequisite_id)
);

CREATE TABLE IF NOT EXISTS proposals (
  id TEXT PRIMARY KEY,
  kind TEXT NOT NULL CHECK (kind IN ('create_unit', 'update_unit')),
  target_unit_id TEXT NOT NULL,
  title TEXT NOT NULL,
  rationale TEXT NOT NULL DEFAULT '',
  payload_json TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'accepted', 'rejected')),
  author_email TEXT NOT NULL,
  author_name TEXT NOT NULL,
  created_at TEXT NOT NULL,
  resolved_at TEXT,
  resolved_by TEXT
);

CREATE TABLE IF NOT EXISTS proposal_files (
  id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  target_unit_id TEXT NOT NULL,
  title TEXT NOT NULL,
  filename TEXT NOT NULL,
  content_type TEXT NOT NULL,
  size_bytes INTEGER NOT NULL,
  sha256 TEXT NOT NULL,
  r2_key TEXT NOT NULL UNIQUE,
  approved INTEGER NOT NULL DEFAULT 0 CHECK (approved IN (0, 1)),
  created_by TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS votes (
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  voter_email TEXT NOT NULL,
  value INTEGER NOT NULL CHECK (value IN (-1, 0, 1)),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  PRIMARY KEY (proposal_id, voter_email)
);

CREATE TABLE IF NOT EXISTS proposal_comments (
  id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
  author_email TEXT NOT NULL,
  author_name TEXT NOT NULL,
  body TEXT NOT NULL,
  created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS user_progress (
  user_email TEXT NOT NULL,
  unit_id TEXT NOT NULL REFERENCES units(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('not_started', 'in_progress', 'completed')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (user_email, unit_id)
);

CREATE TABLE IF NOT EXISTS audit_events (
  id TEXT PRIMARY KEY,
  actor_email TEXT NOT NULL,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  details_json TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_proposals_status_created ON proposals(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_comments_proposal_created ON proposal_comments(proposal_id, created_at);
CREATE INDEX IF NOT EXISTS idx_files_unit_approved ON proposal_files(target_unit_id, approved);
