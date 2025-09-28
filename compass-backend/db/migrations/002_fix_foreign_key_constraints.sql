-- Fix backwards foreign key constraints

-- Drop the incorrect foreign key constraints from projects table
ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_project_specifications_project;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_project_rfis_project;

-- The correct constraints should already exist on the child tables
-- Let's verify and add them if they don't exist

-- Check and add foreign key from project_specifications to projects
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_projects_specifications' 
        AND table_name = 'project_specifications'
    ) THEN
        ALTER TABLE project_specifications 
        ADD CONSTRAINT fk_projects_specifications 
        FOREIGN KEY (project_id) REFERENCES projects(project_id);
    END IF;
END $$;

-- Check and add foreign key from project_rfis to projects  
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints 
        WHERE constraint_name = 'fk_projects_rfis' 
        AND table_name = 'project_rfis'
    ) THEN
        ALTER TABLE project_rfis 
        ADD CONSTRAINT fk_projects_rfis 
        FOREIGN KEY (project_id) REFERENCES projects(project_id);
    END IF;
END $$;