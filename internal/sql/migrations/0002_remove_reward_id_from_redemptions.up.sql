-- Remove foreign key constraint first
ALTER TABLE redemptions DROP CONSTRAINT IF EXISTS redemptions_reward_id_fkey;

-- Drop the index if it exists
DROP INDEX IF EXISTS idx_redemptions_reward_id;

-- Remove the column
ALTER TABLE redemptions DROP COLUMN reward_id;

