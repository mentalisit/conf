package config

import (
	"fmt"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type BotConfig struct {
	Token struct {
		TokenDiscord             string `yaml:"token_discord"`
		TokenTelegram            string `yaml:"token_telegram"`
		WhatsappNumber           string `yaml:"whatsapp_number"`
		WhatsappSessionFile      string `yaml:"whatsapp_session_file"`
		DiscordOAuthClientID     string `yaml:"discord_oauth_client_id"`
		DiscordOAuthClientSecret string `yaml:"discord_oauth_client_secret"`
		WhiteStarStatistic       string `yaml:"white_star_statistic"`
	} `yaml:"token"`
	Logger struct {
		Token   string `yaml:"token"`
		ChatId  int64  `yaml:"chat_id"`
		Webhook string `yaml:"webhook"`
	} `yaml:"logger"`
	Postgres struct {
		Host     string `yaml:"host" env-default:"postgres:5432"`
		Name     string `yaml:"name" env-default:"postgres"`
		Username string `yaml:"username" env-default:"postgres"`
		Password string `yaml:"password" env-default:"root"`
	} `yaml:"postgres"`
}

var Instance *BotConfig
var once sync.Once

func InitConfig() *BotConfig {
	once.Do(func() {
		Instance = &BotConfig{}
		err := cleanenv.ReadConfig("docker/config/configs.yml", Instance)
		if err != nil {
			help, _ := cleanenv.GetDescription(Instance, nil)
			fmt.Println(help)
			panic("NULL")
		}
	})
	return Instance
}

func (b *BotConfig) GetDNS() string {
	if b == nil || b.Postgres.Password == "" {
		return ""
	}
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		b.Postgres.Username, b.Postgres.Password, b.Postgres.Host, b.Postgres.Name)
	return dns
}

func (b *BotConfig) GetDnsLan() string {
	if b == nil || b.Postgres.Password == "" {
		return ""
	}
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		b.Postgres.Username, b.Postgres.Password, "Dell:5434", b.Postgres.Name)
	return dns
}
