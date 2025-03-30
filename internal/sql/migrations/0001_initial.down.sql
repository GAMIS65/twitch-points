DROP TRIGGER IF EXISTS update_users_modtime ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_redemptions_viewer_id;
DROP TABLE IF EXISTS redemptions;
DROP TABLE IF EXISTS rewards;
DROP TABLE IF EXISTS viewers;
DROP TABLE IF EXISTS streamers;

