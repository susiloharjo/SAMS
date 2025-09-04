--
-- PostgreSQL database dump
--

\restrict qQB60KSy2a3RubyMJNb9QHmybupbQGIim8VRuU3SC92OGKI6eUymkzPSQ6ZfpAg

-- Dumped from database version 16.10
-- Dumped by pg_dump version 16.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: update_updated_at_column(); Type: FUNCTION; Schema: public; Owner: sams_user
--

CREATE FUNCTION public.update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_updated_at_column() OWNER TO sams_user;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: assets; Type: TABLE; Schema: public; Owner: sams_user
--

CREATE TABLE public.assets (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    category_id uuid,
    type character varying(100),
    model character varying(100),
    serial_number character varying(100),
    manufacturer character varying(100),
    acquisition_cost numeric(15,2),
    current_value numeric(15,2),
    depreciation_rate numeric(5,2),
    status character varying(50) DEFAULT 'active'::character varying,
    condition character varying(50) DEFAULT 'good'::character varying,
    criticality character varying(50) DEFAULT 'low'::character varying,
    latitude numeric(10,8),
    longitude numeric(11,8),
    address text,
    building_room character varying(100),
    acquisition_date date,
    expected_life_years integer,
    maintenance_schedule text,
    certifications text,
    standards text,
    audit_info text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    department_id uuid,
    deleted_at timestamp with time zone,
    CONSTRAINT assets_condition_check CHECK (((condition)::text = ANY ((ARRAY['excellent'::character varying, 'good'::character varying, 'fair'::character varying, 'poor'::character varying, 'critical'::character varying])::text[]))),
    CONSTRAINT assets_criticality_check CHECK (((criticality)::text = ANY ((ARRAY['low'::character varying, 'medium'::character varying, 'high'::character varying, 'critical'::character varying])::text[]))),
    CONSTRAINT assets_status_check CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'inactive'::character varying, 'maintenance'::character varying, 'disposed'::character varying])::text[]))),
    CONSTRAINT chk_assets_condition CHECK (((condition)::text = ANY ((ARRAY['excellent'::character varying, 'good'::character varying, 'fair'::character varying, 'poor'::character varying, 'critical'::character varying])::text[]))),
    CONSTRAINT chk_assets_criticality CHECK (((criticality)::text = ANY ((ARRAY['low'::character varying, 'medium'::character varying, 'high'::character varying, 'critical'::character varying])::text[]))),
    CONSTRAINT chk_assets_status CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'inactive'::character varying, 'maintenance'::character varying, 'disposed'::character varying])::text[])))
);


ALTER TABLE public.assets OWNER TO sams_user;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: sams_user
--

