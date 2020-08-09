package main

import (
	"fmt"
	"github.com/adityasidharta/gombal/pkg"
)


func main() {
	bot := pkg.NewBot(pkg.MAXIMUM)
	bot.Add("Heya!", "Hello World")
	bot.Add("Heya!", "Hello World")
	bot.Add("Heya!", "Azzubilah")
	err := bot.RemoveQuery("Heya!")
	fmt.Printf("%v, \n", err)
	fmt.Printf("%+v, \n", bot)
}