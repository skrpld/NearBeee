UPDATE %s
SET title = $1, content = $2
WHERE topic_id = $3 AND user_id = $4
RETURNING topic_id, user_id, title, content, idempotency_key, latitude, longitude, created_at, updated_at;