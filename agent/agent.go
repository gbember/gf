// agent.go
package agent

import (
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	l4g "code.google.com/p/log4go"
)

type Agent struct {
	opts        Options
	tcpListener net.Listener
	isClosed    bool
	exitChan    chan bool
	startTime   time.Time
	onlineNum   int32 //在线人数
	clients     map[*Client]struct{}
	l4g.Logger
	sync.WaitGroup
	sync.RWMutex
}

func NewAgent(opts Options) *Agent {
	log := l4g.NewLogger()
	log.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("agent.log", false))
	return &Agent{
		opts:      opts,
		exitChan:  make(chan bool),
		startTime: time.Now(),
		clients:   make(map[*Client]struct{}, opts.MaxOnlineNum),
		Logger:    log,
	}
}

func (self *Agent) Main() {
	self.Lock()
	var tcpListener net.Listener
	tcpListener, err := net.Listen("tcp", self.opts.TCPAddress)
	if err != nil {
		self.Error("FATAL: listen (%s) failed - %s", self.opts.TCPAddress, err)
		os.Exit(1)
	}
	self.tcpListener = tcpListener
	self.Unlock()
	go self.acceptLoop()
}

func (self *Agent) Exit() {
	self.isClosed = true
	self.tcpListener.Close()
	self.Lock()
	clients := make([]*Client, len(self.clients))
	i := 0
	for client, _ := range self.clients {
		clients[i] = client
		i++
	}
	self.Unlock()
	for _, client := range clients {
		client.Kick(0)
	}
	self.Wait()
}

func (self *Agent) acceptLoop() {
	for {
		if self.isClosed {
			return
		}
		conn, err := self.tcpListener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				self.Error("temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				self.Error("listener.Accept() - %s", err)
			}
			break
		}
		self.Debug("client socket addr: %s", conn.RemoteAddr().Network())
		if self.GetOnlineNum() >= self.opts.MaxOnlineNum {
			//TODO 发送人数已满消息
			conn.Close()
			continue
		}

		conn.(*net.TCPConn).SetReadBuffer(self.opts.ReadBuffSize)
		conn.(*net.TCPConn).SetWriteBuffer(self.opts.WriteBuffSize)
		client := NewClient(self, conn)
		self.AddClient(client)
	}
}

func (self *Agent) GetOnlineNum() int32 {
	return atomic.LoadInt32(&self.onlineNum)
}

func (self *Agent) AddClient(client *Client) {
	self.Add(1)
	atomic.AddInt32(&self.onlineNum, 1)
	self.Lock()
	self.clients[client] = struct{}{}
	self.Unlock()
}
func (self *Agent) RemoveClient(client *Client) {
	self.Lock()
	delete(self.clients, client)
	self.Unlock()
	atomic.AddInt32(&self.onlineNum, -1)
	self.Done()
}
