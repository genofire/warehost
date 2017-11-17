package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/genofire/golang-lib/database"
	"github.com/genofire/warehost/data"
	"github.com/genofire/warehost/lib"
)

// authCMD represents the entrance command
var authCMD = &cobra.Command{
	Use:   "auth <username> <password>",
	Short: "validate a username with his password ",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		config := loadConfig()

		err := database.Open(config.Database)
		if err != nil {
			os.Exit(111)
		}
		defer database.Close()
		data.CreateDatabase()

		var realHashedPassword string
		err = database.Read.Raw("select password from login where mail = $1", username).Row().Scan(&realHashedPassword)
		if err != nil {
			os.Exit(3)
		}

		ok, _ := lib.Validate(realHashedPassword, password)
		if ok {
			os.Exit(0)
		} else {
			os.Exit(1)
		}

	},
}

func init() {
	authCMD.Flags().StringVarP(&configPath, "config", "c", "warehost.conf", "Path to configuration file")
	RootCMD.AddCommand(authCMD)
}
