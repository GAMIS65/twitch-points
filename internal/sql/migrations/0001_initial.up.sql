CREATE TABLE streamers(
	twitch_id TEXT PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	verified BOOLEAN DEFAULT FALSE,
	profile_image_url TEXT,
	access_token TEXT,
	refresh_token TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE viewers(
	twitch_id TEXT PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	registered_in TEXT REFERENCES streamers(twitch_id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE rewards(
	reward_id TEXT PRIMARY KEY,
	streamer_id TEXT REFERENCES streamers(twitch_id),
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE redemptions(
	message_id TEXT PRIMARY KEY,
	reward_id TEXT REFERENCES rewards(reward_id),
	streamer_id TEXT REFERENCES streamers(twitch_id),
	viewer_id TEXT REFERENCES viewers(twitch_id),
	redeemed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_redemptions_viewer_id ON redemptions (viewer_id);

-- Function to automatically update the 'updated_at' column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function before updating a row
CREATE TRIGGER update_streamers_modtime
BEFORE UPDATE ON streamers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_viewers_modtime
BEFORE UPDATE ON viewers
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
