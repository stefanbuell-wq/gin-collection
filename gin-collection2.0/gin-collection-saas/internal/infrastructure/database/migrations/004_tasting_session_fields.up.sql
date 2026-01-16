-- Add tonic and botanicals fields to tasting_sessions
ALTER TABLE tasting_sessions
ADD COLUMN tonic VARCHAR(255) NULL COMMENT 'Tonic water used during tasting',
ADD COLUMN botanicals TEXT NULL COMMENT 'Botanicals noticed (comma-separated)';
