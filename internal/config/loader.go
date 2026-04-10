package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Flags holds the parsed CLI flags that can override config file values.
type Flags struct {
	ConfigPath string
	StartPort  int
	EndPort    int
	Interval   time.Duration
	LogFile    string
}

// ParseFlags parses os.Args and returns a Flags struct.
func ParseFlags() *Flags {
	f := &Flags{}
	flag.StringVar(&f.ConfigPath, "config", "", "path to JSON config file")
	flag.IntVar(&f.StartPort, "start", 0, "start of port range (overrides config)")
	flag.IntVar(&f.EndPort, "end", 0, "end of port range (overrides config)")
	flag.DurationVar(&f.Interval, "interval", 0, "scan interval (overrides config)")
	flag.StringVar(&f.LogFile, "log", "", "log output file (overrides config)")
	flag.Parse()
	return f
}

// Resolve builds a final Config by loading the file (if given) and applying
// any non-zero flag overrides on top.
func Resolve(f *Flags) (*Config, error) {
	var cfg *Config
	var err error

	if f.ConfigPath != "" {
		cfg, err = Load(f.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("config: %w", err)
		}
	} else {
		cfg = DefaultConfig()
	}

	if f.StartPort != 0 {
		cfg.StartPort = f.StartPort
	}
	if f.EndPort != 0 {
		cfg.EndPort = f.EndPort
	}
	if f.Interval != 0 {
		cfg.Interval = f.Interval
	}
	if f.LogFile != "" {
		cfg.LogFile = f.LogFile
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: invalid configuration: %v\n", err)
		return nil, err
	}
	return cfg, nil
}
