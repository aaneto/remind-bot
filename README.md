# Remind Bot
A telegram bot for message reminding.

[![Build Status](https://travis-ci.org/aaneto/remind-bot.svg?branch=master)](https://travis-ci.org/aaneto/remind-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/aaneto/remind-bot)](https://goreportcard.com/report/github.com/aaneto/remind-bot)
[![codecov](https://codecov.io/gh/aaneto/remind-bot/branch/master/graph/badge.svg)](https://codecov.io/gh/aaneto/remind-bot)

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
