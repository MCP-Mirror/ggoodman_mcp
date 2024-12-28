CREATE TABLE IF NOT EXISTS integrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  vendor TEXT NOT NULL,
  source_url TEXT NOT NULL,
  homepage TEXT NOT NULL,
  license TEXT NOT NULL,
  instructions BLOB NOT NULL
);
CREATE INDEX idx_integrations_name ON integrations (name);
CREATE INDEX idx_integrations_vendor ON integrations (vendor);
