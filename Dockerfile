FROM golang:alpine

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN apk add --no-cache git
RUN go get github.com/karrick/tparse
RUN go get gopkg.in/tucnak/telebot.v2

RUN go build -o main /app/src/main.go
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./main"]
