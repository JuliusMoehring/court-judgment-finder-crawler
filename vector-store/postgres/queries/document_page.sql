-- name: CreateDocumentPage :one
INSERT INTO document_pages (page, text, embeddings, document_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (document_id, page) DO UPDATE
    SET document_id = $4,
        page        = $1,
        updated_at  = CURRENT_TIMESTAMP
RETURNING id;