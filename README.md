# Sunbot
Stateless Discord bot originally made for the [Cuddle Puddle](https://floof.zone/discord) discord server.

## Setup

Sunbot is intended to be used with one instance per Discord server/guild. You CAN connect it to separate servers, however the databases will be merged (if you choose to use one).

Sunbot is entirely stateless, however it depends on several environment variables to be set.
The `.env.sample` file should contain up-to-date listing in case this readme is neglected (it's possible).

`*` - required

* `DISCORD_AUTH_TOKEN`* - Discord bot API token

* `COMMAND_PREFIX` - Prefix used by users to execute commands (default `.`)

* `DEBUG_OUTPUT` - Verbose debug output (default `true`)

* `SILLY_COMMANDS` - Enable the silly commands which do not use the command prefix (default `true`)

* `REDIS_URL` - Redis database URL:PORT (leave blank to disable)

* `REDIS_PASSWORD` - Redis database password (leave blank if none)

Dockerfile and launch script are included, which will always pull the latest commit on launch. A "stable" release will exist eventually.

## Commands

Use `.help` and `.help [verb]` for an up-to-date list. All commands and respective functionality and help info are defined in `commands.go`.
