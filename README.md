# Sunbot
Stateless Discord bot originally made for the [Cuddle Puddle](https://floof.zone/discord) discord server.

## Setup

Sunbot is entirely stateless, however it depends on several environment variables to be set.
The `docker-env.txt` file should contain up-to-date listing in case this readme is neglected (it's possible).

* `DISCORD_AUTH_TOKEN` - Discord bot API token

* `COMMAND_PREFIX` - Prefix used by users to execute commands

* `DEBUG_OUTPUT` - Verbose debug output


* `SILLY_COMMANDS` - Enable the silly commands which do not use the command prefix

## Commands

Use `.help` and `.help [verb]` for an up-to-date list.