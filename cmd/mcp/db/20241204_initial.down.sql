-- The audit log tracks all of the JSON-RPC operations for each session.
-- It tracks the session ID and the JSON-RPC operation.
DROP INDEX IF EXISTS idx_audit_log_session_id;
DROP INDEX IF EXISTS idx_audit_log_timestamp;
DROP INDEX IF EXISTS idx_audit_log_request_id;
DROP TABLE IF NOT EXISTS audit_log (
  id SERIAL PRIMARY KEY,
  session_id TEXT NOT NULL,
  type TEXT NOT NULL,
  timestamp TIMESTAMPTZ DEFAULT NOW(),
  jsonrpc_version TEXT NOT NULL,
  operation JSONB NOT NULL,
  request_id TEXT
);
-- Adding indices
-- The installed integrations table tracks the integrations that are installed.
DROP INDEX IF EXISTS idx_integrations_name;
DROP INDEX IF EXISTS idx_integrations_vendor;
DROP TABLE IF EXISTS IF NOT EXISTS integrations;
-- Adding indices
