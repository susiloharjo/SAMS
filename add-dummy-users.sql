-- Add dummy users for testing
-- Note: Passwords are hashed versions of "password123" using bcrypt

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'admin',
    'admin@sams.com',
    'System',
    'Administrator',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'admin',
    (SELECT id FROM departments WHERE name = 'Information Technology (IT)' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'manager.finance',
    'finance.manager@sams.com',
    'Sarah',
    'Johnson',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'manager',
    (SELECT id FROM departments WHERE name = 'Finance' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'user.it',
    'it.user@sams.com',
    'Michael',
    'Chen',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'user',
    (SELECT id FROM departments WHERE name = 'Information Technology (IT)' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'manager.hr',
    'hr.manager@sams.com',
    'Emily',
    'Rodriguez',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'manager',
    (SELECT id FROM departments WHERE name = 'Human Capital (HC)' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'user.project',
    'project.user@sams.com',
    'David',
    'Kim',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'user',
    (SELECT id FROM departments WHERE name = 'Project' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'manager.ops',
    'ops.manager@sams.com',
    'Lisa',
    'Thompson',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'manager',
    (SELECT id FROM departments WHERE name = 'Operation' LIMIT 1),
    true,
    NOW(),
    NOW()
);

INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'user.marketing',
    'marketing.user@sams.com',
    'James',
    'Wilson',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'user',
    (SELECT id FROM departments WHERE name = 'Marketing' LIMIT 1),
    true,
    NOW(),
    NOW()
);

-- Add a few more users without departments
INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'guest.user',
    'guest@sams.com',
    'Guest',
    'User',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'user',
    NULL,
    true,
    NOW(),
    NOW()
);

-- Add an inactive user for testing
INSERT INTO users (id, username, email, first_name, last_name, password, role, department_id, is_active, created_at, updated_at) VALUES
(
    gen_random_uuid(),
    'inactive.user',
    'inactive@sams.com',
    'Inactive',
    'User',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password123
    'user',
    (SELECT id FROM departments WHERE name = 'Information Technology (IT)' LIMIT 1),
    false,
    NOW(),
    NOW()
);

-- Display the created users
SELECT 
    u.username,
    u.email,
    u.first_name,
    u.last_name,
    u.role,
    d.name as department,
    u.is_active,
    u.created_at
FROM users u
LEFT JOIN departments d ON u.department_id = d.id
ORDER BY u.created_at DESC;
