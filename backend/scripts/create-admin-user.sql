-- Create default admin user
-- Password: admin123 (hashed with bcrypt)
-- Username: admin
-- Email: admin@sams.com
-- Role: admin

INSERT INTO users (
    id,
    username,
    email,
    first_name,
    last_name,
    password,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    'admin',
    'admin@sams.com',
    'System',
    'Administrator',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- admin123
    'admin',
    true,
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Create default manager user
-- Password: manager123 (hashed with bcrypt)
-- Username: manager
-- Email: manager@sams.com
-- Role: manager

INSERT INTO users (
    id,
    username,
    email,
    first_name,
    last_name,
    password,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    'manager',
    'manager@sams.com',
    'System',
    'Manager',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- manager123
    'manager',
    true,
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;

-- Create default regular user
-- Password: user123 (hashed with bcrypt)
-- Username: user
-- Email: user@sams.com
-- Role: user

INSERT INTO users (
    id,
    username,
    email,
    first_name,
    last_name,
    password,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    'user',
    'user@sams.com',
    'System',
    'User',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- user123
    'user',
    true,
    NOW(),
    NOW()
) ON CONFLICT (username) DO NOTHING;
