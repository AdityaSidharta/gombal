package main

import (
	"fmt"
	"github.com/adityasidharta/gombal/gombal"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	env, err := gombal.LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	c, err := gombal.LoadConfig(gombal.ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	bot, err := gombal.NewBot(c.Strategy, gombal.DataPath)
	if err != nil {
		logrus.Fatal(err)
	}

	go bot.PeriodicSave(gombal.DataPath)

	r := mux.NewRouter()
	r.HandleFunc("/", bot.TestHandler).Methods("GET")
	r.HandleFunc("/webhook", bot.VerificationHandler).Methods("GET")
	r.HandleFunc("/webhook", bot.CallbackHandler).Methods("POST")
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", env.Port), r); err != nil {
		log.Fatal(err)
	}
}
