package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/karrick/tparse"
	tb "gopkg.in/tucnak/telebot.v2"
)

type remindMessage struct {
	time    time.Time
	content string
}

func (remindMsg *remindMessage) String() string {
	return fmt.Sprintf("%v at %v", remindMsg.content, remindMsg.time)
}

func decodeReminderMessage(message string) (*remindMessage, error) {
	splitResponse := strings.Fields(message)

	timeString := splitResponse[0]
	messageString := strings.Join(splitResponse[1:], " ")

	actual, err := tparse.ParseNow(time.RFC3339, timeString)

	if err != nil {
		return nil, err
	}

	remindMsg := remindMessage{time: actual, content: messageString}

	return &remindMsg, nil
}

func handleRemindMe(bot *tb.Bot) func(*tb.Message) {
	return func(message *tb.Message) {
		messageContent := strings.TrimSpace(message.Text[len("/remind"):])
		remindMsg, err := decodeReminderMessage(messageContent)

		if err != nil {
			errorMessage := fmt.Sprintf(
				"Could not parse duration %v: %v",
				remindMsg.time,
				err)

			log.Printf("Sending to Bot %v: %v", message.Sender, errorMessage)
			bot.Send(
				message.Sender,
				errorMessage)
		}

		log.Printf("Scheduled to Send to Bot %v: %v", message.Sender, remindMsg)
		duration := time.Until(remindMsg.time)

		timer1 := time.NewTimer(duration)
		<-timer1.C

		log.Printf("Sending to Bot %v: %v", message.Sender, remindMsg.content)
		bot.Send(
			message.Sender,
			remindMsg.content)
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/remind", handleRemindMe(b))
	b.Start()
}
