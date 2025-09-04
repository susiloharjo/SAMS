-- SAMS Database Initialization Script
-- ISO 55001 Compliant Asset Management Schema

-- Create database if not exists
-- (PostgreSQL creates the database automatically via environment variables)

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create categories table for dynamic asset categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create assets table (ISO 55001 compliant)
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Basic Information (ISO 55001 - Asset Identification)
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    -- Technical Specifications
    type VARCHAR(100),
    model VARCHAR(100),
    serial_number VARCHAR(100) UNIQUE,
    manufacturer VARCHAR(100),
    
    -- Financial Information
    acquisition_cost DECIMAL(15,2),
    current_value DECIMAL(15,2),
    depreciation_rate DECIMAL(5,2),
    
    -- Operational Status
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'maintenance', 'disposed')),
    condition VARCHAR(50) DEFAULT 'good' CHECK (condition IN ('excellent', 'good', 'fair', 'poor', 'critical')),
    criticality VARCHAR(50) DEFAULT 'low' CHECK (criticality IN ('low', 'medium', 'high', 'critical')),
    
    -- Location Information
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    address TEXT,
    building_room VARCHAR(100),
    
    -- Lifecycle Information
    acquisition_date DATE,
    expected_life_years INTEGER,
    maintenance_schedule TEXT,
    
    -- Compliance and Standards
    certifications TEXT,
    standards TEXT,
    audit_info TEXT,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_assets_category_id ON assets(category_id);
CREATE INDEX IF NOT EXISTS idx_assets_status ON assets(status);
CREATE INDEX IF NOT EXISTS idx_assets_location ON assets(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_assets_serial_number ON assets(serial_number);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_categories_updated_at 
    BEFORE UPDATE ON categories 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assets_updated_at 
    BEFORE UPDATE ON assets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- Create users table for authentication and authorization
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    password VARCHAR(255) NOT NULL, -- Changed from password_hash to match GORM model
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('admin', 'manager', 'user')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create departments table
CREATE TABLE IF NOT EXISTS departments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_departments_updated_at
    BEFORE UPDATE ON departments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample categories for MVP demo with explicit UUIDs
INSERT INTO categories (id, name, description) VALUES
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'IT Equipment', 'Computers, servers, and peripherals'),
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b', 'Vehicles', 'Company cars, trucks, and other vehicles'),
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c', 'Furniture', 'Desks, chairs, and office furniture'),
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1d', 'Machinery', 'Industrial and manufacturing machinery'),
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1e', 'Software', 'Software licenses and subscriptions'),
    ('c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f', 'Real Estate', 'Land, buildings, and property')
ON CONFLICT (id) DO NOTHING;

-- Insert default admin user
INSERT INTO users (username, email, first_name, last_name, password, role) VALUES -- Changed from password_hash
    (
        'admin',
        'admin@sams.com',
        'System',
        'Administrator',
        '$2a$10$E/aBilYPDYLWs6ZvQwJ0nex1v7Wj1c0Hi.2M4GRkVuDsHngVsj71q', -- user.1001 (Generated from backend)
        'admin'
    )
ON CONFLICT (username) DO NOTHING;

-- Insert sample departments
INSERT INTO departments (name, description) VALUES
    ('Finance', 'Finance and Accounting Department'),
    ('Marketing', 'Marketing and Sales Department')
ON CONFLICT (name) DO NOTHING;


