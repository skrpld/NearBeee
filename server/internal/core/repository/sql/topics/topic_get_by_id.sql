SELECT topic_id, user_id, title, content, idempotency_key, latitude, longitude, created_at, updated_at
FROM %s
WHERE topic_id = $1;