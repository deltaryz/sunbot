# Sunbot

![](https://orig00.deviantart.net/fdd0/f/2017/183/7/9/untitled_by_hiccupsdoesart-dbeutpr.png)

Stateless Discord bot originally made for the [Cuddle Puddle](https://floof.zone/discord) discord server.

Art by [HiccupsDoesArt](https://twitter.com/HiccupsDoesArt)

## Setup

Sunbot is intended to be used with one instance per Discord server/guild. You CAN connect it to separate servers, however the databases will be merged (if you choose to use one).

Sunbot is entirely stateless, however it depends on several environment variables to be set.
The `.env.sample` file should contain up-to-date listing in case this readme is neglected (it's possible).

`*` - required

* `DISCORD_AUTH_TOKEN`* - Discord bot API token

* `COMMAND_PREFIX` - Prefix used by users to execute commands (default `.`)

* `DEBUG_OUTPUT` - Verbose debug output (default `true`)

* `SILLY_COMMANDS` - Enable the silly commands which do not use the command prefix (default `true`)

* `DERPIBOORU_API_KEY` - API key for Derpibooru queries (leave blank if none)

Dockerfile and launch script are included here as well as `techniponi/sunbot` on [Docker Hub](https://hub.docker.com/r/techniponi/sunbot/), which will always pull the latest commit on launch. A "stable" release will exist eventually.

Note: the Redis database functionality has been *disabled* until further notice. Focus will be directed at the stateless functionality for now.

## Commands

Use `.help` and `.help [verb]` for an up-to-date list. All commands and help info are defined in `commands.go`.
