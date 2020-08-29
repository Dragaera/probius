# Probius

This repository hosts 'Probius', a Discord bot interacting with the API of
https://sc2replaystats.com - a replay-processing site for StarCraft 2.

## Usage

This section serves to document usage of the bot as a enduser.

### Authorizing access to API

To start with, you must grant the bot access to SC2ReplayStats' API on your
behalf. To do so:

- Sign in to your profile on https://sc2replaystats.com
- Select 'My Account' -> 'Settings' in the top right
- Select 'API Access'
- Copy the key shown below 'Authorization Key'
- **In a direct message** with the bot, use the following command: `!auth PASTE_YOUR_KEY_HERE`

### Basic usage

You can now use the `!last` command in any channel the bot is in, to make it
show your most-recently uploaded replay.

### Managing automatic notifications

You can now optionally configure the bot to automatically look for new replays,
and post a message in a channel of your choice when a new replay is detected.
To do so:

- Use the `!subscribe` command in the channel it should post messages to
- Use the `!unsubscribe` command to make it stop posting messages
- Use the `!subscriptions` command to list channels where it will automatically post your replays to

## Configuration

Configuration is done exclusively via environment variables, documented in the
following section.

### Discord configuration

| Environment variable | Default value | Comment           |
| -------------------- | ------------- | ----------------- |
| `DISCORD_CLIENT_ID`  | -             | Discord client ID |
| `DISCORD_TOKEN`      | -             | Discord token     |

### DB configuration

Only Postgres is supported.

| Environment variable | Default value | Comment                                                                                                              |
| -------------------- | ------------- | -------------------------------------------------------------------------------------------------------------------- |
| `DB_USER`            | -             | Database username                                                                                                    |
| `DB_PASSWORD`        | -             | Database password                                                                                                    |
| `DB_HOST`            | 127.0.0.1     | Database host                                                                                                        |
| `DB_PORT`            | 5432          | Database port                                                                                                        |
| `DB_DATABASE`        | probius       | Database name                                                                                                        |
| `DB_SSL_MODE`        | disable       | Postgres SSL mode as per [the documentation](https://www.postgresql.org/docs/12/libpq-ssl.html#LIBPQ-SSL-PROTECTION) |
| `DB_LOG_SQL`         | false         | Set to `true` to enable logging of all SQL statements to STDOUT                                                      |

### Redis configuration

| Environment variable    | Default value | Comment    |
| ----------------------- | ------------- | ---------- |
| `REDIS_HOST`            | 127.0.0.1     | Redis host |
| `REDIS_PORT`            | 5432          | Redis port |

### Background worker configuration

| Environment variable    | Default value | Comment                                     |
| ----------------------- | ------------- | ------------------------------------------- |
| `WORKER_CONCURRENCY`    | 5             | Number of background workers to spawn       |
| `WORKER_NAMESPACE`      | probius       | Redis key prefix used for worker management |

### SC2ReplayStats configuration

| Environment variable                  | Default value | Comment                                                                                    |
| ------------------------------------- | ------------- | ------------------------------------------------------------------------------------------ |
| `SC2_REPLAY_STATS_UPDATE_INTERVAL`    | 300           | Update interval in seconds in which to check whether there are new replays per player      |
| `SC2_REPLAY_STATS_LOCK_TTL`           | 900           | Duration in seconds after which to consider update job to have silently died               |
| `SC2_REPLAY_STATS_RATE_LIMIT_AVERAGE` | 1             | Amount of average requests per second to use for rate-limiting towards SC2ReplayStats' API |
| `SC2_REPLAY_STATS_RATE_LIMIT_BURST`   | 2             | Amount of burst requests to allow for rate-limiting towards SC2ReplayStats' API            |
