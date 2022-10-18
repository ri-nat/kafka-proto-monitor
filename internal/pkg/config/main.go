package config

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	Brokers        []string `short:"b" long:"broker" description:"Kafka broker URL (you can use this option multiple times)" default:"localhost:9092"`
	Topics         []string `short:"t" long:"topic" description:"Kafka topic to read from (you can use this option multiple times)" required:"true"`
	ReadEarliest   bool     `short:"a" long:"read-earliest" description:"Read topic starting from the beginning (false by default)"`
	PrintHeaders   bool     `short:"r" long:"print-headers" description:"Print message headers (false by default)"`
	ProtoFile      string   `short:"p" long:"proto-file" description:"Path to proto file" required:"true"`
	ProtoMsg       string   `short:"m" long:"proto-message" description:"Proto message to use"`
	ProtoMsgHeader string   `short:"e" long:"proto-message-header" description:"Name of Kafka message header, that contains message's proto name"`

	// This flag is used only to detect if the app is running in `--help` mode or not
	IsHelp bool
}

func ParseOptions() (*Options, error) {
	opts := &Options{}

	_, err := flags.Parse(opts)
	if err != nil {
		// Check for a special error type to detect `--help` mode
		if err.(*flags.Error).Type == flags.ErrHelp {
			opts.IsHelp = true
		} else {
			return nil, fmt.Errorf("arguments parsing error: %w", err)
		}
	}

	return opts, nil
}
