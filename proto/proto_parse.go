package proto

import (
	"errors"
	"fmt"
)

type Messager interface {
	Read(p *Packet) error
	WriteMsgID(p *Packet)
	Write(p *Packet)
}
type PMessage struct {
	ID  uint16
	Msg Messager
}

func EncodeProto(msg Messager) []byte {
	p := NewWriter()
	msg.WriteMsgID(p)
	msg.Write(p)
	return p.Data()
}
func EncodeProtoPacket(msg Messager, p *Packet) []byte {
	msg.WriteMsgID(p)
	msg.Write(p)
	bs := p.Data()
	p.Reset()
	return bs
}
func DecodeProto(bin []byte) (msgID uint16, msg Messager, err error) {
	p := NewReader(bin)
	msgID, err = p.readUint16()
	if err != nil {
		return
	}
	mid := msgID / 100
	switch mid {
	case 100:
		switch msgID {
		case CS_ACCOUNT_CREATE:
			v := &Cs_account_create{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_CREATE:
			v := &Sc_account_create{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_KICK:
			v := &Sc_account_kick{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_HEART:
			v := &Cs_account_heart{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_HEART:
			v := &Sc_account_heart{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case CS_ACCOUNT_AUTH:
			v := &Cs_account_auth{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		case SC_ACCOUNT_AUTH:
			v := &Sc_account_auth{}
			err = v.Read(p)
			if err == nil {
				msg = v
			}
		default:
			err = errors.New(fmt.Sprintf("error: invalid msgID %d", msgID))
		}
	default:
		err = errors.New(fmt.Sprintf("error: invalid msgID %d", msgID))
	}
	return
}
