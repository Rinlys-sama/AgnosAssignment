-- =============================================
-- Hospital Middleware Database Schema
-- =============================================

-- Hospitals table: stores hospital information
CREATE TABLE IF NOT EXISTS hospitals (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,  -- e.g., 'hospital_a', 'hospital_b'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Staff table: hospital employees who can log in and search patients
CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    hospital_id INTEGER NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Patients table: patient records linked to a hospital
CREATE TABLE IF NOT EXISTS patients (
    id SERIAL PRIMARY KEY,
    first_name_th VARCHAR(255),
    middle_name_th VARCHAR(255),
    last_name_th VARCHAR(255),
    first_name_en VARCHAR(255),
    middle_name_en VARCHAR(255),
    last_name_en VARCHAR(255),
    date_of_birth DATE,
    patient_hn VARCHAR(50),
    national_id VARCHAR(20),
    passport_id VARCHAR(50),
    phone_number VARCHAR(20),
    email VARCHAR(255),
    gender VARCHAR(1) CHECK (gender IN ('M', 'F')),
    hospital_id INTEGER NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for faster searching
CREATE INDEX IF NOT EXISTS idx_patients_national_id ON patients(national_id);
CREATE INDEX IF NOT EXISTS idx_patients_passport_id ON patients(passport_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients(hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_name_en ON patients(first_name_en, last_name_en);
CREATE INDEX IF NOT EXISTS idx_staff_hospital_id ON staff(hospital_id);

-- Seed data: create some hospitals
INSERT INTO hospitals (name, code) VALUES
    ('Hospital A', 'hospital_a'),
    ('Hospital B', 'hospital_b')
ON CONFLICT (code) DO NOTHING;

-- Seed data: create some test patients
INSERT INTO patients (first_name_th, last_name_th, first_name_en, last_name_en, date_of_birth, patient_hn, national_id, passport_id, phone_number, email, gender, hospital_id) VALUES
    ('สมชาย', 'ใจดี', 'Somchai', 'Jaidee', '1990-01-15', 'HN0001', '1234567890123', 'AA1234567', '0812345678', 'somchai@email.com', 'M', 1),
    ('สมหญิง', 'รักไทย', 'Somying', 'Rakthai', '1985-06-20', 'HN0002', '9876543210987', 'BB9876543', '0898765432', 'somying@email.com', 'F', 1),
    ('วิชัย', 'สุขสันต์', 'Wichai', 'Suksan', '1992-03-10', 'HN0003', '1111222233334', NULL, '0855551234', 'wichai@email.com', 'M', 2)
ON CONFLICT DO NOTHING;
