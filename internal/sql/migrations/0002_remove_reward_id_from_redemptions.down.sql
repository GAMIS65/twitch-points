-- Add the column back
ALTER TABLE redemptions ADD COLUMN reward_id TEXT;

-- Recreate the foreign key constraint
ALTER TABLE redemptions ADD CONSTRAINT redemptions_reward_id_fkey 
FOREIGN KEY (reward_id) REFERENCES rewards(reward_id);

-- Recreate any indexes if needed
CREATE INDEX idx_redemptions_reward_id ON redemptions (reward_id);