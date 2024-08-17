-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS documents
(
    id         uuid PRIMARY KEY         DEFAULT gen_random_uuid() NOT NULL,
    file_path  text UNIQUE                                        NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS document_pages
(
    id          uuid PRIMARY KEY         DEFAULT gen_random_uuid() NOT NULL,
    page        int                                                NOT NULL,
    text        text                                               NOT NULL,
    embeddings  vector(1536)                                       NOT NULL,
    document_id uuid                                               NOT NULL,
    created_at  timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents (id),
    UNIQUE (document_id, page)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS document_pages;
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
