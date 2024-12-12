-- The audit log tracks all of the JSON-RPC operations for each session.
-- It tracks the session ID and the JSON-RPC operation.
CREATE TABLE IF NOT EXISTS audit_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  session_id TEXT NOT NULL,
  type TEXT NOT NULL,
  timestamp TIMESTAMP CURRENT_TIMESTAMP,
  jsonrpc_version TEXT NOT NULL,
  operation JSONB NOT NULL,
  request_id TEXT
);
-- Adding indices
CREATE INDEX idx_audit_log_session_id ON audit_log (session_id);
CREATE INDEX idx_audit_log_timestamp ON audit_log (timestamp);
CREATE INDEX idx_audit_log_request_id ON audit_log (request_id);
-- The installed integrations table tracks the integrations that are installed.
CREATE TABLE IF NOT EXISTS integrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT,
  vendor TEXT,
  source_url TEXT,
  homepage TEXT,
  license TEXT,
  runtime TEXT
);
-- Adding indices
CREATE INDEX idx_integrations_name ON integrations (name);
CREATE INDEX idx_integrations_vendor ON integrations (vendor);
