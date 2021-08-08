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

func (Node) Tag() byte        { return 0x4E }
func (Node) FieldCount() uint { return 3 }

func (n Node) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, n); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, n.ID); err != nil {
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

type Relationship struct {
	ID          int
	StartNodeID int
	EndNodeID   int
	Type        string
	Properties  Dictionary
}

func (Relationship) Tag() byte        { return 0x52 }
func (Relationship) FieldCount() uint { return 5 }

func (r Relationship) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, r); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, r.ID); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, r.StartNodeID); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, r.EndNodeID); err != nil {
		return nil, err
	}

	if err := encodeString(buf, r.Type); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, r.Properties); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type UnboundRelationship struct {
	ID         int
	Type       string
	Properties Dictionary
}

func (UnboundRelationship) Tag() byte        { return 0x72 }
func (UnboundRelationship) FieldCount() uint { return 3 }

func (r UnboundRelationship) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, r); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, r.ID); err != nil {
		return nil, err
	}

	if err := encodeString(buf, r.Type); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, r.Properties); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Path struct {
	// Nodes is alist of nodes.
	Nodes List
	// Rels is a list of unbound relationships.
	Rels List
	// IDs is a list of relationship id and node id to represent the path.
	IDs List
}

func (Path) Tag() byte        { return 0x50 }
func (Path) FieldCount() uint { return 3 }

func (p Path) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, p); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, p.Nodes); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, p.Rels); err != nil {
		return nil, err
	}

	if err := encodeMarshaller(buf, p.IDs); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Date struct {
	Days int
}

func (Date) Tag() byte        { return 0x44 }
func (Date) FieldCount() uint { return 1 }

func (d Date) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, d); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, d.Days); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Time struct {
	Nanoseconds     int
	TZOffsetSeconds int
}

func (Time) Tag() byte        { return 0x54 }
func (Time) FieldCount() uint { return 2 }

func (t Time) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, t); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Nanoseconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.TZOffsetSeconds); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t Time) ToUTCNanoseconds() int {
	return t.Nanoseconds - (t.TZOffsetSeconds * 1_000_000_000)
}

type LocalTime struct {
	Nanoseconds int
}

func (LocalTime) Tag() byte        { return 0x74 }
func (LocalTime) FieldCount() uint { return 1 }

func (t LocalTime) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, t); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Nanoseconds); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type DateTime struct {
	Seconds         int
	Nanoseconds     int
	TZOffsetSeconds int
}

func (DateTime) Tag() byte        { return 0x46 }
func (DateTime) FieldCount() uint { return 3 }

func (t DateTime) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, t); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Seconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Nanoseconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.TZOffsetSeconds); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t DateTime) ToUTCNanoseconds() int {
	return (t.Seconds * 1_000_000_000) + t.Nanoseconds - (t.TZOffsetSeconds * 1_000_000_000)
}

type DateTimeZoneID struct {
	Seconds     int
	Nanoseconds int
	TimeZoneID  string
}

func (DateTimeZoneID) Tag() byte        { return 0x66 }
func (DateTimeZoneID) FieldCount() uint { return 3 }

func (t DateTimeZoneID) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, t); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Seconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Nanoseconds); err != nil {
		return nil, err
	}

	if err := encodeString(buf, t.TimeZoneID); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type LocalDateTime struct {
	Seconds     int
	Nanoseconds int
}

func (LocalDateTime) Tag() byte        { return 0x64 }
func (LocalDateTime) FieldCount() uint { return 2 }

func (t LocalDateTime) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, t); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Seconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, t.Nanoseconds); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Duration struct {
	Months      int
	Days        int
	Seconds     int
	Nanoseconds int
}

func (Duration) Tag() byte        { return 0x45 }
func (Duration) FieldCount() uint { return 4 }

func (d Duration) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, d); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, d.Months); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, d.Days); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, d.Seconds); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, d.Nanoseconds); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Point2D struct {
	SRID int
	X    float64
	Y    float64
}

func (Point2D) Tag() byte        { return 0x58 }
func (Point2D) FieldCount() uint { return 3 }

func (p Point2D) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, p); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, p.SRID); err != nil {
		return nil, err
	}

	if err := encodeFloat(buf, p.X); err != nil {
		return nil, err
	}

	if err := encodeFloat(buf, p.Y); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Point3D struct {
	SRID int
	X    float64
	Y    float64
	Z    float64
}

func (Point3D) Tag() byte        { return 0x59 }
func (Point3D) FieldCount() uint { return 4 }

func (p Point3D) MarshalPackstream() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := writeStructHeader(buf, p); err != nil {
		return nil, err
	}

	if err := encodeInt(buf, p.SRID); err != nil {
		return nil, err
	}

	if err := encodeFloat(buf, p.X); err != nil {
		return nil, err
	}

	if err := encodeFloat(buf, p.Y); err != nil {
		return nil, err
	}

	if err := encodeFloat(buf, p.Z); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
