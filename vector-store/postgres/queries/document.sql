-- name: CreateDocument :one
INSERT INTO documents (file_path)
VALUES ($1)
ON CONFLICT (file_path) DO UPDATE
    SET updated_at = CURRENT_TIMESTAMP
RETURNING id;

-- name: GetDocumentIDByFilePath :one
SELECT id
FROM documents
WHERE file_path = $1;