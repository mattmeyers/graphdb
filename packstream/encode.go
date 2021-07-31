package packstream

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type encoder struct {
	buf *bytes.Buffer
}

type Marshaller interface {
	MarshalPackstream() ([]byte, error)
}

func Marshal(v interface{}) ([]byte, error) {
	e := encoder{buf: new(bytes.Buffer)}

	err := e.marshal(v)
	if err != nil {
		return nil, err
	}

	return e.buf.Bytes(), nil
}

func (e encoder) marshal(v interface{}) error {
	var err error

	switch val := v.(type) {
	case nil:
		e.buf.WriteByte(0xC0)
	case bool:
		e.encodeBool(val)
	case *bool:
		e.encodeBool(*val)
	case []byte:
		e.encodeBytes(val)
	case string:
		err = e.encodeString(val)
	case *string:
		err = e.encodeString(*val)
	case Marshaller:
		err = e.encodeMarshaller(val)
	default:
		err = fmt.Errorf("unable to marshal value of type %T", v)
	}

	return err
}

func (e encoder) encodeBool(v bool) {
	if v {
		e.buf.WriteByte(0xC3)
	} else {
		e.buf.WriteByte(0xC2)
	}
}

func (e encoder) encodeBytes(v []byte) error {
	if len(v) < (1 << 8) {
		e.buf.WriteByte(0xCC)
		e.buf.WriteByte(byte(len(v)))
	} else if len(v) < (1 << 16) {
		e.buf.WriteByte(0xCD)
		binary.Write(e.buf, binary.BigEndian, uint16(len(v)))
	} else if len(v) < (1 << 32) {
		e.buf.WriteByte(0xCE)
		binary.Write(e.buf, binary.BigEndian, uint32(len(v)))
	} else {
		return errors.New("cannot encode byte slices of length greater than 2,147,483,647 bytes")
	}

	e.buf.Write(v)

	return nil
}

// shortStringMarkers is a quick lookup table for strings of length less than
// 16 bytes.
var shortStringMarkers = [16]byte{
	0x80, 0x81, 0x82, 0x83,
	0x84, 0x85, 0x86, 0x87,
	0x88, 0x89, 0x8A, 0x8B,
	0x8C, 0x8D, 0x8E, 0x8F,
}

func (e encoder) encodeString(v string) error {
	if len(v) < (1 << 4) {
		e.buf.WriteByte(shortStringMarkers[len(v)])
	} else if len(v) < (1 << 8) {
		e.buf.WriteByte(0xD0)
		e.buf.WriteByte(byte(len(v)))
	} else if len(v) < (1 << 16) {
		e.buf.WriteByte(0xD1)
		binary.Write(e.buf, binary.BigEndian, uint16(len(v)))
	} else if len(v) < (1 << 32) {
		e.buf.WriteByte(0xD2)
		binary.Write(e.buf, binary.BigEndian, uint32(len(v)))
	} else {
		return errors.New("cannot encode string of length greater than 2,147,483,647 bytes")
	}

	e.buf.Write([]byte(v))

	return nil
}

func (e encoder) encodeMarshaller(v Marshaller) error {
	b, err := v.MarshalPackstream()
	if err != nil {
		return err
	}

	e.buf.Write(b)

	return nil
}
