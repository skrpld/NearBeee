SELECT topic_id, user_id, title, content, idempotency_key, latitude, longitude, created_at, updated_at
FROM %s
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2;