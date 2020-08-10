package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func (bot *Bot) HandleMessage(inMessage Message, id string) {
	env, err := LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(fmt.Sprintf("inMessage : %+v", inMessage))
	logrus.Info(fmt.Sprintf("id : %v", id))

	var outMessage Message
	if inMessage.Text == "" {
		outMessage = Message{
			Text: "Sorry, I am not yet trained to answer those! :( Are there anything else that you want to talk about?",
		}
	} else {
		query := inMessage.Text
		previousQuery, err := bot.GetLastMessage(id)

		if err != nil {
			logrus.Info(fmt.Sprintf("No Previous Conversation detected with %v", id))
		} else {
			logrus.Info(fmt.Sprintf("Previous Conversation detected with %v", id))
			bot.Add(previousQuery, query)
		}

		resp, err := bot.Get(query)
		if err != nil {
			logrus.Fatal(err)
		}
		outMessage = Message{
			Text: resp,
		}
		if err != nil {
			logrus.Fatal(err)
		}

		bot.UpdateLastMessage(id, resp)
	}

	logrus.Info(fmt.Sprintf("outMessage : %+v", outMessage))

	outResponse := Response{
		Recipient: User{ID: id},
		Message:   outMessage,
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

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(fmt.Sprintf("Response : %v", string(bodyBytes)))

}

func (bot *Bot) VerificationHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("VerificationHandler is called")
	env, err := LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	challengeQuery := r.URL.Query().Get("hub.challenge")
	modeQuery := r.URL.Query().Get("hub.mode")
	verifyTokenQuery := r.URL.Query().Get("hub.verify_token")

	verifyToken := env.VerifyToken

	logrus.Info(fmt.Sprintf("Challenge Query : %v", challengeQuery))
	logrus.Info(fmt.Sprintf("modeQuery : %v", modeQuery))
	logrus.Info(fmt.Sprintf("verifyTokenQuery : %v", verifyTokenQuery))

	if modeQuery == "subscribe" && verifyTokenQuery == verifyToken {
		w.WriteHeader(200)
		_, err := w.Write([]byte(challengeQuery))
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Info("Verification is Successful")
	} else {
		w.WriteHeader(404)
		_, err := w.Write([]byte("Unauthorized Token"))
		logrus.Info("Verification has failed.")
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func (bot *Bot) TestHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("TestHandler is called")
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(fmt.Sprintf("Request Received : %v", string(requestDump)))

	w.WriteHeader(200)
	_, err = w.Write([]byte("Hello World"))
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("TestHandler is Successful")
}

func (bot *Bot) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("CallbackHandler is called")
	var callback Callback
	err := json.NewDecoder(r.Body).Decode(&callback)
	if err != nil {
		logrus.Fatal(err)
	}

	if callback.Object == "page" {
		for _, entry := range callback.Entries {
			for _, messaging := range entry.Messagings {
				if messaging.Message != (Message{}) {
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
			logrus.Info("Callback is Successful")
		}
	} else {
		w.WriteHeader(400)
		_, err := w.Write([]byte("callback.Object must be page"))
		if err != nil {
			logrus.Fatal(err)
		}
	}
}
