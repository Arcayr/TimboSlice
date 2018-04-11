package timboslice

import (
	"crypto/tls"
	"fmt"
	"github.com/elliotspeck/markov"
	"github.com/elliotspeck/markov/storage"
	irc "github.com/fluffle/goirc/client"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Tim struct {
	Config     Configuration
	Chain      markov.Chain
	IRCClient  *irc.Conn
	LastPosted int64
}

func NewTim() Tim {
	var tim Tim

	tim.Config = NewConfiguration()

	tim.Chain = markov.Chain{PrefixLen: tim.Config.ChainPrefixLength}
	if storage, err := storage.Load("sqlite3", tim.Config.DBPath); err != nil {
		panic(err.Error())
	} else {
		tim.Chain.Storage = storage
	}

	ircConfig := irc.NewConfig(tim.Config.IRCNickname, tim.Config.IRCUsername, tim.Config.IRCRealName)
	ircConfig.SSL = tim.Config.ServerTLS
	ircConfig.SSLConfig = &tls.Config{ServerName: tim.Config.ServerHost}
	ircConfig.Server = fmt.Sprintf("%s:%d", tim.Config.ServerHost, tim.Config.ServerPort)

	tim.IRCClient = irc.Client(ircConfig)
	tim.IRCClient.EnableStateTracking()

	return tim
}

func (tim *Tim) handleConnect(conn *irc.Conn, line *irc.Line) {
	ghostMsg := fmt.Sprintf("GHOST %s %s", tim.Config.IRCNickname, tim.Config.IRCPassword)
	identMsg := fmt.Sprintf("IDENTIFY %s %s", tim.Config.IRCNickname, tim.Config.IRCPassword)

	time.Sleep(2 * time.Second)
	conn.Privmsg("NICKSERV", ghostMsg)
	time.Sleep(2 * time.Second)
	conn.Privmsg("NICKSERV", identMsg)

	for _, channel := range tim.Config.Channels {
		time.Sleep(1 * time.Second)
		log.Println("Joining ", channel)
		conn.Join(channel)
	}
}

func (tim *Tim) handleDisconnect(conn *irc.Conn, line *irc.Line) {
	log.Println("Disconnected from the server, retrying...")
	if err := tim.IRCClient.Connect(); err != nil {
		log.Fatal("No dice. Shutting down.")
	}
}

func (tim *Tim) handleKick(conn *irc.Conn, line *irc.Line) {
	if strings.ToLower(line.Args[1]) == strings.ToLower(tim.IRCClient.Config().Me.Nick) {
		log.Printf("Kicked from %s, rejoining...", line.Target())
		conn.Join(line.Target())
	} else {
		conn.Privmsg(line.Target(), "Great! Bye now")
	}
}

func (tim *Tim) handlePrivmsg(conn *irc.Conn, line *irc.Line) {
	// Ensure the message was to a channel we're in, not a query.
	if line.Public() == false {
		return
	}

	// Check if the user is annoying. Break out of this if so.
	for _, nick := range tim.Config.Ignored {
		if line.Nick == nick {
			return
		}
	}

	// Check if we're being told to stop posting.
	if strings.ToLower(line.Text()) == strings.ToLower(tim.Config.SilenceToggle) {
		for _, nick := range tim.Config.Moderators {
			if line.Nick == nick {
				tim.Config.Silenced = !tim.Config.Silenced
				log.Printf("Toggling silenced. Now: %t", tim.Config.Silenced)
				return
			}
		}
	}

	// All modifiers are passed, add the text to the chain pool.
	tim.Chain.AddLine(line.Text())

	// No need to go any further if we're silenced.
	if tim.Config.Silenced {
		return
	}

	// If the bot wasn't pinged, we want to test for minwait and rand.
	if strings.Contains(strings.ToLower(line.Text()), strings.ToLower(tim.IRCClient.Config().Me.Nick)) == false {
		// Checking minwait
		if tim.LastPosted+tim.Config.MinWaitSeconds >= time.Now().Unix() {
			return
		}

		// Checking rand
		// Rand functions as a "one in x" chance, where "x" is the configured chance.
		if r := rand.Intn(tim.Config.Chance); r != 1 {
			return
		}
	}

	// If we reach this point, it's time to post.
	generatedLine, err := tim.Chain.GenerateLine(tim.Config.MaxLineLength)
	if err != nil {
		return
	}

	tim.IRCClient.Privmsg(line.Target(), generatedLine)
	tim.LastPosted = time.Now().Unix()
}

func (tim *Tim) Connect() {
	if err := tim.IRCClient.Connect(); err != nil {
		panic(err.Error())
	}

	tim.IRCClient.HandleFunc(irc.CONNECTED, tim.handleConnect)
	tim.IRCClient.HandleFunc(irc.DISCONNECTED, tim.handleDisconnect)
	tim.IRCClient.HandleFunc(irc.PRIVMSG, tim.handlePrivmsg)
	tim.IRCClient.HandleFunc(irc.KICK, tim.handleKick)
}
