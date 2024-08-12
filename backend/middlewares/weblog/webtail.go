package weblog

import (
	"github.com/fasthttp/websocket"
	"github.com/linxlib/fw"
	"sync"

	"github.com/sirupsen/logrus"
)

// codebeat:disable[TOO_MANY_IVARS]

// Config defines local application flags
type Config struct {
	Root        string `long:"root"  default:"log/"  description:"Root directory for log files"`
	Bytes       int64  `long:"bytes" default:"5000"  description:"tail from the last Nth location"`
	Lines       int    `long:"lines" default:"100"   description:"keep N old lines for new consumers"`
	MaxLineSize int    `long:"split" default:"180"   description:"split line if longer"`
	ListCache   int    `long:"cache" default:"2"      description:"Time to cache file listing (sec)"`
	Poll        bool   `long:"poll"  description:"use polling, instead of inotify"`
	Trace       bool   `long:"trace" description:"trace worker channels"`

	ClientBufferSize  int `long:"out_buf"      default:"256"  description:"Client Buffer Size"`
	WSReadBufferSize  int `long:"ws_read_buf"  default:"1024" description:"WS Read Buffer Size"`
	WSWriteBufferSize int `long:"ws_write_buf" default:"1024" description:"WS Write Buffer Size"`
}

// codebeat:enable[TOO_MANY_IVARS]

// Service holds WebTail service
type Service struct {
	cfg *Config
	hub *Hub
	wg  *sync.WaitGroup
	log logrus.FieldLogger
}

// New creates WebTail service
func New(log logrus.FieldLogger, cfg *Config) (*Service, error) {
	tail, err := NewTailService(log, cfg)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	hub := NewHub(log, tail, &wg)
	service := Service{cfg: cfg, hub: hub, log: log, wg: &wg}
	return &service, nil
}

// Run runs a message hub
func (wt *Service) Run() {
	wt.hub.Run()
}

// Close stops a message hub
func (wt *Service) Close() {
	wt.log.Info("Service Exiting")
	wt.hub.Close()
	wt.wg.Wait()
}

// Handle handles websocket requests from the peer
func (wt *Service) ServeHTTP(context *fw.Context) {
	wsUpgrader := upgrader(wt.cfg.WSReadBufferSize, wt.cfg.WSWriteBufferSize)
	err := wsUpgrader.Upgrade(context.GetFastContext(), func(conn *websocket.Conn) {
		client := &Client{
			conn: conn,
			send: make(chan []byte, wt.cfg.ClientBufferSize),
			log:  wt.log,
		}
		wt.hub.register <- client
		// Allow collection of memory referenced by the caller by doing all work in
		// new goroutines.
		go client.runWritePump(wt.wg)
		client.runReadPump(wt.wg, wt.hub.unregister, wt.hub.broadcast)
	})
	if err != nil {
		wt.log.Error(err, " Upgrade connection")
		return
	}

}
