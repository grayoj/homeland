CREATE TABLE incidents (
    id SERIAL PRIMARY KEY,
    agent_id VARCHAR(255) NOT NULL,
    department VARCHAR(255) NOT NULL,
    incident_type VARCHAR(50) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    caller_full_name VARCHAR(255) NOT NULL,
    caller_phone_number VARCHAR(20) NOT NULL,
    caller_location VARCHAR(255) NOT NULL,
    people_involved INT NOT NULL,
    incident_report TEXT NOT NULL,
    staff_id INT NOT NULL REFERENCES staff(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

