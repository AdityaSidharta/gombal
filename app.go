package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)


func Verification(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from AdityaSidharta :)")
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Verification).Methods("GET")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.Fatal(err)
	}
}
