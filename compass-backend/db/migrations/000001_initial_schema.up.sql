-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user')),
    account_status VARCHAR(20) DEFAULT 'pending' CHECK (account_status IN ('pending', 'active', 'disabled')),
    invited_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    project_id BIGSERIAL PRIMARY KEY,
    project_name VARCHAR(200) NOT NULL,
    company_name VARCHAR(200),
    company_address TEXT,
    project_status VARCHAR(20) DEFAULT 'not_yet_started' CHECK (project_status IN ('not_yet_started', 'progress', 'completed')),
    project_type VARCHAR(20) CHECK (project_type IN ('windows', 'doors')),
    created_by BIGINT NOT NULL,
    last_updated_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_projects_creator FOREIGN KEY (created_by) REFERENCES users(user_id),
    CONSTRAINT fk_projects_last_updater FOREIGN KEY (last_updated_by) REFERENCES users(user_id)
);

-- Create project_specifications table
CREATE TABLE IF NOT EXISTS project_specifications (
    specification_id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL,
    version_no BIGINT NOT NULL,
    colour VARCHAR(100),
    ironmongery VARCHAR(150),
    u_value DECIMAL(5,2),
    g_value DECIMAL(5,2),
    vents VARCHAR(100),
    acoustics VARCHAR(100),
    sbd VARCHAR(100),
    pas24 VARCHAR(100),
    restrictors VARCHAR(100),
    special_comments TEXT,
    attachment_url TEXT,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_project_version UNIQUE (project_id, version_no),
    CONSTRAINT fk_project_specifications_creator FOREIGN KEY (created_by) REFERENCES users(user_id),
    CONSTRAINT fk_projects_specifications FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create project_rfis table
CREATE TABLE IF NOT EXISTS project_rfis (
    rfi_id BIGSERIAL PRIMARY KEY,
    project_id BIGINT NOT NULL,
    question_text TEXT NOT NULL,
    answer_value VARCHAR(10) CHECK (answer_value IN ('yes', 'no')),
    answered_by BIGINT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_project_rfis_answerer FOREIGN KEY (answered_by) REFERENCES users(user_id),
    CONSTRAINT fk_projects_rfis FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by);
CREATE INDEX IF NOT EXISTS idx_project_specifications_project_id ON project_specifications(project_id);
CREATE INDEX IF NOT EXISTS idx_project_rfis_project_id ON project_rfis(project_id);

-- Add foreign key for users.invited_by (self-referencing)
ALTER TABLE users ADD CONSTRAINT fk_users_invited_by FOREIGN KEY (invited_by) REFERENCES users(user_id);
