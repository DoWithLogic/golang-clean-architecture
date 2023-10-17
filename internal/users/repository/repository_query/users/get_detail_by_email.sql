SELECT
    u.id,
    u.email,
    u.password
FROM users u
WHERE u.email = ?