-- Insert sample assets for MVP demo (Government/Company focused)
INSERT INTO assets (name, description, category_id, type, model, serial_number, manufacturer, acquisition_cost, current_value, status, condition, criticality, latitude, longitude, address, building_room, acquisition_date, expected_life_years) VALUES
    ('Dell Latitude Laptop', 'Government issued laptop for office staff', 
     (SELECT id FROM categories WHERE name = 'IT Equipment'), 
     'Laptop', 'Latitude 5520', 'DL-001-2024', 'Dell Technologies', 
     1200.00, 900.00, 'active', 'good', 'medium', 
     40.7128, -74.0060, '123 Government Plaza, New York, NY', 'Floor 5, Room 501', 
     '2024-01-15', 4),
    
    ('HP LaserJet Printer', 'Office printer for document management', 
     (SELECT id FROM categories WHERE name = 'IT Equipment'), 
     'Printer', 'LaserJet Pro M404n', 'HP-002-2024', 'HP Inc.', 
     300.00, 250.00, 'active', 'excellent', 'low', 
     40.7128, -74.0060, '123 Government Plaza, New York, NY', 'Floor 5, Room 501', 
     '2024-01-20', 5),
    
    ('Ford Transit Van', 'Government vehicle for maintenance crew', 
     (SELECT id FROM categories WHERE name = 'Vehicles'), 
     'Van', 'Transit 350', 'FT-003-2024', 'Ford Motor Company', 
     35000.00, 32000.00, 'active', 'excellent', 'high', 
     40.7128, -74.0060, '123 Government Plaza, New York, NY', 'Parking Garage A', 
     '2024-01-10', 8),
    
    ('Office Building A', 'Main government office building', 
     (SELECT id FROM categories WHERE name = 'Buildings'), 
     'Office Building', 'Government Plaza A', 'OB-004-2024', 'Government Construction', 
     5000000.00, 5200000.00, 'active', 'excellent', 'critical', 
     40.7128, -74.0060, '123 Government Plaza, New York, NY', 'Main Building', 
     '2020-06-01', 50),
    
    ('Industrial Drill Press', 'Heavy machinery for maintenance department', 
     (SELECT id FROM categories WHERE name = 'Machinery'), 
     'Drill Press', 'DP-5000', 'ID-005-2024', 'Industrial Tools Co.', 
     2500.00, 2000.00, 'active', 'good', 'medium', 
     40.7128, -74.0060, '123 Government Plaza, New York, NY', 'Maintenance Shop', 
     '2023-08-15', 15)
ON CONFLICT (serial_number) DO NOTHING;

