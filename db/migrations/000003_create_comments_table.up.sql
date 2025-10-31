CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    author_id BIGSERIAL REFERENCES users(id),
    body VARCHAR(2000) NOT NULL,
    parent_post_id BIGSERIAL REFERENCES posts(id),
    parent_comment_id BIGSERIAL REFERENCES comments(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);