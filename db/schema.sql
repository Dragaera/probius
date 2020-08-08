DROP TABLE IF EXISTS sc2replaystats_users;

DROP TABLE IF EXISTS discord_users;
DROP TABLE IF EXISTS discord_channels;

DROP TABLE IF EXISTS discord_guilds;

/* Discord's IDs for users / channels / guilds etc all follow the 'snowflake'
 * format, which is a 64-bit integer - albeit the API returns it as strings, to
* prevent integer overflow in consumers. */

CREATE TABLE discord_guilds (
	id         SERIAL PRIMARY KEY,
	discord_id BIGINT      UNIQUE NOT NULL,
	name       VARCHAR(100)       NOT NULL,
	owner_id   BIGINT             NOT NULL,
	created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP
);
/* Fake guild for direct messages, allowing a NOT NULL constraint for DM channels also */
INSERT INTO discord_guilds (discord_id, name, owner_id) VALUES (0, 'Direct Message', 0);

CREATE TABLE discord_channels (
	id         SERIAL      PRIMARY KEY,
	discord_id BIGINT      UNIQUE NOT NULL,
	name       VARCHAR(100)       NOT NULL,
	is_dm      BOOLEAN            NOT NULL DEFAULT FALSE,
	created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE discord_users (
	id            SERIAL      PRIMARY KEY,
	discord_id    BIGINT      UNIQUE NOT NULL,
	/* No length limit mentioned in API docs */
	name          VARCHAR(500)       NOT NULL,
	discriminator VARCHAR(4)         NOT NULL,
	created_at    TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sc2replaystats_users (
	id              SERIAL       PRIMARY KEY,
	discord_user_id INT          UNIQUE NOT NULL REFERENCES discord_users(id) ON DELETE CASCADE,
	api_key         VARCHAR(128)        NOT NULL,
	last_replay_id  INT                 NOT NULL DEFAULT 0,
	created_at      TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP
);

