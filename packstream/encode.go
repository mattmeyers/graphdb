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
		encodeBool(e.buf, val)
	case *bool:
		encodeBool(e.buf, *val)
	case []byte:
		encodeBytes(e.buf, val)
	case string:
		err = encodeString(e.buf, val)
	case *string:
		err = encodeString(e.buf, *val)
	case int:
		err = encodeInt(e.buf, val)
	case int8:
		err = encodeInt(e.buf, int(val))
	case int16:
		err = encodeInt(e.buf, int(val))
	case int32:
		err = encodeInt(e.buf, int(val))
	case int64:
		err = encodeInt(e.buf, int(val))
	case Marshaller:
		err = encodeMarshaller(e.buf, val)
	default:
		err = fmt.Errorf("unable to marshal value of type %T", v)
	}

	return err
}

func encodeBool(buf *bytes.Buffer, v bool) {
	if v {
		buf.WriteByte(0xC3)
	} else {
		buf.WriteByte(0xC2)
	}
}

func encodeBytes(buf *bytes.Buffer, v []byte) error {
	if len(v) < (1 << 8) {
		buf.WriteByte(0xCC)
		buf.WriteByte(byte(len(v)))
	} else if len(v) < (1 << 16) {
		buf.WriteByte(0xCD)
		binary.Write(buf, binary.BigEndian, uint16(len(v)))
	} else if len(v) < (1 << 32) {
		buf.WriteByte(0xCE)
		binary.Write(buf, binary.BigEndian, uint32(len(v)))
	} else {
		return errors.New("cannot encode byte slices of length greater than 2,147,483,647 bytes")
	}

	buf.Write(v)

	return nil
}

func encodeInt(buf *bytes.Buffer, v int) error {
	if -16 <= v && v <= 128 {
		buf.WriteByte(byte(v))
	} else {
		return errors.New("unable to encode int")
	}

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

func encodeString(buf *bytes.Buffer, v string) error {
	if len(v) < (1 << 4) {
		buf.WriteByte(shortStringMarkers[len(v)])
	} else if len(v) < (1 << 8) {
		buf.WriteByte(0xD0)
		buf.WriteByte(byte(len(v)))
	} else if len(v) < (1 << 16) {
		buf.WriteByte(0xD1)
		binary.Write(buf, binary.BigEndian, uint16(len(v)))
	} else if len(v) < (1 << 32) {
		buf.WriteByte(0xD2)
		binary.Write(buf, binary.BigEndian, uint32(len(v)))
	} else {
		return errors.New("cannot encode string of length greater than 2,147,483,647 bytes")
	}

	buf.Write([]byte(v))

	return nil
}

func encodeMarshaller(buf *bytes.Buffer, v Marshaller) error {
	b, err := v.MarshalPackstream()
	if err != nil {
		return err
	}

	buf.Write(b)

	return nil
}
