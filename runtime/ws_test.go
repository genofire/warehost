package runtime

import (
	"errors"
	"sync"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/genofire/golang-lib/websocket"
)

type ThreadSafeBool struct {
	sync.Mutex
	value bool
}

func (t *ThreadSafeBool) Get() bool {
	t.Lock()
	defer t.Unlock()
	return t.value
}
func (t *ThreadSafeBool) Set(v bool) {
	t.Lock()
	t.value = v
	t.Unlock()
}

func TestWebsocketHandler(t *testing.T) {
	assert := assert.New(t)
	inputMSG := make(chan *websocket.Message, 1)
	runned := ThreadSafeBool{}

	go WebsocketHandler(inputMSG)

	websockethandlermap["errortest"] = func(log *log.Entry, msg *websocket.Message) error {
		runned.Set(true)
		assert.Equal("errortest", msg.Subject)
		return errors.New("some")
	}

	assert.False(runned.Get())
	inputMSG <- &websocket.Message{
		Subject: "errortest",
	}
	time.Sleep(time.Millisecond)
	assert.True(runned.Get())

	runned.Set(false)
	inputMSG <- &websocket.Message{
		Subject: "notfound",
	}
	time.Sleep(time.Millisecond)
	assert.False(runned.Get())
}
