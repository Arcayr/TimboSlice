package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func (bot *Bot) run() {
	err := bot.Conn.Connect()
	if err != nil {
		panic(err.Error())
	}
}

func (bot *Bot) handleConnect(conn *irc.Conn, line *irc.Line) {
	ghostMsg := fmt.Sprintf("GHOST %s %s", bot.Nickname, bot.NickservPass)
	identMsg := fmt.Sprintf("IDENTIFY %s %s", bot.Nickname, bot.NickservPass)

	// Some IRC networks get upset if you flood NICKSERV without a delay.
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

func (bot *Bot) handleDisconnect(conn *irc.Conn, line *irc.Line) {
	log.Println("Disconnected from server. Retrying...")
	err := bot.Conn.Connect()
	if err != nil {
		log.Println("No dice. Shutting down.")
	}
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

func (bot *Bot) addLinks(links []Link) {
	for _, link := range links {
		go bot.addLink(link)
	}
}

func (bot *Bot) processLine(line string) {
	links := splitLine(line)
	bot.addLinks(links)
}

func (bot *Bot) processTrainingFile() {
	file, err := os.Open(bot.TrainFile)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bot.processLine(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (bot *Bot) handlePrivmsg(conn *irc.Conn, line *irc.Line) {
	// Check the PRIVMSG wasn't a query.
	if line.Public() == false {
		return
	}

	// Check if the user is annoying. Break out of this if so.
	for _, nick := range bot.Ignored {
		if line.Nick == nick {
			return
		}
	}

	// Check if we're being told to stop posting.
	if strings.Contains(line.Text(), strings.ToLower(bot.Nickname)) && strings.Contains(line.Text(), strings.ToLower(bot.Trigger)) {
		for _, nick := range bot.Moderators {
			if line.Nick == nick {
				bot.Posting = !bot.Posting
				log.Printf("Toggling posting. Now: %t", bot.Posting)
				return
			}
		}
	}

	// Only add lines from people who aren't annoying.
	bot.processLine(line.Text())

	// Check we're allowed to post. Otherwise, quit out.
		if bot.Posting == false {
		return
	}

	// If the bot was pinged, override the randomness and delay. This gets skipped if the line contained the nick.
	if strings.Contains(strings.ToLower(line.Text()), strings.ToLower(bot.Nickname)) == false {
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

	// Fire on all cylinders, we're good to go.
	generatedLine := bot.generateLine()
	bot.Conn.Privmsg(line.Target(), generatedLine)
	bot.LastTime = time.Now().Unix()
}

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

	bot.Conn.HandleFunc(irc.CONNECTED, bot.handleConnect)
	bot.Conn.HandleFunc(irc.DISCONNECTED, bot.handleDisconnect)
	bot.Conn.HandleFunc(irc.PRIVMSG, bot.handlePrivmsg)
}
