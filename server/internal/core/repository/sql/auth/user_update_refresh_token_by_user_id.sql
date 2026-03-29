UPDATE %s
SET refresh_token = $1, refresh_token_expiry_time = $2
WHERE user_id = $3;