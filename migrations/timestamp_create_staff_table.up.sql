CREATE TABLE staff (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    middle_name VARCHAR(255),
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    agent_id VARCHAR(255) UNIQUE NOT NULL,
    profile_photo TEXT,
    position VARCHAR(50) NOT NULL CHECK (position IN ('SSA', 'Director', 'IT', 'Call Center', 'Staff')),
    address TEXT,
    department VARCHAR(50) NOT NULL CHECK (department IN ('Homeland Security', 'AVS', 'EMS', 'Fire Service')),
    date_of_birth DATE NOT NULL,
    state_of_origin VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('Admin', 'SSA', 'Director', 'Staff')),
    must_change_password BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

