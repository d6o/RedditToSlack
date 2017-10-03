package slack

import (
	"sync/atomic"

	"fmt"
	"github.com/disiqueira/RedditToSlack/pkg/slack/rtm"
	"golang.org/x/net/websocket"
	"time"
)

const (
	originURL         = "https://api.slack.com/"
	counterInitial    = 1
	counterIncrement  = 1
	webSocketProtocol = ""
)

//New TODO
func New(rtm *rtm.Response) (a *Agent, err error) {
	a = &Agent{
		counter: counterInitial,
	}
	err = a.connect(rtm.URL)
	return
}

//Agent TODO
type Agent struct {
	ws      *websocket.Conn
	counter uint64
}

func (a *Agent) connect(rtm string) (err error) {
	a.ws, err = websocket.Dial(rtm, webSocketProtocol, originURL)
	return
}

//Message TODO
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subtype,omitempty"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
	Ts      string `json:"ts"`
}

//SendMessage sends a Message to a channel.
func (a *Agent) SendMessage(m Message) error {
	m.ID = atomic.AddUint64(&a.counter, counterIncrement)
	m.Ts = fmt.Sprintf("%d", time.Now().Unix())
	return websocket.JSON.Send(a.ws, m)
}
