-- Migration to convert boolean fields to text in project_specifications table
-- This needs to be run manually after updating the model

-- First, add temporary columns
ALTER TABLE project_specifications ADD COLUMN sbd_temp VARCHAR(100);
ALTER TABLE project_specifications ADD COLUMN pas24_temp VARCHAR(100);
ALTER TABLE project_specifications ADD COLUMN restrictors_temp VARCHAR(100);

-- Convert boolean values to text
UPDATE project_specifications 
SET 
    sbd_temp = CASE WHEN sbd = true THEN 'Yes' ELSE 'No' END,
    pas24_temp = CASE WHEN pas24 = true THEN 'Yes' ELSE 'No' END,
    restrictors_temp = CASE WHEN restrictors = true THEN 'Yes' ELSE 'No' END;

-- Drop old columns
ALTER TABLE project_specifications DROP COLUMN sbd;
ALTER TABLE project_specifications DROP COLUMN pas24;
ALTER TABLE project_specifications DROP COLUMN restrictors;

-- Rename temporary columns
ALTER TABLE project_specifications RENAME COLUMN sbd_temp TO sbd;
ALTER TABLE project_specifications RENAME COLUMN pas24_temp TO pas24;
ALTER TABLE project_specifications RENAME COLUMN restrictors_temp TO restrictors;