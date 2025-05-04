-- name: CreateStreamer :one
INSERT INTO streamers (twitch_id, username, verified, access_token, refresh_token, profile_image_url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateViewer :one
INSERT INTO viewers (twitch_id, username, registered_in)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetViewerByID :one
SELECT * FROM viewers WHERE twitch_id = $1;

-- name: GetStreamerByID :one
SELECT * FROM streamers WHERE twitch_id = $1;

-- name: GetAllStreamers :many
SELECT username, twitch_id, profile_image_url FROM streamers WHERE verified = TRUE;

-- name: GetAllStreamersWithTokens :many
SELECT * FROM streamers;

-- name: UpdateStreamerTokens :one
UPDATE streamers 
SET access_token = $2, 
    refresh_token = $3
WHERE twitch_id = $1
RETURNING *;

-- name: CreateReward :one
INSERT INTO rewards (reward_id, streamer_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetRewardsByStreamer :many
SELECT * FROM rewards WHERE streamer_id = $1 ORDER BY created_at DESC;

-- name: DeleteRewardsByStreamerID :exec
DELETE FROM rewards
WHERE streamer_id = $1;

-- name: CreateRedemption :one
INSERT INTO redemptions (message_id, streamer_id, viewer_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetViewerLeaderboard :many
SELECT
    v.username, -- Retrieve the username from the viewers table
    COUNT(r.*) AS total_redemptions
FROM
    redemptions r
JOIN
    viewers v ON r.viewer_id = v.twitch_id -- Join the tables based on viewer_id
GROUP BY
    r.viewer_id, v.username -- Group by both viewer_id and username
ORDER BY
    total_redemptions DESC;

-- name: GetRecentRedemptionsWithUsernames :many
SELECT
    r.message_id,
    s.username AS streamer_username, -- Get streamer username from streamers table
    v.username AS viewer_username,   -- Get viewer username from viewers table
    r.redeemed_at
FROM
    redemptions r
JOIN
    streamers s ON r.streamer_id = s.twitch_id
JOIN
    viewers v ON r.viewer_id = v.twitch_id
ORDER BY
    r.redeemed_at DESC
LIMIT $1;

-- name: GetTotalRedemptionsCount :one
SELECT COUNT(*) AS total_redemptions
FROM redemptions;

-- name: GetTotalParticipantsCount :one
SELECT COUNT(*) AS total_participants
FROM viewers;

