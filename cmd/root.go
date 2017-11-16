package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "warehost",
	Short: "warehost to manage a server",
	Long:  `Warehost is a little web gui to manage your (mail) server`,
}

func Execute() {
	if err := RootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
