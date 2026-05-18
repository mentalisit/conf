package conf

import (
	"context"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/mentalisit/conf/config"
	"github.com/mentalisit/conf/logger"
)

func InitConf(serviceName string) (cfg *config.BotConfig, log *logger.Logger, clientDb *sqlx.DB) {
	cfg = config.InitConfig()
	log = logger.LoggerZap(cfg.Logger.Token, cfg.Logger.ChatId, cfg.Logger.Webhook, serviceName)
	db, err := newClientDb(log, cfg)
	if err != nil {
		return nil, nil, nil
	}

	return cfg, log, db
}

func InitConfDev() (cfg *config.BotConfig, log *logger.Logger, clientDb *sqlx.DB) {
	cfg = config.InitConfig()
	log = logger.LoggerZapDEV()
	db, err := newClientDb(log, cfg)
	if err != nil {
		return nil, nil, nil
	}

	return cfg, log, db
}

func newClientDb(log *logger.Logger, conf *config.BotConfig) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	err = doWithTries(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		db, err = sqlx.ConnectContext(ctx, "postgres", conf.GetDNS())
		if err != nil {
			db, err = sqlx.ConnectContext(ctx, "postgres", conf.GetDnsLan())
			if err != nil {
				log.ErrorErr(err)
				os.Exit(1)
			}
		}
		return nil
	}, 5, 5*time.Second)
	if err != nil {
		log.Fatal(err.Error())
	}

	return db, nil
}

func doWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
