SELECT 
    u.id,
    u.email,
    u.fullname,
    u.phone_number,
    u.user_type,
    u.is_active,
    u.created_at
FROM users u
WHERE u.id = ?
