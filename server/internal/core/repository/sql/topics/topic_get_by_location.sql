SELECT topic_id, user_id, title, content, idempotency_key, latitude, longitude, created_at, updated_at
FROM %s
WHERE calculate_distance($1, $2, latitude, longitude) <= $3
ORDER BY calculate_distance($1, $2, latitude, longitude)
LIMIT $4;