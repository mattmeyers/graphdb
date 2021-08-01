package packstream

import (
	"bytes"
	"errors"
)

// structMarkers is a quick lookup table for structs.
var structMarkers = [16]byte{
	0xB0, 0xB1, 0xB2, 0xB3,
	0xB4, 0xB5, 0xB6, 0xB7,
	0xB8, 0xB9, 0xBA, 0xBB,
	0xBC, 0xBD, 0xBE, 0xBF,
}

func writeStructHeader(buf *bytes.Buffer, s Structure) error {
	if s.FieldCount() > (1 << 4) {
		return errors.New("cannot encode structure with more than 16 fields")
	}

	buf.WriteByte(0x80 + byte(s.FieldCount()))
	buf.WriteByte(s.Tag())

	return nil
}

type Structure interface {
	Tag() byte
	FieldCount() uint
}

type Node struct {
	ID         int
	Labels     List
	Properties Dictionary
}

func (n Node) Tag() byte        { return 0x4E }
func (n Node) FieldCount() uint { return 3 }

func (n Node) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, n); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, n.Labels); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, n.Properties); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
