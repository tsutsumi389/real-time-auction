-- Migration 007: Remove automatic point history trigger
-- This trigger causes duplicate entries when the application
-- explicitly creates point_history records

DROP TRIGGER IF EXISTS trigger_record_point_history ON bidder_points;
DROP FUNCTION IF EXISTS record_point_history();
