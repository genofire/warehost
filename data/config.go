package data

import "github.com/genofire/golang-lib/database"

type Config struct {
	Address  string          `toml:"address"`
	Webroot  string          `toml:"webroot"`
	Database database.Config `toml:"database"`
}
