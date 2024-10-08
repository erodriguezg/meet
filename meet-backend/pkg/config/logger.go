package config

import "go.uber.org/zap"

var (
	log *zap.Logger
)

func configLoggers() {
	var err error
	if env == "PROD" {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}
