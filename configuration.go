package timboslice

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	DBPath            string   `mapstructure:"db_path"`
	IRCNickname       string   `mapstructure:"irc_nickname"`
	IRCUsername       string   `mapstructure:"irc_username"`
	IRCRealName       string   `mapstructure:"irc_realname"`
	IRCPassword       string   `mapstructure:"irc_password"`
	ServerHost        string   `mapstructure:"server_host"`
	ServerPort        int      `mapstructure:"server_port"`
	ServerTLS         bool     `mapstructure:"server_tls"`
	ServerPassword    string   `mapstructure:"server_password"`
	Channels          []string `mapstructure:"channels"`
	Ignored           []string `mapstructure:"ignored"`
	Moderators        []string `mapstructure:"moderators"`
	Chance            int      `mapstructure:"chance"`
	Silenced          bool     `mapstructure:"silenced"`
	SilenceToggle     string   `mapstructure:"silence_toggle"`
	MinWaitSeconds    int64    `mapstructure:"min_wait_seconds"`
	MaxLineLength     int      `mapstructure:"max_line_length"`
	ChainPrefixLength int      `mapstructure:"chain_prefix_length"`
}

func NewConfiguration() Configuration {
	var config Configuration

	viper.SetConfigName("timboslice")
	viper.AddConfigPath("/etc/timboslice/")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(err.Error())
	}

	viper.Unmarshal(&config)

	log.Printf(`Loaded configuration:
- Nickname: %s
- Server: %s:%d
- TLS: %t`,
		config.IRCNickname, config.ServerHost, config.ServerPort, config.ServerTLS)

	return config
}
