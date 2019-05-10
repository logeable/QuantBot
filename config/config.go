package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Database struct {
		Driver string
		DSN    string
	}
	Server struct {
		Addr string
	}
}

var Config *config

func init() {
	f, err := os.Open(os.ExpandEnv(`$HOME/.quantbot/config.yml`))
	if err != nil {
		log.Fatal("open config file failed:", err)
	}

	Config = new(config)
	err = yaml.NewDecoder(f).Decode(Config)
	if err != nil {
		log.Fatal("decode config file failed:", err)
	}
}
