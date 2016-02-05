// client.go
package agent

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gbember/util"
)

type Client struct {
	agent          *Agent
	msgNum         int
	startTime      time.Time
	checkSpeedTime time.Time
	hearNum        int //上一次心跳检查到当前的心跳次数
	conn           net.Conn
	exitChan       chan bool
	isClosed       bool
	gtimer         *util.GTimer
	sync.RWMutex
}

func NewClient(agent *Agent, conn net.Conn) *Client {
	now := time.Now()
	client := &Client{
		agent:          agent,
		checkSpeedTime: now,
		startTime:      now,
		conn:           conn,
		exitChan:       make(chan bool),
		gtimer:         util.NewGTimer(),
	}
	go client.start()

	return client
}
func (self *Client) start() {
	self.check_heart()
	go self.ioLoop()
}
func (self *Client) stop() {
	self.agent.RemoveClient(self)
	self.conn.Close()
}

func (self *Client) ioLoop() {
	defer self.stop()
	header := make([]byte, 2)
	maxSize := self.agent.opts.MaxMsgSize
	buf := make([]byte, maxSize)

	for {
		n, err := io.ReadFull(self.conn, header)
		if err != nil {
			self.agent.Warn("read header failed, reason:%v size:%v", err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)
		if size > maxSize {
			self.agent.Error("msg too long, max:%v size:%v", maxSize, size)
			return
		}
		body := buf[:size]
		n, err = io.ReadFull(self.conn, body)
		if err != nil {
			self.agent.Warn("read body failed, reason:%v size:%v", err, n)
			return
		}
		//TODO 处理消息
		//读取消息头
		//转发消息或执行消息处理函数
		//只处理心跳和登录
		//登录验证服务器验证
	}

}

//心跳检测
func (self *Client) check_heart() {
	self.RLock()
	if self.isClosed {
		self.Unlock()
		return
	}
	if self.hearNum == 0 {
		self.RUnlock()
		self.Kick(1)
		return
	}
	self.RUnlock()
	self.Lock()
	self.hearNum = 0
	self.Unlock()
	self.gtimer.AddAfter(self.agent.opts.MaxHeartbeatInterval*2, self.check_heart)
}

func (self *Client) Kick(code int) {
	close(self.exitChan)
	self.Lock()
	self.isClosed = true
	self.Unlock()
}
