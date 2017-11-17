package runtime

import (
	log "github.com/sirupsen/logrus"

	"github.com/genofire/golang-lib/websocket"
)

type WebsocketHandlerFunc func(*log.Entry, *websocket.Message) error

var websockethandlermap = make(map[string]WebsocketHandlerFunc)

func WebsocketHandler(inputMSG chan *websocket.Message) {
	for msg := range inputMSG {
		logger := log.WithFields(log.Fields{"session": msg.Session, "id": msg.ID})
		if handler, ok := websockethandlermap[msg.Subject]; ok {
			err := handler(logger, msg)
			if err != nil {
				logger.Errorf("websocket message '%s' cound not handle: %s", msg.Subject, err)
			}
		} else {
			logger.Warnf("websocket message '%s' cound not handle", msg.Subject)
		}
	}
}
