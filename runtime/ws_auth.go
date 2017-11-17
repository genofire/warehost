package runtime

import (
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"github.com/genofire/golang-lib/database"
	"github.com/genofire/golang-lib/websocket"
	"github.com/genofire/warehost/data"
	"github.com/genofire/warehost/lib"
)

var loggedIn = make(map[uuid.UUID]*data.Login)

func loginHandler(logger *log.Entry, msg *websocket.Message) error {
	_, ok := loggedIn[msg.Session]
	if ok {
		msg.Answer(msg.Subject, true)
		logger.Warn("already loggedIn")
		return nil
	}
	requestlogin := data.RequestLogin{}
	err := mapstructure.Decode(msg.Body, &requestlogin)
	if err != nil {
		msg.Answer(msg.Subject, false)
		return nil
	}

	logger = logger.WithField("username", requestlogin.Username)

	login := data.Login{}
	if database.Read.Where("mail = ?", requestlogin.Username).First(&login).RecordNotFound() {
		logger.Warn("user not found")
		msg.Answer(msg.Subject, false)
		return nil
	}
	if !login.Active {
		logger.Warn("user not active")
		msg.Answer(msg.Subject, false)
		return nil
	}
	ok, err = lib.Validate(login.Password, requestlogin.Password)
	if ok {
		loggedIn[msg.Session] = &login
		login.LastLoginAt = time.Now()
		if err == lib.ErrorHashDeprecated {
			login.Password = lib.NewHash(requestlogin.Password)
			logger.Info("loggedin with saving new hashed password")
		} else if err != nil {
			logger.Errorf("loggedin with password problem: %s", err)
		} else {
			logger.Info("loggedin success")
		}
		database.Write.Save(&login)
	} else {
		if err != nil {
			logger.Errorf("login failed with password problem: %s", err)
		} else {
			logger.Warn("wrong password")
		}
	}
	msg.Answer(msg.Subject, ok)
	return nil
}

func authStatusHandler(logger *log.Entry, msg *websocket.Message) error {
	login, ok := loggedIn[msg.Session]
	if !ok {
		msg.Answer(msg.Subject, false)
		return nil
	}
	logger.Infof("welcome back %d", login.ID)
	msg.Answer(msg.Subject, login)
	return nil
}
func logoutHandler(logger *log.Entry, msg *websocket.Message) error {
	_, ok := loggedIn[msg.Session]
	if !ok {
		msg.Answer(msg.Subject, false)
		logger.Warn("logout without login")
		return nil
	}
	delete(loggedIn, msg.Session)
	logger.Info("logout")
	msg.Answer(msg.Subject, true)
	return nil
}

func init() {
	websockethandlermap["login"] = loginHandler
	websockethandlermap["auth_status"] = authStatusHandler
	websockethandlermap["logout"] = logoutHandler
}
