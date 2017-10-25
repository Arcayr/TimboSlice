package main

import (
	"database/sql"
	"flag"
	"sync"

	irc "github.com/fluffle/goirc/client"
)

type Bot struct {
	Nickname      string   `json:"nickname"`
	Username      string   `json:"username"`
	RealName      string   `json:"real_name"`
	NickservPass  string   `json:"nickserv_pass"`
	ServerName    string   `json:"server_name"`
	Port          int      `json:"port"`
	TLS           bool     `json:"tls"`
	Channels      []string `json:"channels"`
	Ignored       []string `json:"ignored"`
	Moderators    []string `json:"moderators"`
	MaxLineLength int      `json:"max_line_length"`
	Chance        int      `json:"chance"`
	Posting       bool     `json:"posting"`
	Trigger       string   `json:"trigger"`
	MinGap        int64    `json:"min_gap_seconds"`
	DBName        string   `json:"db_name"`
	DBUser        string   `json:"db_username"`
	DBPassword    string   `json:"db_password"`
	Training      bool     `json:"training"`
	TrainFile     string   `json:"trainfile"`
	DB            *sql.DB
	Conn          *irc.Conn
	LastTime      int64
	Mutex         sync.RWMutex
}

func main() {
	var bot Bot

	var configFilePath = flag.String("configFilePath", "timhortons.json", "The path to the JSON config file.")
	flag.Parse()

	bot.loadConfiguration(*configFilePath)
	bot.dialDB()
	if bot.Training {
		bot.processTrainingFile()
	}

	bot.dial()
	bot.run()

	select {}
}
