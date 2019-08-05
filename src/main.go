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

func handleRemindMe(bot *tb.Bot) func(*tb.Message) {
	return func(message *tb.Message) {
		messageContent := strings.TrimSpace(message.Text[len("/remind"):])
		remindMsg, err := decodeReminderMessage(messageContent)

		if err != nil {
			errorMessage := fmt.Sprintf(
				"Could not parse time for message \"%v\"",
				messageContent)
			log.Printf("Error parsing message %v: %v", messageContent, err.Error())

			trySendMessage(bot, message.Sender, errorMessage)
			return
		}

		log.Printf("Scheduled sending %v to %v", remindMsg, message.Sender)
		duration := time.Until(remindMsg.time)

		timer1 := time.NewTimer(duration)
		<-timer1.C

		trySendMessage(bot, message.Sender, remindMsg.content)
	}
}

func trySendMessage(bot *tb.Bot, sender *tb.User, message string) {
	log.Printf("Sending \"%v\" to \"%v\"", message, sender.ID)
	msg, err := bot.Send(sender, message)

	if err != nil {
		log.Printf("Failed to send message \"%v\" to \"%v\".", message, sender.ID)
	}

	log.Printf("Message \"%v\" sent to \"%v\"", msg.Text, sender.ID)
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

type remindMessage struct {
	time    time.Time
	content string
}

func (remindMsg *remindMessage) String() string {
	return fmt.Sprintf("%v at %v", remindMsg.content, remindMsg.time)
}
