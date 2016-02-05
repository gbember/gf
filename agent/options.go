// options.go
package agent

import (
	"log"
	"os"
	"time"

	"github.com/gbember/util"
)

type Options struct {
	Id                   uint64        //agent编号id
	Name                 string        //名字
	TCPAddress           string        //与客户端的监听端口地址
	CenterTCPAddress     string        //中心服务器端口地址
	MaxMsgSize           uint16        //接受客户端的最大消息
	MaxHeartbeatInterval time.Duration //心跳间隔
	MsgSpeedCtrlNum      int           //速度控制 个数
	MsgSpeedCtrlSec      time.Duration //速度控制 时间
	MsgSpeedETimes       int           //速度控制 最大超出次数(超出次数后踢掉)
	MaxOnlineNum         int32         //最高在线
	ReadBuffSize         int
	WriteBuffSize        int
}

func NewOptions() *Options {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	return &Options{
		Id:                   util.UUID(),
		Name:                 hostname,
		TCPAddress:           "0.0.0.0:40015",
		CenterTCPAddress:     "0.0.0.0:50015",
		MaxMsgSize:           16 * 1024,
		MaxHeartbeatInterval: 60 * time.Second,
		MsgSpeedCtrlNum:      100,
		MsgSpeedCtrlSec:      time.Second,
		MsgSpeedETimes:       10,
		MaxOnlineNum:         100000,
		ReadBuffSize:         8 * 1024,
		WriteBuffSize:        32 * 1024,
	}
}
