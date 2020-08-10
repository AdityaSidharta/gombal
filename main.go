package main

import (
	"fmt"
	"github.com/adityasidharta/gombal/pkg"
)

func main() {
	c, err := pkg.LoadConfig(pkg.ConfigPath)
	if err != nil {
		panic(err)
	}
	bot, err := pkg.NewBot(c.Strategy, pkg.DataPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v, \n", bot)
	err = bot.Save(pkg.DataPath)
	fmt.Printf("%v, \n", err)

	e, err := pkg.LoadEnv()
	if err!= nil {
		panic(err)
	}
	fmt.Printf("%+v" , e)
}
