# UWA server checker

Just another golang project

# Development

## Environment

- copy the provided example `cp env_example.sh development_env.sh`

## Telegram Bot

- The bot can only be used on a single server. Turn of the bot for testing.
- The bot only recognize users that are defined on the env.

### Commands

- `/containers` list out all the running containers
- `/get` get the opening TCP port, for tunneling/ssh

## Docker socket, for development

- The program default to `unix:///var/run/docker.sock`. To enable remote, setup `socat` and tunnel it via ssh to the server.
- Run `socat`
```
sudo socat \
  TCP-LISTEN:23750,bind=127.0.0.1,reuseaddr,fork \
  UNIX-CONNECT:/var/run/docker.sock
```
- Open new terminal and expose it through ssh
```
ssh -N -L 23750:localhost:23750 user@server
```
- Then set the env `DOCKER_HOST` before running on development