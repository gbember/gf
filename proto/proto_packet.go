package proto

import (
	"errors"
	"math"
)

const (
	_DEFAULT_BUFF_SIZE = 512
)

type Packet struct {
	pos  int
	l    int
	data []byte
}

func (p *Packet) SeekTo(n int) {
	p.pos = p.pos + n
}

func (p *Packet) Data() []byte {
	return p.data[0:p.pos]
}

func (p *Packet) Length() int {
	return p.pos
}

func (p *Packet) Reset() {
	p.pos = 0
}

func NewReader(data []byte) *Packet {
	return &Packet{data: data}
}

func NewWriter() *Packet {
	return &Packet{data: make([]byte, _DEFAULT_BUFF_SIZE, _DEFAULT_BUFF_SIZE), l: _DEFAULT_BUFF_SIZE}
}

//=============================================== Readers
func (p *Packet) readBool() (ret bool, err error) {
	b, _err := p.readInt8()

	if b != 1 {
		return false, _err
	}

	return true, _err
}

func (p *Packet) readArrayBool() (ret []bool, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]bool, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readBool()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readInt8() (ret int8, err error) {
	_ret, err := p.readUint8()
	ret = int8(_ret)
	return
}

func (p *Packet) readArrayInt8() (ret []int8, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]int8, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readInt8()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readUint8() (ret uint8, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read int8 failed")
		return
	}

	ret = uint8(p.data[p.pos])
	p.pos++
	return
}
func (p *Packet) readArrayUint8() (ret []uint8, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]uint8, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readUint8()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.readUint16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}
func (p *Packet) readArrayString() (ret []string, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]string, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readString()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readUint16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read int16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.pos += 2
	return
}

func (p *Packet) readArrayUint16() (ret []uint16, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]uint16, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readUint16()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readInt16() (ret int16, err error) {
	_ret, _err := p.readUint16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *Packet) readArrayInt16() (ret []int16, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]int16, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readInt16()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readUint32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read int32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	p.pos += 4
	return
}
func (p *Packet) readArrayUint32() (ret []uint32, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]uint32, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readUint32()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readInt32() (ret int32, err error) {
	_ret, _err := p.readUint32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *Packet) readArrayInt32() (ret []int32, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]int32, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readInt32()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readUint64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read int64 failed")
		return
	}

	ret = 0
	buf := p.data[p.pos : p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.pos += 8
	return
}
func (p *Packet) readArrayUint64() (ret []uint64, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]uint64, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readUint64()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readInt64() (ret int64, err error) {
	_ret, _err := p.readUint64()
	ret = int64(_ret)
	err = _err
	return
}
func (p *Packet) readArrayInt64() (ret []int64, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]int64, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readInt64()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readFloat32() (ret float32, err error) {
	bits, _err := p.readUint32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}
func (p *Packet) readArrayFloat32() (ret []float32, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]float32, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readFloat32()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

func (p *Packet) readFloat64() (ret float64, err error) {
	bits, _err := p.readUint64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *Packet) readArrayFloat64() (ret []float64, err error) {
	size, err := p.readUint16()
	if err != nil {
		return
	}
	ret = make([]float64, 0, size)
	for i := uint16(0); i < size; i++ {
		b, err := p.readFloat64()
		if err != nil {
			return ret, err
		}
		ret = append(ret, b)
	}
	return
}

//================================================ Writers

func (p *Packet) writeBool(v bool) {
	p.grow(1)
	if v {
		p.data[p.pos] = 1
	} else {
		p.data[p.pos] = 0
	}
	p.pos++
}

func (p *Packet) writeArrayBool(v []bool) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeBool(v[i])
	}
}

func (p *Packet) writeInt8(v int8) {
	p.grow(1)
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayInt8(v []int8) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeInt8(v[i])
	}
}

func (p *Packet) writeUint8(v uint8) {
	p.grow(1)
	p.data[p.pos] = v
	p.pos++
}

func (p *Packet) writeArrayUint8(v []uint8) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeUint8(v[i])
	}
}

func (p *Packet) writeString(v string) {
	p.writeUint16(uint16(len(v)))
	bs := []byte(v)
	l := len(bs)
	p.grow(l)
	for i := 0; i < l; i++ {
		p.data[p.pos] = bs[i]
		p.pos++
	}
}
func (p *Packet) writeArrayString(v []string) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeString(v[i])
	}
}

func (p *Packet) writeUint16(v uint16) {
	p.grow(2)
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayUint16(v []uint16) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeUint16(v[i])
	}
}

func (p *Packet) writeInt16(v int16) {
	p.grow(2)
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayInt16(v []int16) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeInt16(v[i])
	}
}

func (p *Packet) writeInt32(v int32) {
	p.grow(4)
	p.data[p.pos] = byte(v >> 24)
	p.pos++
	p.data[p.pos] = byte(v >> 16)
	p.pos++
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayInt32(v []int32) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeInt32(v[i])
	}
}

func (p *Packet) writeUint32(v uint32) {
	p.grow(4)
	p.data[p.pos] = byte(v >> 24)
	p.pos++
	p.data[p.pos] = byte(v >> 16)
	p.pos++
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayUint32(v []uint32) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeUint32(v[i])
	}
}

func (p *Packet) writeUint64(v uint64) {
	p.grow(8)
	p.data[p.pos] = byte(v >> 56)
	p.pos++
	p.data[p.pos] = byte(v >> 48)
	p.pos++
	p.data[p.pos] = byte(v >> 40)
	p.pos++
	p.data[p.pos] = byte(v >> 32)
	p.pos++
	p.data[p.pos] = byte(v >> 24)
	p.pos++
	p.data[p.pos] = byte(v >> 16)
	p.pos++
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayUint64(v []uint64) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeUint64(v[i])
	}
}

func (p *Packet) writeInt64(v int64) {
	p.grow(8)
	p.data[p.pos] = byte(v >> 56)
	p.pos++
	p.data[p.pos] = byte(v >> 48)
	p.pos++
	p.data[p.pos] = byte(v >> 40)
	p.pos++
	p.data[p.pos] = byte(v >> 32)
	p.pos++
	p.data[p.pos] = byte(v >> 24)
	p.pos++
	p.data[p.pos] = byte(v >> 16)
	p.pos++
	p.data[p.pos] = byte(v >> 8)
	p.pos++
	p.data[p.pos] = byte(v)
	p.pos++
}
func (p *Packet) writeArrayInt64(v []int64) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeInt64(v[i])
	}
}

func (p *Packet) writeFloat32(f float32) {
	v := math.Float32bits(f)
	p.writeUint32(v)
}
func (p *Packet) writeArrayFloat32(v []float32) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeFloat32(v[i])
	}
}

func (p *Packet) writeFloat64(f float64) {
	v := math.Float64bits(f)
	p.writeUint64(v)
}
func (p *Packet) writeArrayFloat64(v []float64) {
	l := uint16(len(v))
	p.writeUint16(l)
	for i := uint16(0); i < l; i++ {
		p.writeFloat64(v[i])
	}
}

func (p *Packet) grow(n int) {
	if n+p.pos > p.l {
		l := p.l + _DEFAULT_BUFF_SIZE + n
		bs := make([]byte, l, l)
		copy(bs, p.data)
		p.data = bs
		p.l = l
	}
}
