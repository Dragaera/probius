DROP TABLE IF EXISTS sc2replaystats_users;
DROP TABLE IF EXISTS discord_users;

CREATE TABLE discord_users (
	id           SERIAL      PRIMARY KEY,
  /* Discord user IDs are 64-bit, but discordgo exposes them as a string.
   * 64-bit integers have <= 20 digits. */
	discord_id   VARCHAR(20) UNIQUE NOT NULL,
	created_at   TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP
);



CREATE TABLE sc2replaystats_users (
	id              SERIAL       PRIMARY KEY,
	discord_user_id INT          UNIQUE NOT NULL REFERENCES discord_users(id) ON DELETE CASCADE,
	api_key         VARCHAR(128)        NOT NULL,
	last_replay_id  INT                 NOT NULL DEFAULT 0,
	created_at      TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

