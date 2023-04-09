# Is server alive?

A simple Discord bot to tell you when your server-chan is up, and where you can find her (server public IP)

Use Discord slash command `/ip` to ask for server public IP.

## Using Docker Compose

Easy enough
- clone the project
- cd into the folder
- add your bot token to `docker-compose.yml`
- run `docker-compose up -d`
- done

## Build it manually

Easy,
- `go build`