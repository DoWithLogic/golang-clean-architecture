UPDATE users SET 
    fullname = CASE WHEN ? != '' THEN ? ELSE fullname END,
    phone_number = CASE WHEN ? != '' THEN ? ELSE phone_number END,
    user_type = CASE WHEN ? != '' THEN ? ELSE user_type END,
    updated_at = ?
WHERE id = ?