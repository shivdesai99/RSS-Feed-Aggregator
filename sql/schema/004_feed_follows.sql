-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Delete the follow if the user is deleted
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE, -- Delete the follow if the feed is deleted
    UNIQUE(user_id, feed_id) -- Integrity Restraint to prevent duplicate follows
);


-- +goose Down
DROP TABLE feed_follows;
