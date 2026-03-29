INSERT INTO %s (email, password_hash, refresh_token, refresh_token_expiry_time)
VALUES ($1, $2, $3, $4)
RETURNING user_id, email, password_hash, refresh_token, refresh_token_expiry_time;