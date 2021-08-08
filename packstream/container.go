package packstream

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// List is a heterogeneous sequence of values.
//
// Lists can contain up to 2,147,483,647 items (1 << 32). If more items exist
// in the list at the time of encoding, then an error will occur.
type List []interface{}

// shortListMarkers is a quick lookup table for lists of length less than
// 16 bytes.
var shortListMarkers = [16]byte{
	0x90, 0x91, 0x92, 0x93,
	0x94, 0x95, 0x96, 0x97,
	0x98, 0x99, 0x9A, 0x9B,
	0x9C, 0x9D, 0x9E, 0x9F,
}

func (l List) writeMarker(buf *bytes.Buffer) (err error) {
	if len(l) < (1 << 4) {
		buf.WriteByte(shortListMarkers[len(l)])
	} else if len(l) < (1 << 8) {
		buf.WriteByte(0xD4)
		buf.WriteByte(byte(len(l)))
	} else if len(l) < (1 << 16) {
		buf.WriteByte(0xD5)
		binary.Write(buf, binary.BigEndian, uint16(len(l)))
	} else if len(l) < (1 << 32) {
		buf.WriteByte(0xD6)
		binary.Write(buf, binary.BigEndian, uint32(len(l)))
	} else {
		err = errors.New("cannot encode List with more than 2,147,483,647 elements")
	}

	return err
}

func (l List) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := l.writeMarker(buf); err != nil {
		return nil, err
	}

	for _, item := range l {
		b, err := Marshal(item)
		if err != nil {
			return nil, err
		}

		buf.Write(b)
	}

	return buf.Bytes(), nil
}

type Dictionary map[string]interface{}

// shortDictionaryMarkers is a quick lookup table for dictionaries of length less than
// 16 bytes.
var shortDictionaryMarkers = [16]byte{
	0xA0, 0xA1, 0xA2, 0xA3,
	0xA4, 0xA5, 0xA6, 0xA7,
	0xA8, 0xA9, 0xAA, 0xAB,
	0xAC, 0xAD, 0xAE, 0xAF,
}

func (d Dictionary) writeMarker(buf *bytes.Buffer) (err error) {
	if len(d) < (1 << 4) {
		buf.WriteByte(shortDictionaryMarkers[len(d)])
	} else if len(d) < (1 << 8) {
		buf.WriteByte(0xD8)
		buf.WriteByte(byte(len(d)))
	} else if len(d) < (1 << 16) {
		buf.WriteByte(0xD9)
		binary.Write(buf, binary.BigEndian, uint16(len(d)))
	} else if len(d) < (1 << 32) {
		buf.WriteByte(0xDA)
		binary.Write(buf, binary.BigEndian, uint32(len(d)))
	} else {
		err = errors.New("cannot encode Dictionary with more than 2,147,483,647 key-value pairs")
	}

	return err
}

func (d Dictionary) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := d.writeMarker(buf); err != nil {
		return nil, err
	}

	for k, v := range d {
		ek, err := Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(ek)

		ev, err := Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(ev)
	}

	return buf.Bytes(), nil
}