CREATE TABLE public.categories (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.categories OWNER TO sams_user;

--
-- Name: asset_summary; Type: VIEW; Schema: public; Owner: sams_user
--

CREATE VIEW public.asset_summary AS
 SELECT a.id,
    a.name,
    a.description,
    c.name AS category_name,
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
   FROM (public.assets a
     LEFT JOIN public.categories c ON ((a.category_id = c.id)));


ALTER VIEW public.asset_summary OWNER TO sams_user;

--
-- Name: departments; Type: TABLE; Schema: public; Owner: sams_user
--

CREATE TABLE public.departments (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);


ALTER TABLE public.departments OWNER TO sams_user;

--
-- Name: users; Type: TABLE; Schema: public; Owner: sams_user
--

CREATE TABLE public.users (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(100) NOT NULL,
    first_name character varying(50) NOT NULL,
    last_name character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    role character varying(20) DEFAULT 'user'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    department_id uuid,
    is_active boolean DEFAULT true,
    last_login timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT users_role_check CHECK (((role)::text = ANY (ARRAY[('admin'::character varying)::text, ('manager'::character varying)::text, ('user'::character varying)::text])))
);


ALTER TABLE public.users OWNER TO sams_user;

--
-- Data for Name: assets; Type: TABLE DATA; Schema: public; Owner: sams_user
--

COPY public.assets (id, name, description, category_id, type, model, serial_number, manufacturer, acquisition_cost, current_value, depreciation_rate, status, condition, criticality, latitude, longitude, address, building_room, acquisition_date, expected_life_years, maintenance_schedule, certifications, standards, audit_info, created_at, updated_at, department_id, deleted_at) FROM stdin;
7fb15206-6a17-407b-9d05-1b09b83ca722	Dell Latitude Laptop	Government issued laptop for office staff	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Laptop	Latitude 5520	DL-001-2024	Dell Technologies	1200.00	900.00	\N	active	good	medium	40.71280000	-74.00600000	123 Government Plaza, New York, NY	Floor 5, Room 501	2024-01-15	4	\N	\N	\N	\N	2025-09-04 12:42:02.949516+00	2025-09-04 12:42:02.949516+00	\N	\N
8731545e-e68e-4ac4-a9cd-216415e34927	HP LaserJet Printer	Office printer for document management	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Printer	LaserJet Pro M404n	HP-002-2024	HP Inc.	300.00	250.00	\N	active	excellent	low	40.71280000	-74.00600000	123 Government Plaza, New York, NY	Floor 5, Room 501	2024-01-20	5	\N	\N	\N	\N	2025-09-04 12:42:02.949516+00	2025-09-04 12:42:02.949516+00	\N	\N
255218aa-b00e-46f3-a069-0ff08ef5a982	Ford Transit Van	Government vehicle for maintenance crew	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b	Van	Transit 350	FT-003-2024	Ford Motor Company	35000.00	32000.00	\N	active	excellent	high	40.71280000	-74.00600000	123 Government Plaza, New York, NY	Parking Garage A	2024-01-10	8	\N	\N	\N	\N	2025-09-04 12:42:02.949516+00	2025-09-04 12:42:02.949516+00	\N	\N
8f039ad9-db40-4731-8605-3fc0222a8000	Office Building A	Main government office building	\N	Office Building	Government Plaza A	OB-004-2024	Government Construction	5000000.00	5200000.00	\N	active	excellent	critical	40.71280000	-74.00600000	123 Government Plaza, New York, NY	Main Building	2020-06-01	50	\N	\N	\N	\N	2025-09-04 12:42:02.949516+00	2025-09-04 12:42:02.949516+00	\N	\N
d956c38a-afbb-4a58-8654-27fc25492c4e	Industrial Drill Press	Heavy machinery for maintenance department	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1d	Drill Press	DP-5000	ID-005-2024	Industrial Tools Co.	2500.00	2000.00	\N	active	good	medium	40.71280000	-74.00600000	123 Government Plaza, New York, NY	Maintenance Shop	2023-08-15	15	\N	\N	\N	\N	2025-09-04 12:42:02.949516+00	2025-09-04 12:42:02.949516+00	\N	\N
550e8400-e29b-41d4-a716-446655440001	Dell Latitude 5520	Business laptop for office use	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Laptop	Latitude 5520	DL001	Dell	1200000.00	800000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Office Building A - Floor 3	2023-01-15	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440002	HP LaserJet Pro M404n	Office printer for document printing	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Printer	LaserJet Pro M404n	HP001	HP	2500000.00	1800000.00	\N	active	good	low	-6.20880000	106.84560000	Jakarta	Office Building A - Floor 1	2023-02-20	7	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440003	Cisco Catalyst 2960	Network switch for office connectivity	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Network Equipment	Catalyst 2960	CS001	Cisco	5000000.00	3500000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Server Room	2023-03-10	8	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440004	Toyota Avanza	Company vehicle for transportation	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b	Vehicle	Avanza	TOY001	Toyota	250000000.00	200000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Parking Lot A	2023-04-05	10	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440005	Samsung Galaxy Tab S7	Tablet for field operations	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Tablet	Galaxy Tab S7	SG001	Samsung	8000000.00	6000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Office Building B - Floor 2	2023-05-12	4	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440006	Canon EOS R6	Professional camera for documentation	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Camera	EOS R6	CN001	Canon	35000000.00	28000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Media Room	2023-06-18	6	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440007	Apple MacBook Pro 16	Development workstation	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Laptop	MacBook Pro 16	AP001	Apple	45000000.00	38000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Development Lab	2023-07-22	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440008	LG 55" OLED TV	Conference room display	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c	Display	OLED55C1	LG001	LG	15000000.00	12000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Conference Room A	2023-08-30	8	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440009	Bosch Drill Set	Construction tools	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1d	Tools	Professional Drill Set	BS001	Bosch	3000000.00	2200000.00	\N	active	good	low	-6.20880000	106.84560000	Jakarta	Tool Storage	2023-09-14	10	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440010	Yamaha PSR-E373	Digital piano for events	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f	Musical Instrument	PSR-E373	YM001	Yamaha	8000000.00	6500000.00	\N	active	good	low	-6.20880000	106.84560000	Jakarta	Event Hall	2023-10-25	12	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440011	Lenovo ThinkPad X1	Executive laptop	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Laptop	ThinkPad X1 Carbon	LN001	Lenovo	28000000.00	22000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Executive Office	2023-11-08	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440012	Epson WorkForce Pro	Large format printer	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Printer	WorkForce Pro WF-3720	EP001	Epson	12000000.00	9000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Printing Room	2023-12-12	6	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440013	Honda CR-V	Management vehicle	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b	Vehicle	CR-V	HND001	Honda	450000000.00	380000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Executive Parking	2024-01-15	8	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440014	iPad Pro 12.9	Design team tablet	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Tablet	iPad Pro 12.9	AP002	Apple	25000000.00	20000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Design Studio	2024-02-20	4	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440015	Sony A7 IV	Marketing camera	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Camera	A7 IV	SN001	Sony	40000000.00	32000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Marketing Office	2024-03-10	7	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440016	Dell PowerEdge R740	Server for data processing	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Server	PowerEdge R740	DL002	Dell	80000000.00	65000000.00	\N	active	good	critical	-6.20880000	106.84560000	Jakarta	Data Center	2024-04-05	6	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440017	Microsoft Surface Pro	Sales team device	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Tablet	Surface Pro 8	MS001	Microsoft	22000000.00	18000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Sales Office	2024-05-18	4	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440018	Brother HL-L2350DW	Department printer	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Printer	HL-L2350DW	BR001	Brother	4000000.00	3000000.00	\N	active	good	low	-6.20880000	106.84560000	Jakarta	HR Department	2024-06-22	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440019	Asus ROG Strix	Gaming station for testing	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Desktop	ROG Strix G15	AS001	Asus	35000000.00	28000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Testing Lab	2024-07-30	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440020	JBL Professional	Audio system for events	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f	Audio Equipment	Professional Series	JB001	JBL	18000000.00	15000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Event Hall	2024-08-14	10	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440021	HP EliteBook 840	IT Support laptop	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Laptop	EliteBook 840 G8	HP002	HP	25000000.00	20000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	IT Department	2024-09-05	5	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440022	Canon imageRUNNER	Multifunction printer	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Printer	imageRUNNER 2630	CN002	Canon	35000000.00	28000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Main Office	2024-10-12	7	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440023	Nikon Z6 II	Product photography camera	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Camera	Z6 II	NK001	Nikon	45000000.00	36000000.00	\N	active	good	high	-6.20880000	106.84560000	Jakarta	Product Studio	2024-11-18	8	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440024	Dell OptiPlex	Reception computer	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	Desktop	OptiPlex 7090	DL003	Dell	15000000.00	12000000.00	\N	active	good	low	-6.20880000	106.84560000	Jakarta	Reception Area	2024-12-25	6	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
550e8400-e29b-41d4-a716-446655440025	Samsung QLED TV	Meeting room display	c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c	Display	QLED 65" Q80T	SG002	Samsung	25000000.00	20000000.00	\N	active	good	medium	-6.20880000	106.84560000	Jakarta	Meeting Room B	2025-01-08	8	\N	\N	\N	\N	2025-09-04 12:42:02.95382+00	2025-09-04 12:42:02.95382+00	\N	\N
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: sams_user
--

COPY public.categories (id, name, description, created_at, updated_at, deleted_at) FROM stdin;
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1a	IT Equipment	Computers, servers, and peripherals	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1b	Vehicles	Company cars, trucks, and other vehicles	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1c	Furniture	Desks, chairs, and office furniture	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1d	Machinery	Industrial and manufacturing machinery	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1e	Software	Software licenses and subscriptions	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
c1a7e21a-3e7a-4e1a-8c1a-2e4b6e8a2d1f	Real Estate	Land, buildings, and property	2025-09-04 12:42:02.942373+00	2025-09-04 12:42:02.942373+00	\N
\.


--
-- Data for Name: departments; Type: TABLE DATA; Schema: public; Owner: sams_user
--

COPY public.departments (id, name, description, created_at, updated_at, deleted_at) FROM stdin;
a0b217f2-17d6-4ca0-99ad-61663cce2dab	Finance	Finance and Accounting Department	2025-09-04 12:42:02.947335+00	2025-09-04 12:42:02.947335+00	\N
a1c1baca-49a1-4c71-8893-4e578604dbf1	Marketing	Marketing and Sales Department	2025-09-04 12:42:02.947335+00	2025-09-04 12:42:02.947335+00	\N
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: sams_user
--

COPY public.users (id, username, email, first_name, last_name, password, role, created_at, updated_at, department_id, is_active, last_login, deleted_at) FROM stdin;
06ae57f7-f0f0-4416-92a1-efe5db981f1d	admin	admin@sams.com	System	Administrator	$2a$10$E/aBilYPDYLWs6ZvQwJ0nex1v7Wj1c0Hi.2M4GRkVuDsHngVsj71q	admin	2025-09-04 12:42:02.945023+00	2025-09-04 12:42:02.945023+00	\N	t	\N	\N
\.


--
-- Name: assets assets_pkey; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_pkey PRIMARY KEY (id);


--
-- Name: assets assets_serial_number_key; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_serial_number_key UNIQUE (serial_number);


--
-- Name: categories categories_name_key; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_name_key UNIQUE (name);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: departments departments_name_key; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_name_key UNIQUE (name);


--
-- Name: departments departments_pkey; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.departments
    ADD CONSTRAINT departments_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_assets_category_id; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_assets_category_id ON public.assets USING btree (category_id);


--
-- Name: idx_assets_deleted_at; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_assets_deleted_at ON public.assets USING btree (deleted_at);


--
-- Name: idx_assets_location; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_assets_location ON public.assets USING btree (latitude, longitude);


--
-- Name: idx_assets_serial_number; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_assets_serial_number ON public.assets USING btree (serial_number);


--
-- Name: idx_assets_status; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_assets_status ON public.assets USING btree (status);


--
-- Name: idx_categories_deleted_at; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_categories_deleted_at ON public.categories USING btree (deleted_at);


--
-- Name: idx_departments_deleted_at; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_departments_deleted_at ON public.departments USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: sams_user
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: assets update_assets_updated_at; Type: TRIGGER; Schema: public; Owner: sams_user
--

CREATE TRIGGER update_assets_updated_at BEFORE UPDATE ON public.assets FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: categories update_categories_updated_at; Type: TRIGGER; Schema: public; Owner: sams_user
--

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: departments update_departments_updated_at; Type: TRIGGER; Schema: public; Owner: sams_user
--

CREATE TRIGGER update_departments_updated_at BEFORE UPDATE ON public.departments FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: users update_users_updated_at; Type: TRIGGER; Schema: public; Owner: sams_user
--

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();


--
-- Name: assets assets_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE SET NULL;


--
-- Name: assets fk_assets_department; Type: FK CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT fk_assets_department FOREIGN KEY (department_id) REFERENCES public.departments(id);


--
-- Name: assets fk_categories_assets; Type: FK CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT fk_categories_assets FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: users fk_users_department; Type: FK CONSTRAINT; Schema: public; Owner: sams_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_department FOREIGN KEY (department_id) REFERENCES public.departments(id);


--
-- PostgreSQL database dump complete
--

\unrestrict qQB60KSy2a3RubyMJNb9QHmybupbQGIim8VRuU3SC92OGKI6eUymkzPSQ6ZfpAg

