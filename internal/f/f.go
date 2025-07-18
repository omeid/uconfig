// Package f provides simple test fixtures for uconfig.
package f

import (
	"encoding"
	"strings"
	"time"
)

// Anon is part of text fixtures.
type Anon struct {
	Version string
}

// Host is part of text fixtures.
type Host struct {
	Address string
	Port    string
}

// RethinkConfig is part of text fixtures.
type RethinkConfig struct {
	Host     Host
	Db       string `default:"primary" usage:"main database used by our application"`
	Password string `secret:""`
}

// Redis is part of text fixtures.
type Redis struct {
	Host string `uconfig:".Address"`
	Port int
}

// Config is part of text fixtures.
type Config struct {
	Command string `flag:",command" default:"run"` // expose this as the cli command.
	Anon
	GoHard  bool
	Redis   Redis
	Rethink RethinkConfig
}

// TextUnmarshalerStringSlice is an example of encoding.TextUnmarshaler
type TextUnmarshalerStringSlice []string

// UnmarshalText is part of encoding.TextUnmarshaler
func (l *TextUnmarshalerStringSlice) UnmarshalText(value []byte) error {
	*l = strings.Split(string(value), ".")
	return nil
}

type ReadableDirection int

func NewReadableDirection(value int) *ReadableDirection {
	dir := ReadableDirection(value)
	return &dir
}

// UnmarshalText is part of encoding.TextUnmarshaler
func (l *ReadableDirection) UnmarshalText(value []byte) error {
	switch string(value) {
	case "north":
		*l = ReadableDirection(0)
	case "east":
		*l = ReadableDirection(90)
	case "south":
		*l = ReadableDirection(180)
	case "west":
		*l = ReadableDirection(270)
	default:
		*l = ReadableDirection(0)
	}

	return nil
}

type (
	ElemUnmarshalerSlice    []ReadableDirection
	ElemPtrUnmarshalerSlice []*ReadableDirection
)

// ensure the interfae is implemented properly.
var (
	_ encoding.TextUnmarshaler = &TextUnmarshalerStringSlice{}
	_ encoding.TextUnmarshaler = NewReadableDirection(0)
)

// Types is part of text fixtures.
type Types struct {
	String   string
	Bool     bool
	Duration time.Duration
	Int      int
	Int8     int8
	Int16    int16
	Int32    int32
	Int64    int64
	Uint     uint
	Uint8    uint8
	Uint16   uint16
	Uint32   uint32
	Uint64   uint64
	Float32  float32
	Float64  float64

	// MapStringInt    map[string]int

	SliceString   []string
	SliceInt      []int
	SliceInt32    []int
	SliceUint     []uint
	SliceFloat32  []float32
	SliceDuration []time.Duration

	SliceTextUnmarshaler    *TextUnmarshalerStringSlice
	SliceElemUnmarshaler    ElemUnmarshalerSlice
	SliceElemPtrUnmarshaler ElemPtrUnmarshalerSlice
}
