INSERT INTO %s (user_id, title, content, idempotency_key, latitude, longitude)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING topic_id, user_id, title, content, idempotency_key, latitude, longitude, created_at, updated_at;