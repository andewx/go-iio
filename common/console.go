package common

// ConsoleOptions are used to control display options in console
type ConsoleOptions struct {
	AttrVal  bool
	AttrName bool
	Host     string
	Device   string
}

// NewConsoleOptions creates a new ConsoleOptions struct
func NewConsoleOptions() *ConsoleOptions {
	return &ConsoleOptions{
		AttrVal:  false,
		AttrName: false,
		Host:     "192.168.2.1",
		Device:   "ad9361",
	}
}
