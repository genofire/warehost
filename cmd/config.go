package cmd

import (
	"log"

	"github.com/genofire/golang-lib/file"
	"github.com/genofire/warehost/data"
)

var configPath string

func loadConfig() *data.Config {
	config := &data.Config{}
	err := file.ReadTOML(configPath, config)
	if err != nil {
		log.Panic(err)
	}
	return config
}
