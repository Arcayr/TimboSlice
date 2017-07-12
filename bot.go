package main

import (
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"math/rand"
	"strings"
	"time"
)

func (bot *Bot) dial() {
	bot.LastTime = 0

	config := irc.NewConfig(bot.Nickname, bot.Username, bot.RealName)
	if bot.TLS {
		config.SSL = true
		config.SSLConfig = &tls.Config{ServerName: bot.ServerName}
	}

	config.Server = bot.ServerName

	bot.Conn = irc.Client(config)
	bot.Conn.EnableStateTracking()

	bot.Conn.HandleFunc(irc.CONNECTED, bot.HandleConnect)
	bot.Conn.HandleFunc(irc.DISCONNECTED, bot.HandleDisconnect)
	bot.Conn.HandleFunc(irc.PRIVMSG, bot.HandlePrivmsg)
}

func (bot *Bot) run() {
	err := bot.Conn.Connect()
	if err != nil {
		panic(err.Error())
	}
}

func (bot *Bot) HandleConnect(conn *irc.Conn, line *irc.Line) {
	ghostMsg := fmt.Sprintf("GHOST %s %s", bot.Nickname, bot.NickservPass)
	identMsg := fmt.Sprintf("IDENTIFY %s %s", bot.Nickname, bot.NickservPass)

	time.Sleep(3 * time.Second)
	conn.Privmsg("NICKSERV", ghostMsg)
	time.Sleep(3 * time.Second)
	conn.Privmsg("NICKSERV", identMsg)
	time.Sleep(5 * time.Second)

	for _, channel := range bot.Channels {
		log.Println("Joining ", channel)
		conn.Join(channel)
	}
}

func (bot *Bot) HandlePrivmsg(conn *irc.Conn, line *irc.Line) {
	if line.Public() == false { // Check the PRIVMSG wasn't a query.
		return
	}

	links := splitLine(line.Text())

	// Check if the user is annoying. Break out of this if so.
	for _, nick := range bot.Ignored {
		if line.Nick == nick {
			return
		}
	}

	// Check if we're being told to stop posting.
	if line.Text() == bot.Trigger {
		for _, nick := range bot.Moderators {
			if line.Nick == nick {
				bot.Posting = !bot.Posting
				log.Printf("Toggling posting. Now: %t", bot.Posting)
				return
			}
		}
	}

	// Only add lines from people who aren't annoying.
	bot.addLinks(links)

	// Check we're allowed to post
	if bot.Posting == false {
		return
	}


	if strings.HasPrefix(strings.ToLower(line.Text()), strings.ToLower(bot.Nickname)) == false { // If he was pinged, override the randomness and delay.
		// Check it hasn't been too short a time since the previous post.
		if bot.LastTime+bot.MinGap >= time.Now().Unix() {
			return
		}

		// XKCD 221
		r := rand.Intn(100)
		if r > bot.Chance {
			return
		}		
	}

	// Generate a line.
	generatedLine := bot.generateLine()
	bot.Conn.Privmsg(line.Target(), generatedLine)
	bot.LastTime = time.Now().Unix()
}

func (bot *Bot) HandleDisconnect(conn *irc.Conn, line *irc.Line) {
	log.Println("Disconnected from server. Shutting down...")
	bot.DB.Close()
	quit <- true
}

func splitLine(line string) []Link {
	var links []Link
	words := strings.Fields(line)

	for i := 0; i <= len(words)-3; i++ {
		prefix := strings.Join(words[i:i+2], " ")
		suffix := words[i+2]

		link := Link{Prefix: prefix, Suffix: suffix}
		links = append(links, link)
	}

	return links
}
