package f

//Anon is part of text fixtures.
type Anon struct {
	Version string
}

//Host is part of text fixtures.
type Host struct {
	Address string
	Port    string
}

//RethinkConfig is part of text fixtures.
type RethinkConfig struct {
	Host Host
	Db   string
}

//Redis is part of text fixtures.
type Redis struct {
	Host string
	Port int
}

//Config is part of text fixtures.
type Config struct {
	Anon
	GoHard  bool
	Redis   Redis
	Rethink RethinkConfig
}

//Types is part of text fixtures.
type Types struct {
	String  string
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64

	// MapStringInt    map[string]int

	// SliceString []string
	// SliceInt    []int
}
