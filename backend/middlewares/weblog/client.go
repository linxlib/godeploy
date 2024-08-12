package weblog

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	log logrus.FieldLogger
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// runReadPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) runReadPump(wg *sync.WaitGroup, quit chan *Client, inbox chan *Message) {
	wg.Add(1)
	defer func() {
		c.log.Error(" Client runReadPump err ", recover())
		quit <- c
		err := c.conn.Close()
		if err != nil {
			c.log.Error(err, " ReadPump conn Close")
		}
		wg.Done()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		c.log.Error(err, " SetReadDeadline")
		return
	}

	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				c.log.Error(err, " UnexpectedCloseError")
			}
			c.log.Error(err, " ReadMessage")
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		inbox <- &Message{Client: c, Message: message}
	}
}

// runWritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) runWritePump(wg *sync.WaitGroup) {
	wg.Add(1)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			c.log.Error(err, " WritePump conn Close")
		}
		wg.Done()
	}()
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				c.log.Error(err, " SetWriteDeadline")
				return
			}
			if !ok {
				// The hub closed the channel. Send Bye and exit
				//c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.sendMessage(message); err != nil {
				c.log.Error(err, " sendMessage")
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				err = c.conn.WriteMessage(websocket.PingMessage, nil)
			}
			if err != nil {
				return
			}
		}
	}
}

func (c *Client) sendMessage(message []byte) error {
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		c.log.Error(err, " NextWriter")
		return err
	}
	_, err = w.Write(message)
	if err != nil {
		c.log.Error(err, " Write")
		return err
	}
	// Add queued chat messages to the current websocket message.
	n := len(c.send)
	for i := 0; i < n; i++ {
		_, err = w.Write(newline)
		if err == nil {
			_, err = w.Write(<-c.send)
		}
		if err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

func upgrader(readBufferSize, writeBufferSize int) websocket.FastHTTPUpgrader {
	return websocket.FastHTTPUpgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return true
		},
	}
}
