package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func (bot *Bot) loadConfiguration(filepath string) {
	if filepath == "" {
		filepath = "./timhortons.json"
	}

	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	err = json.Unmarshal(configFile, &bot)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	bot.Trigger = fmt.Sprintf("%s: %s", bot.Nickname, bot.Trigger)

	log.Println("Configuration loaded!")
}
