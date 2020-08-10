package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adityasidharta/gombal/pkg/facebook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func (bot *Bot) HandleMessage(inMessage facebook.Message, id string) {
	env, err := LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	var r facebook.Message
	if inMessage.Text == "" {
		r = facebook.Message{
			Text: "Sorry, I am not yet trained to answer those! :( Are there anything else that you want to talk about?",
		}
	} else {
		query := inMessage.Text
		resp, err := bot.Get(query)
		if err != nil {
			logrus.Fatal(err)
		}
		r = facebook.Message{
			Text: resp,
		}
		if err != nil {
			logrus.Fatal(err)
		}
	}
	outResponse := facebook.Response{
		Recipient: facebook.User{ID: id},
		Message:   r,
	}

	requestBody, err := json.Marshal(outResponse)
	if err != nil {
		logrus.Fatal(err)
	}

	url := fmt.Sprintf(FacebookApi, env.PageAccessToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logrus.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (bot *Bot) VerificationHandler(w http.ResponseWriter, r *http.Request) {
	env, err := LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	challengeQuery := r.URL.Query().Get("hub.challenge")
	modeQuery := r.URL.Query().Get("hub.mode")
	verifyTokenQuery := r.URL.Query().Get("hub.verify_token")

	verifyToken := env.VerifyToken

	if modeQuery != "subscribe" && verifyTokenQuery == verifyToken {
		w.WriteHeader(200)
		_, err := w.Write([]byte(challengeQuery))
		if err != nil {
			logrus.Fatal(err)
		}
	} else {
		w.WriteHeader(401)
		_, err := w.Write([]byte("Unauthorized Token"))
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func (bot *Bot) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var callback facebook.Callback
	err := json.NewDecoder(r.Body).Decode(&callback)
	if err != nil {
		logrus.Fatal(err)
	}

	if callback.Object == "page" {
		for _, entry := range callback.Entries {
			for _, messaging := range entry.Messagings {
				if messaging.Message != (facebook.Message{}) {
					bot.HandleMessage(messaging.Message, messaging.Sender.ID)
				} else {
					logrus.Fatal(emptyMessagingError)
				}
			}
			w.WriteHeader(200)
			_, err := w.Write([]byte("Message received successfully"))
			if err != nil {
				logrus.Fatal(err)
			}
		}
	} else {
		w.WriteHeader(400)
		_, err := w.Write([]byte("callback.Object must be page"))
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
