package runtime

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/genofire/golang-lib/database"
	"github.com/genofire/golang-lib/websocket"
	"github.com/genofire/warehost/data"
)

type LogRead struct {
	Entry *log.Entry
}

func NewLogRead() *LogRead {
	return &LogRead{}
}
func (l *LogRead) Levels() []log.Level {
	return []log.Level{log.InfoLevel, log.WarnLevel, log.ErrorLevel}
}
func (l *LogRead) Fire(e *log.Entry) error {
	l.Entry = e
	return nil
}

func TestLoginHandler(t *testing.T) {
	assert := assert.New(t)
	logread := NewLogRead()
	log.AddHook(logread)
	logger := log.WithField("a", 1)

	database.Open(database.Config{
		Type:       "sqlite3",
		Logging:    false,
		Connection: fmt.Sprintf("file:databaseWSAuth?mode=memory"),
	})

	data.CreateDatabase()

	out := make(chan *websocket.Message, 1)
	dummyClient := websocket.NewTestClient(out)

	// Invalid
	err := loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body:    "b",
	})
	assert.NoError(err)
	msg := <-out
	assert.False(msg.Body.(bool))

	// user not found
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "notextisinguser",
			Password: "root",
		},
	})
	assert.NoError(err)
	assert.Equal("user not found", logread.Entry.Message)
	msg = <-out
	assert.False(msg.Body.(bool))

	// wrong password
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "roota",
		},
	})
	assert.NoError(err)
	assert.Equal("wrong password", logread.Entry.Message)
	msg = <-out
	assert.False(msg.Body.(bool))

	// login
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("loggedin success", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.True(msg.Body.(bool))

	// login again
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("already loggedIn", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.True(msg.Body.(bool))

	login := loggedIn[uuid.Nil]
	login.Active = false
	login.Password = "pbkdf2_sha1$10000$a5viM+Paz3o=$orD4shu1Ss+1wPAhAt8hkZ/fH7Y="
	database.Write.Save(login)
	delete(loggedIn, uuid.Nil)

	// login again
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("user not active", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.False(msg.Body.(bool))

	login.Active = true
	database.Write.Save(login)

	// login again
	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("loggedin with saving new hashed password", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.True(msg.Body.(bool))

	// login - check if new hashed password was saved correct
	delete(loggedIn, uuid.Nil)

	err = loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("loggedin success", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.True(msg.Body.(bool))

	delete(loggedIn, uuid.Nil)
}

func TestLogoutHandler(t *testing.T) {
	assert := assert.New(t)
	logread := NewLogRead()
	log.AddHook(logread)
	logger := log.WithField("a", 1)

	delete(loggedIn, uuid.Nil)

	database.Open(database.Config{
		Type:       "sqlite3",
		Logging:    false,
		Connection: fmt.Sprintf("file:databaseWSAuth?mode=memory"),
	})

	data.CreateDatabase()

	out := make(chan *websocket.Message, 1)
	dummyClient := websocket.NewTestClient(out)

	// login
	err := loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("loggedin success", logread.Entry.Message)
	assert.NoError(err)
	msg := <-out
	assert.True(msg.Body.(bool))

	// during login
	err = logoutHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "logout",
	})
	assert.Equal("logout", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.True(msg.Body.(bool))

	// during logout
	err = logoutHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "logout",
	})
	assert.Equal("logout without login", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.False(msg.Body.(bool))
}

func TestAuthStatusHandler(t *testing.T) {
	assert := assert.New(t)
	logread := NewLogRead()
	log.AddHook(logread)
	logger := log.WithField("a", 1)

	delete(loggedIn, uuid.Nil)

	database.Open(database.Config{
		Type:       "sqlite3",
		Logging:    false,
		Connection: fmt.Sprintf("file:databaseWSAuth?mode=memory"),
	})

	data.CreateDatabase()

	out := make(chan *websocket.Message, 1)
	dummyClient := websocket.NewTestClient(out)

	// login
	err := loginHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "login",
		Body: data.RequestLogin{
			Username: "root",
			Password: "root",
		},
	})
	assert.Equal("loggedin success", logread.Entry.Message)
	assert.NoError(err)
	msg := <-out
	assert.True(msg.Body.(bool))

	// during login
	err = authStatusHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "auth_status",
	})
	assert.Equal("welcome back 1", logread.Entry.Message)
	assert.NoError(err)
	msg = <-out
	assert.Equal(int64(1), msg.Body.(*data.Login).ID)

	// during logout
	delete(loggedIn, uuid.Nil)
	err = authStatusHandler(logger, &websocket.Message{
		From:    dummyClient,
		Subject: "auth_status",
	})
	assert.NoError(err)
	msg = <-out
	assert.False(msg.Body.(bool))
}
