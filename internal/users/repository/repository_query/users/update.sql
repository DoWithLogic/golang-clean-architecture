UPDATE users SET 
    email = CASE WHEN ? != '' THEN ? ELSE email END,
    fullname = CASE WHEN ? != '' THEN ? ELSE fullname END,
    phone_number = CASE WHEN ? != '' THEN ? ELSE phone_number END,
    user_type = CASE WHEN ? != '' THEN ? ELSE user_type END,
    updated_at = ?,
    updated_by = ?
WHERE id = ?