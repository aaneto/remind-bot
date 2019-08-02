FROM golang:alpine as builder

RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN apk add git --no-cache

RUN go get github.com/karrick/tparse
RUN go get gopkg.in/tucnak/telebot.v2
RUN go build -o main /build/src/main.go


FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates 2>/dev/null || true

RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
