package main

import (
	"fmt"
	"github.com/adityasidharta/gombal/pkg"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	env, err := pkg.LoadEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	c, err := pkg.LoadConfig(pkg.ConfigPath)
	if err != nil {
		logrus.Fatal(err)
	}

	bot, err := pkg.NewBot(c.Strategy, pkg.DataPath)
	if err != nil {
		logrus.Fatal(err)
	}

	go bot.PeriodicSave(pkg.DataPath)

	r := mux.NewRouter()
	r.HandleFunc("/", bot.TestHandler).Methods("GET")
	r.HandleFunc("/webhook", bot.VerificationHandler).Methods("GET")
	r.HandleFunc("/webhook", bot.CallbackHandler).Methods("POST")
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", env.Port), r); err != nil {
		log.Fatal(err)
	}
}
