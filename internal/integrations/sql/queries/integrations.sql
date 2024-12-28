-- name: ListIntegrations :many
SELECT id,
  name,
  description,
  vendor,
  source_url,
  homepage,
  license,
  instructions
FROM integrations;
-- name: CreateIntegration :one
INSERT INTO integrations (
    name,
    description,
    vendor,
    source_url,
    homepage,
    license,
    instructions
  )
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;
