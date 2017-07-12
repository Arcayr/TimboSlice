package main

import (
	"database/sql"
	irc "github.com/fluffle/goirc/client"
	"os"
	"sync"
)

var quit chan (bool)

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
	DB            *sql.DB
	Conn          *irc.Conn
	LastTime      int64
	Mutex         sync.RWMutex
}

func main() {
	var bot Bot

	configFilePath := os.Args[1]
	bot.loadConfiguration(configFilePath)
	bot.dialDB()
	bot.dial()
	bot.run()

	<-quit
}
