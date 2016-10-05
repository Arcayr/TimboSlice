package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

var config Configuration

type Configuration struct {
	Memory       string
	Server       string
	Port         int
	TLS          bool
	Nickname     string
	Username     string
	RealName     string
	NickServPass string
	Channels     []string
	Interval     int64
	Delay        int
	ChancePool   int
	LineLength   int
	IgnoreList   []string
}

func (config *Configuration) LoadFromFile(path string) {
	if path == "" {
		path = "./config.json"
	}

	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal(err.Error())
		panic(err)
	}

	log.Println("Configuration loaded!")

	TLSPlus := ""

	if config.TLS {
		TLSPlus = "+"
	}

	log.Printf("%s!%s (%s) @ %s:%s%d\n", config.Nickname, config.Username, config.RealName, config.Server, TLSPlus, config.Port)
	log.Printf("Responding, at most, once every %d seconds with a %d second delay, and a one in %d chance of responding.", config.Interval, config.Delay, config.ChancePool)
	log.Printf("Joining: %s", strings.Join(config.Channels, ", "))
	log.Printf("Ignoring: %s", strings.Join(config.IgnoreList, ", "))
}
