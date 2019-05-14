package config

import (
	"log"
	"os"
	"path/filepath"

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
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(filepath.Join(home, `/.quantbot/config.yml`))
	if err != nil {
		log.Fatal("open config file failed:", err)
	}

	Config = new(config)
	err = yaml.NewDecoder(f).Decode(Config)
	if err != nil {
		log.Fatal("decode config file failed:", err)
	}
}
