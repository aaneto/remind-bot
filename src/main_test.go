package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var mockGetMe = `
{
	"ok": true,
	"result": {"id": 0, "is_bot": true, "first_name": "Bot"}
}
`

var mockGetUpdates = `
{
	"ok": true,
	"result": [%v]
}
`

var mockUpdate = `
{
	"update_id": %v,
	"message": {
		"message_id": %v,
		"date": 0,
		"from": {
			"id": 1,
			"is_bot": false,
			"first_name": "User"
		},
		"chat": {
			"id": 1,
			"type": "group"
		},
		"text": "%v"
	}
}
`

func TestCommandParsing(t *testing.T) {
	now := time.Now()

	expectedDurationMillis := 50
	expectedDuration := time.Duration(expectedDurationMillis) * time.Millisecond
	expectedTime := now.Add(expectedDuration)
	acceptedDelay := time.Millisecond

	messageText := "MessageContent Foobar"
	message := fmt.Sprintf("now+%vms %v", expectedDurationMillis, messageText)
	remindMessage, err := decodeReminderMessage(message)

	if err != nil {
		t.Errorf("Could not create remindMessage from %v", message)
	}

	if !timeAlmostEqual(expectedTime, remindMessage.time, acceptedDelay) {
		t.Error("Parsed duration is different than expected duration!")
	}

	if remindMessage.content != "MessageContent Foobar" {
		t.Error("Content not parsed correctly.")
	}
}

func TestInvalidCommandParsing(t *testing.T) {
	_, err := decodeReminderMessage("one-second my message")

	if err == nil {
		t.Error("Message should not be parsable.")
	}
}

func TestRemindMessageStringer(t *testing.T) {
	utc, _ := time.LoadLocation("UTC")
	remindMsg := remindMessage{time: time.Unix(0, 0).In(utc), content: "My Custom Message"}
	expectedMessage := "My Custom Message at 1970-01-01 00:00:00 +0000 UTC"

	if remindMsg.String() != expectedMessage {
		t.Errorf("remindMsg String() should be equal to: %v, got %v instead", expectedMessage, remindMsg.String())
	}
}

func TestInvalidCommandHandle(t *testing.T) {
	messageContent := "Content of the message"
	message := fmt.Sprintf("/remind inammoment %v", messageContent)

	server := startFakeServer(message)
	defer server.Close()

	bot, _ := tb.NewBot(tb.Settings{
		Token:  "DUMMY TOKEN",
		URL:    server.getURL(),
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	bot.Handle("/remind", handleRemindMe(bot))

	// Text needs to be escaped inside JSON
	expectedContent := `Could not parse time for message \"inammoment Content of the message\"`
	expectedMessage := newMessage(expectedContent)
	go stopBotAfterDuration(time.Second, bot)

	bot.Start()

	if expectedMessage != *server.sentMessage {
		t.Errorf(
			"Wrong message received, expected %#v got %#v",
			expectedMessage,
			server.sentMessage)
	}
}

func TestCommandHandle(t *testing.T) {
	messageContent := "Content of the message"
	messageDelaySeconds := 1
	messageDelay := time.Duration(messageDelaySeconds) * time.Second
	message := fmt.Sprintf("/remind now+%vs %v", messageDelaySeconds, messageContent)
	server := startFakeServer(message)
	defer server.Close()

	bot, _ := tb.NewBot(tb.Settings{
		Token:  "DUMMY TOKEN",
		URL:    server.getURL(),
		Poller: &tb.LongPoller{Timeout: 1 * time.Second},
	})

	bot.Handle("/remind", handleRemindMe(bot))
	expectedMessage := newMessage(messageContent)

	acceptedDelay := messageDelay + time.Second
	go stopBotAfterDuration(acceptedDelay, bot)
	bot.Start()

	expectedArrivalTime := server.receivedMessageTime.Add(messageDelay)

	if !timeAlmostEqual(*server.sentMessageTime, expectedArrivalTime, acceptedDelay) {
		t.Error("Parsed time was not respected.")
	}

	if expectedMessage != *server.sentMessage {
		t.Errorf("Wrong message received, expected %#v got %#v", expectedMessage, server.sentMessage)
	}

}

type testServerData struct {
	sentMessage         *string
	receivedMessageTime *time.Time
	sentMessageTime     *time.Time
	server              *httptest.Server
}

func (serverData testServerData) Close() {
	serverData.server.Close()
}

func (serverData testServerData) getURL() string {
	return serverData.server.URL
}

func newMessage(text string) string {
	return fmt.Sprintf("{\"chat_id\":\"1\",\"text\":\"%v\"}\n", text)
}

// Returns true if the distance between two time.Time objects are within a tolerance.
func timeAlmostEqual(lhand time.Time, rhand time.Time, tolerance time.Duration) bool {
	var withinTolerance bool
	if rhand.After(lhand) {
		// lhand < rhand < lhand + tolerance
		withinTolerance = rhand.Before(lhand.Add(tolerance))
	} else {
		// rhand < lhand < rhand + tolerance
		withinTolerance = lhand.Before(rhand.Add(tolerance))
	}

	return withinTolerance
}

func startFakeServer(responseMessage string) testServerData {
	nonce := 1
	msgReceived := false
	var sentTime = new(time.Time)
	var receivedTime = new(time.Time)
	var sentMessage = new(string)

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(200 * time.Millisecond)

			switch {
			case strings.HasSuffix(r.URL.Path, "getMe"):
				w.Header().Add("Content-Type", "application/json")
				_, err := w.Write([]byte(mockGetMe))
				if err != nil {
					panic("Could not read mockGetMe bytes!")
				}

			case strings.HasSuffix(r.URL.Path, "sendMessage"):
				*sentTime = time.Now()

				buffer, _ := ioutil.ReadAll(r.Body)
				*sentMessage = string(buffer)
				r.Body.Close()

				messageResponse := fmt.Sprintf(`{"ok": true, "result": %v}`, *sentMessage)
				_, err := w.Write([]byte(messageResponse))

				if err != nil {
					panic("Could not write messageResponse bytes!")
				}

			case strings.HasSuffix(r.URL.Path, "getUpdates"):
				w.Header().Add("Content-Type", "application/json")
				if !msgReceived {
					*receivedTime = time.Now()
					updateString := fmt.Sprintf(mockUpdate, nonce, nonce, responseMessage)
					responseString := fmt.Sprintf(mockGetUpdates, updateString)
					_, err := w.Write([]byte(responseString))

					if err != nil {
						panic("Could not read mockUpdate bytes!")
					}
					msgReceived = true
				} else {
					updateString := fmt.Sprintf(mockGetUpdates, "")
					_, err := w.Write([]byte(updateString))

					if err != nil {
						panic("Could not read empty mockUpdate bytes!")
					}
				}
			}
		}))

	return testServerData{
		server:              ts,
		sentMessageTime:     sentTime,
		sentMessage:         sentMessage,
		receivedMessageTime: receivedTime}
}

func stopBotAfterDuration(duration time.Duration, bot *tb.Bot) {
	time.Sleep(duration)
	bot.Stop()
}
