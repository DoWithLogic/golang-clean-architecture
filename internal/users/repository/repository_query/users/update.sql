UPDATE users SET 
    fullname = CASE WHEN $2 != '' THEN $2 ELSE fullname END,
    phone_number = CASE WHEN $3 != '' THEN $3 ELSE phone_number END,
    user_type = CASE WHEN $4 != '' THEN $4 ELSE user_type END,
    is_active = CASE WHEN $5 IS NULL THEN $5 ELSE user_type END,
    updated_at = ?,
    created_by = ?
WHERE id = $1