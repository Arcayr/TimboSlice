package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

func (bot *Bot) loadConfiguration(filepath string) {

	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		flag.Usage()
		log.Fatal(err.Error())
	}

	err = json.Unmarshal(configFile, &bot)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	bot.Trigger = fmt.Sprintf("%s: %s", bot.Nickname, bot.Trigger)

	log.Println("Configuration loaded!")
}