-- Insert sample assets
INSERT INTO assets (id, name, description, category_id, type, model, serial_number, manufacturer, acquisition_cost, current_value, status, condition, criticality, latitude, longitude, address, building_room, acquisition_date, expected_life_years, created_at, updated_at) VALUES
       ('550e8400-e29b-41d4-a716-446655440001', 'Dell Latitude 5520', 'Business laptop for office use', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Laptop', 'Latitude 5520', 'DL001', 'Dell', 1200000, 800000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Office Building A - Floor 3', '2023-01-15', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440002', 'HP LaserJet Pro M404n', 'Office printer for document printing', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Printer', 'LaserJet Pro M404n', 'HP001', 'HP', 2500000, 1800000, 'active', 'good', 'low', -6.2088, 106.8456, 'Jakarta', 'Office Building A - Floor 1', '2023-02-20', 7, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440003', 'Cisco Catalyst 2960', 'Network switch for office connectivity', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Network Equipment', 'Catalyst 2960', 'CS001', 'Cisco', 5000000, 3500000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Server Room', '2023-03-10', 8, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440004', 'Toyota Avanza', 'Company vehicle for transportation', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b', 'Vehicle', 'Avanza', 'TOY001', 'Toyota', 250000000, 200000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Parking Lot A', '2023-04-05', 10, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440005', 'Samsung Galaxy Tab S7', 'Tablet for field operations', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Tablet', 'Galaxy Tab S7', 'SG001', 'Samsung', 8000000, 6000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Office Building B - Floor 2', '2023-05-12', 4, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440006', 'Canon EOS R6', 'Professional camera for documentation', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Camera', 'EOS R6', 'CN001', 'Canon', 35000000, 28000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Media Room', '2023-06-18', 6, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440007', 'Apple MacBook Pro 16', 'Development workstation', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Laptop', 'MacBook Pro 16', 'AP001', 'Apple', 45000000, 38000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Development Lab', '2023-07-22', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440008', 'LG 55" OLED TV', 'Conference room display', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c', 'Display', 'OLED55C1', 'LG001', 'LG', 15000000, 12000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Conference Room A', '2023-08-30', 8, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440009', 'Bosch Drill Set', 'Construction tools', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1d', 'Tools', 'Professional Drill Set', 'BS001', 'Bosch', 3000000, 2200000, 'active', 'good', 'low', -6.2088, 106.8456, 'Jakarta', 'Tool Storage', '2023-09-14', 10, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440010', 'Yamaha PSR-E373', 'Digital piano for events', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f', 'Musical Instrument', 'PSR-E373', 'YM001', 'Yamaha', 8000000, 6500000, 'active', 'good', 'low', -6.2088, 106.8456, 'Jakarta', 'Event Hall', '2023-10-25', 12, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440011', 'Lenovo ThinkPad X1', 'Executive laptop', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Laptop', 'ThinkPad X1 Carbon', 'LN001', 'Lenovo', 28000000, 22000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Executive Office', '2023-11-08', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440012', 'Epson WorkForce Pro', 'Large format printer', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Printer', 'WorkForce Pro WF-3720', 'EP001', 'Epson', 12000000, 9000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Printing Room', '2023-12-12', 6, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440013', 'Honda CR-V', 'Management vehicle', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b', 'Vehicle', 'CR-V', 'HND001', 'Honda', 450000000, 380000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Executive Parking', '2024-01-15', 8, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440014', 'iPad Pro 12.9', 'Design team tablet', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Tablet', 'iPad Pro 12.9', 'AP002', 'Apple', 25000000, 20000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Design Studio', '2024-02-20', 4, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440015', 'Sony A7 IV', 'Marketing camera', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Camera', 'A7 IV', 'SN001', 'Sony', 40000000, 32000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Marketing Office', '2024-03-10', 7, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440016', 'Dell PowerEdge R740', 'Server for data processing', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Server', 'PowerEdge R740', 'DL002', 'Dell', 80000000, 65000000, 'active', 'good', 'critical', -6.2088, 106.8456, 'Jakarta', 'Data Center', '2024-04-05', 6, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440017', 'Microsoft Surface Pro', 'Sales team device', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Tablet', 'Surface Pro 8', 'MS001', 'Microsoft', 22000000, 18000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Sales Office', '2024-05-18', 4, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440018', 'Brother HL-L2350DW', 'Department printer', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Printer', 'HL-L2350DW', 'BR001', 'Brother', 4000000, 3000000, 'active', 'good', 'low', -6.2088, 106.8456, 'Jakarta', 'HR Department', '2024-06-22', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440019', 'Asus ROG Strix', 'Gaming station for testing', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Desktop', 'ROG Strix G15', 'AS001', 'Asus', 35000000, 28000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Testing Lab', '2024-07-30', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440020', 'JBL Professional', 'Audio system for events', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f', 'Audio Equipment', 'Professional Series', 'JB001', 'JBL', 18000000, 15000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Event Hall', '2024-08-14', 10, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440021', 'HP EliteBook 840', 'IT Support laptop', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Laptop', 'EliteBook 840 G8', 'HP002', 'HP', 25000000, 20000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'IT Department', '2024-09-05', 5, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440022', 'Canon imageRUNNER', 'Multifunction printer', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Printer', 'imageRUNNER 2630', 'CN002', 'Canon', 35000000, 28000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Main Office', '2024-10-12', 7, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440023', 'Nikon Z6 II', 'Product photography camera', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Camera', 'Z6 II', 'NK001', 'Nikon', 45000000, 36000000, 'active', 'good', 'high', -6.2088, 106.8456, 'Jakarta', 'Product Studio', '2024-11-18', 8, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440024', 'Dell OptiPlex', 'Reception computer', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a', 'Desktop', 'OptiPlex 7090', 'DL003', 'Dell', 15000000, 12000000, 'active', 'good', 'low', -6.2088, 106.8456, 'Jakarta', 'Reception Area', '2024-12-25', 6, NOW(), NOW()),
       ('550e8400-e29b-41d4-a716-446655440025', 'Samsung QLED TV', 'Meeting room display', 'c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c', 'Display', 'QLED 65" Q80T', 'SG002', 'Samsung', 25000000, 20000000, 'active', 'good', 'medium', -6.2088, 106.8456, 'Jakarta', 'Meeting Room B', '2025-01-08', 8, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Create a view for asset summary (useful for AI queries)
CREATE OR REPLACE VIEW asset_summary AS
SELECT 
    a.id,
    a.name,
    a.description,
    c.name as category_name,
    a.type,
    a.model,
    a.serial_number,
    a.status,
    a.condition,
    a.criticality,
    a.latitude,
    a.longitude,
    a.address,
    a.building_room,
    a.current_value,
    a.acquisition_date,
    a.expected_life_years
FROM assets a
LEFT JOIN categories c ON a.category_id = c.id;

-- Grant permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO sams_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO sams_user;
