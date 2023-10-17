INSERT INTO users (
    email,
    password,
    fullname,
    phone_number,
    user_type,
    is_active,
    created_at,
    created_by
)VALUE (?, ?, ?, ?, ?, ?, ?, ?)