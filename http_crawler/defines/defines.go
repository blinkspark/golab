package defines

const (
	SigTerminate = iota
)

type SaveRule struct {
	ContentType string
}

type DigRule struct {
	Deep int32
}

// Target used in Config struct
type Target struct {
	URL      string
	SaveRule SaveRule
	DigRule  DigRule
}

// Config is a struct to read config file
type Config struct {
	Proxy      string
	SavePath   string
	Retry      int
	MaxRoutine int32
	Targets    []Target
}
