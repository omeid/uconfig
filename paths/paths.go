package paths

// Kind represents the type of a path resolver.
type Kind int

const (
	Unknown Kind = iota
	Absolute
	Relative
	Workspace
)

func (k Kind) String() string {
	switch k {
	case Absolute:
		return "absolute"
	case Relative:
		return "relative"
	case Workspace:
		return "workspace"
	default:
		return "unknown"
	}
}

// One represents a resolver for a single file path.
type One interface {
	// Name returns the display name of the path (e.g., as written by the user).
	Name() string
	// Kind returns the kind of the path (e.g., paths.Absolute).
	Kind() Kind
	// Resolve returns the actual absolute filesystem path to the file.
	Resolve() string
}

// Set represents a resolver for a set of file paths.
type Set interface {
	// Name returns the display name of the path set (e.g., a glob pattern).
	Name() string
	// Kind returns the kind of the path set (e.g., paths.Relative).
	Kind() Kind
	// Resolve returns a slice of absolute filesystem paths matching the set.
	Resolve() ([]string, error)
}
