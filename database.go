package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Link struct {
	Prefix string
	Suffix string
}

func (link *Link) Slide() {
	prefixWords := strings.Split(link.Prefix, " ")
	var newPrefix []string
	newPrefix = append(newPrefix, prefixWords[1])
	newPrefix = append(newPrefix, link.Suffix)
	link.Prefix = strings.Join(newPrefix, " ")
}

func (bot *Bot) dialDB() {
	connstring := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", bot.DBUser, bot.DBPassword, bot.DBName)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec("CREATE TABLE markov (prefix varchar(32), suffix varchar(32))")
	if err != nil {
		log.Println(err.Error())
	}

	bot.DB = db
}

func (bot *Bot) generateLine() string {
	var link Link
	var line []string

	query, _ := bot.DB.Prepare("SELECT prefix, suffix FROM markov ORDER BY RANDOM()")
	query.QueryRow().Scan(&link.Prefix, &link.Suffix)
	line = append(line, link.Prefix)
	line = append(line, link.Suffix)

	for i := 0; i <= bot.MaxLineLength; i++ {
		link.Slide()
		query, err := bot.DB.Prepare("SELECT suffix FROM markov WHERE prefix = $1 ORDER BY RANDOM()")
		if err != nil {
			break
		}

		err = query.QueryRow(link.Prefix).Scan(&link.Suffix)
		if err != nil {
			break
		}

		line = append(line, link.Suffix)
	}

	return strings.Join(line, " ")
}

func (bot *Bot) addLink(link Link) {
	bot.Mutex.Lock()
	_, err := bot.DB.Exec("INSERT INTO markov (prefix, suffix) VALUES ($1, $2)", link.Prefix, link.Suffix)
	if err != nil {
		log.Println(err.Error())
	}
	bot.Mutex.Unlock()
}
