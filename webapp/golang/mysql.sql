CREATE INDEX post_id_created_at_key ON comments(post_id, created_at);
CREATE INDEX created_at_key ON posts(created_at);
