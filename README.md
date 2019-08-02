# Remind Bot
A telegram bot for message reminding.

## Running

```bash
# Install packages
go get github.com/karrick/tparse
go get gopkg.in/tucnak/telebot.v2

# Run binary
go run src/main.go
```

## Runing with docker

```bash
# Build locally
docker build . -t adilsinho/remind-bot
# Pull from Dockerhub
docker pull adilsinho/remind-bot

docker run -e TELEGRAM_BOT_TOKEN='YOUR_BOT_TOKEN' adilsinho/remind-bot
```
