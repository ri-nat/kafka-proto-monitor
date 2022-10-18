package monitor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ri-nat/kafka-proto-monitor/internal/pkg/config"
	"github.com/ri-nat/kafka-proto-monitor/internal/pkg/consumer"
	"github.com/ri-nat/kafka-proto-monitor/internal/pkg/parser"
)

// Monitor is a type for representing application instance.
// It holds options, connections and supporting objects.
type Monitor struct {
	Opts     *config.Options
	Consumer *kafka.Consumer
	Parser   *parser.Parser

	ShutdownCh chan (os.Signal)
}

// NewMonitor instantiates new `Monitor` object.
func NewMonitor() (monitor Monitor, err error) {
	opts, err := config.ParseOptions()
	if err != nil {
		return
	}

	monitor.Opts = opts
	return
}

// RunMonitor instantiates new `Monitor` object and calls `.Run()` on it.
func RunMonitor() error {
	monitor, err := NewMonitor()
	if err != nil {
		return err
	}

	return monitor.Run()
}

// Run initiates all the supporting objects, connections and starts consuming Kafka messages.
func (m *Monitor) Run() error {
	// Do not run monitor if starting with `--help` flag
	if m.Opts.IsHelp {
		return nil
	}

	// Handle Ctrl + C
	m.initInterruptHandler()

	// Initiate parser, parse proto file and look for a proto message in it
	if err := m.initParser(); err != nil {
		return err
	}

	// Init Kafka consumer
	if err := m.initConsumer(); err != nil {
		return err
	}

	// Start Kafka messages consumption loop
	m.loop()

	return nil
}

func (m *Monitor) loop() {
	fmt.Println("Starting consumption...")

	for {
		msg, err := m.Consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}

		if err := m.handleMessage(msg); err != nil {
			fmt.Printf("Message processing error: %v\n", err)
		}
	}
}

func (m *Monitor) handleMessage(msg *kafka.Message) error {
	protoMsgName, err := m.detectProtoMsgName(msg)
	if err != nil {
		return err
	}

	parsed, err := m.Parser.NewMessageFromBytes(protoMsgName, msg.Value)
	if err != nil {
		return nil
		// TODO: implement an option that will enable this error
		// return fmt.Errorf("kafka message unmarhaling error: %w", err)
	}

	fmt.Println(strings.Repeat("=", 64))
	fmt.Println()
	fmt.Printf("-- ID: %v\n", msg.TopicPartition)
	fmt.Printf("-- Timestamp: %v\n", msg.Timestamp)

	if len(m.Opts.ProtoMsgHeader) > 0 {
		fmt.Printf("-- Proto: %s\n", protoMsgName)
	}

	if m.Opts.PrintHeaders && msg.Headers != nil {
		fmt.Printf("-- Headers: %v\n", msg.Headers)
	}

	empJSON, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling proto to JSON: %w", err)
	}

	fmt.Printf("\n%s\n\n", string(empJSON))

	return nil
}

func (m *Monitor) detectProtoMsgName(msg *kafka.Message) (string, error) {
	protoMsgName := m.Opts.ProtoMsg

	if len(m.Opts.ProtoMsgHeader) > 0 {
		header, err := consumer.GetHeader(msg.Headers, m.Opts.ProtoMsgHeader)
		if err != nil {
			return "", err
		}

		protoMsgName = string(header.Value)
	}

	return protoMsgName, nil
}

func (m *Monitor) initConsumer() error {
	consumer, err := consumer.NewConsumer(m.Opts)
	if err != nil {
		return err
	}

	m.Consumer = consumer

	return nil
}

func (m *Monitor) initParser() error {
	parser, err := parser.NewParser(m.Opts)
	if err != nil {
		return err
	}

	m.Parser = parser

	return nil
}

func (m *Monitor) initInterruptHandler() {
	m.ShutdownCh = make(chan os.Signal)
	signal.Notify(m.ShutdownCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-m.ShutdownCh

		fmt.Println("")
		fmt.Println("Doing cleanup...")

		m.Consumer.Close()

		fmt.Println("Done!")

		os.Exit(0)
	}()
}
