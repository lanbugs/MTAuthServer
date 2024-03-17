package mtauthserver

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig() Config {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	var cnf Config

	err := viper.Unmarshal(&cnf)

	if err != nil {
		log.Fatalf("unable to decode config, %v", err)
	}

	return cnf
}
