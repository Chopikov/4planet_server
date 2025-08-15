-- Add cover_url field to projects table
-- This allows projects to have a cover image similar to news articles

ALTER TABLE projects ADD COLUMN cover_url text;
