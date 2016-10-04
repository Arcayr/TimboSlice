package main

import (
	"crypto/tls"
	"math/rand"
	"fmt"
	"time"
	"strings"
	irc "github.com/fluffle/goirc/client"
)

var client *irc.Conn
var LastTime int64

func NewBot(config Configuration) {
	LastTime = time.Now().Unix()
	clientConfig := irc.NewConfig(config.Nickname, config.Username, config.RealName)
	if config.TLS {
		clientConfig.SSL = config.TLS
		clientConfig.SSLConfig = &tls.Config{ServerName: config.Server}
	}

	clientConfig.Server = config.Server

	client = irc.Client(clientConfig)
	client.EnableStateTracking()

	client.HandleFunc(irc.CONNECTED, ConnectHandler)
	client.HandleFunc(irc.DISCONNECTED, DisconnectHandler)
	client.HandleFunc(irc.PRIVMSG, PrivmsgHandler)
}

func ConnectHandler(conn *irc.Conn, line *irc.Line) {
	ghostMsg := fmt.Sprintf("GHOST %s %s", config.Nickname, config.NickServPass)
	identMsg := fmt.Sprintf("IDENTIFY %s %s", config.Nickname, config.NickServPass)

	conn.Privmsg("NICKSERV", ghostMsg)
	conn.Privmsg("NICKSERV", identMsg)

	for _, channel := range config.Channels {
		conn.Join(channel)
	}
}

func DisconnectHandler(conn *irc.Conn, line *irc.Line) {
	quit <- true
}

func PrivmsgHandler(conn *irc.Conn, line *irc.Line) {
	// Roll the dice : 0 <= n < config.RandPool
	dice := rand.Intn(config.ChancePool)
	if dice != 0 {
		return
	}

	// Make sure it's not still time to stop posting.
	if time.Now().Unix() < LastTime + config.Interval {
		return
	}

	// Check the last poster is supposed to be posting.
	for _, n := range config.IgnoreList {
		if n == line.Nick {
			return
		}
	}

	// Check someone isn't PMing us.
	if line.Public() == false {
		return
	}

	// All systems go!
	// Set LastTime to now.
	LastTime = time.Now().Unix()

	// Get some input, somewhat randomly.
	postWords := strings.Split(line.Text(), " ")
	dice = rand.Intn(len(postWords))
	prefix := strings.Join(postWords[dice:dice+2], " ")

	markov.AddLine(line.Text())
	markov.GenerateLine(prefix, config.LineLength)
}
