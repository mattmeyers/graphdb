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
	case float32:
		err = encodeFloat(e.buf, float64(val))
	case float64:
		err = encodeFloat(e.buf, val)
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

// encodeInt provides an optimal int64 encoding based on the packstream spec.
// The encoding table is as such:
// 	Minimum                         Maximum                     Optimal Representation
// 	-9_223_372_036_854_775_808      -2_147_483_649              INT_64
// 	-2_147_483_648                  -32_769                     INT_32
// 	-32_768                         -129                        INT_16
// 	-128                            -17                         INT_8
// 	-16                             +127                        TINY_INT
// 	+128                            +32_767                     INT_16
// 	+32_768                         +2_147_483_647              INT_32
// 	+2_147_483_648                  +9_223_372_036_854_775_807  INT_64
func encodeInt(buf *bytes.Buffer, v int) error {
	if -16 <= v && v <= 127 { // TINY_INT
		buf.WriteByte(byte(v))
	} else if -128 <= v && v <= -17 { // INT_8
		buf.WriteByte(0xC8)
		binary.Write(buf, binary.BigEndian, byte(v))
	} else if -32_768 <= v && v <= 32_767 { // INT_16
		buf.WriteByte(0xC9)
		binary.Write(buf, binary.BigEndian, int16(v))
	} else if -2_147_483_648 <= v && v <= 2_147_483_647 { // INT_32
		buf.WriteByte(0xCA)
		binary.Write(buf, binary.BigEndian, int32(v))
	} else if -9_223_372_036_854_775_808 <= v && v <= 9_223_372_036_854_775_807 { // INT_64
		buf.WriteByte(0xCB)
		binary.Write(buf, binary.BigEndian, int64(v))
	} else {
		return errors.New("unable to encode int")
	}

	return nil
}

func encodeFloat(buf *bytes.Buffer, v float64) error {
	buf.WriteByte(0xC1)
	binary.Write(buf, binary.BigEndian, v)

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
