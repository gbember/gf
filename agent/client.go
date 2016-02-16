// client.go
package agent

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gbember/gf/proto"
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
	pkw            *proto.Packet
	stream         interface{} //与game之间的流
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
		pkw:            proto.NewWriter(),
	}
	go client.start()

	return client
}
func (self *Client) start() {
	self.check_heart()
	go self.ioLoop()
}
func (self *Client) stop() {
	close(self.exitChan)
	self.gtimer.Stop()
	self.agent.RemoveClient(self)
	self.Lock()
	self.isClosed = true
	self.Unlock()
	self.conn.Close()
}

func (self *Client) ioLoop() {
	defer util.PrintPanicStack(self.agent.Error)
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
		if size < 2 {
			self.agent.Error("msg too short, short:%v size:%v", 2, size)
			return
		}
		body := buf[:size]
		n, err = io.ReadFull(self.conn, body)
		if err != nil {
			self.agent.Warn("read body failed, reason:%v size:%v", err, n)
			return
		}
		self.dispatcher(body)
	}

}

//消息路由
func (self *Client) dispatcher(msgRaw []byte) {
	mid := binary.BigEndian.Uint16(msgRaw[:2])
	id := mid / 100
	if id == 100 {
		mid, msg, err := proto.DecodeProto(msgRaw)
		if err != nil {
			self.agent.Error(err)
			self.Kick(0)
			return
		}
		switch mid {
		case proto.CS_ACCOUNT_AUTH:
			//验证
			self.auth(msg.(*proto.Cs_account_auth))
		case proto.CS_ACCOUNT_CREATE:
			//创建角色
			self.create(msg.(*proto.Cs_account_create))
		case proto.CS_ACCOUNT_HEART:
			//心跳
			self.heart(msg.(*proto.Cs_account_heart))
		default:
			self.agent.Error("无效消息:~d", mid)
			self.Kick(0)
			return
		}
		return
	}

	//通过流发送到game处理
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

func (self *Client) Kick(code int8) {
	if code != 0 {
		self.Send(&proto.Sc_account_kick{Reason: code})
	}
	self.conn.Close()
	self.Lock()
	self.isClosed = true
	self.Unlock()
}
func (self *Client) Send(msg proto.Messager) {
	self.pkw.SeekTo(2)
	data := proto.EncodeProtoPacket(msg, self.pkw)
	bs := data[:2]
	binary.BigEndian.PutUint16(bs, uint16(len(data)-2))
	self.conn.Write(data)
}

//验证
func (self *Client) auth(authData *proto.Cs_account_auth) {

}

//心跳
func (self *Client) heart(pingData *proto.Cs_account_heart) {
	self.hearNum++
	self.Send(&proto.Sc_account_heart{UnixTime: int64(util.NowStamp())})
}

//创建角色
func (self *Client) create(createData *proto.Cs_account_create) {

}
