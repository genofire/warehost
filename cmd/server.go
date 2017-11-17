package cmd

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/genofire/golang-lib/database"
	"github.com/genofire/golang-lib/websocket"
	"github.com/genofire/warehost/data"
	"github.com/genofire/warehost/runtime"
)

// serveCMD represents the entrance command
var serverCMD = &cobra.Command{
	Use:   "server",
	Short: "server by a websocket",
	Long:  `server`,

	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()

		err := database.Open(config.Database)
		if err != nil {
			os.Exit(111)
		}
		defer database.Close()
		data.CreateDatabase()

		inputMSG := make(chan *websocket.Message)
		ws := websocket.NewServer(inputMSG, websocket.NewSessionManager())

		http.HandleFunc("/ws", ws.Handler)
		http.Handle("/", http.FileServer(http.Dir(config.Webroot)))

		srv := &http.Server{
			Addr: config.Address,
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil {
				panic(err)
			}
		}()

		log.Println("started under:", srv.Addr)

		go runtime.WebsocketHandler(inputMSG)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigs
		log.Println("received", sig)

		close(inputMSG)
		srv.Close()

		log.Println("shutdown warehost")
	},
}

func init() {
	serverCMD.Flags().StringVarP(&configPath, "config", "c", "warehost.conf", "Path to configuration file")
	RootCMD.AddCommand(serverCMD)
}
