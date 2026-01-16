-- Remove tonic and botanicals fields from tasting_sessions
ALTER TABLE tasting_sessions
DROP COLUMN tonic,
DROP COLUMN botanicals;
