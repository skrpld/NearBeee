SELECT user_id, email, password_hash, refresh_token, refresh_token_expiry_time
FROM %s
WHERE user_id = $1;