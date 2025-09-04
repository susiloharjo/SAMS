-- Enable UUID extension for PostgreSQL
-- This script ensures uuid_generate_v4() function is available

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
