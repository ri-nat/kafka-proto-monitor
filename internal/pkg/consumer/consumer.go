package consumer

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ri-nat/kafka-proto-monitor/internal/pkg/config"
)

func NewConsumer(opts *config.Options) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(newConfigMap(opts))
	if err != nil {
		return nil, fmt.Errorf("consumer init error: %w", err)
	}

	if err := c.SubscribeTopics(opts.Topics, nil); err != nil {
		return nil, fmt.Errorf("topics subscribing error: %w", err)
	}

	return c, nil
}

// GetHeader returns a header out of headers slice by given key
func GetHeader(headers []kafka.Header, key string) (*kafka.Header, error) {
	for _, header := range headers {
		if header.Key == key {
			return &header, nil
		}
	}

	return nil, fmt.Errorf("header `%s` not found", key)
}

func newConfigMap(opts *config.Options) *kafka.ConfigMap {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(opts.Brokers, ","),
		"group.id":          generateConsumerGroupID(),
		"auto.offset.reset": "latest",
	}

	if opts.ReadEarliest {
		fmt.Println("Reading from the earliest (oldest) messages")
		_ = config.SetKey("auto.offset.reset", "earliest")
	} else {
		fmt.Println("Reading from the latest (newest) messages")
	}

	return config
}

// Produces IDs like `kpm-user-hostname-1777070123`
func generateConsumerGroupID() string {
	groupIDParts := []string{"kpm"}

	user, err := user.Current()
	username := "unknown"
	if err == nil {
		username = user.Username
	}
	groupIDParts = append(groupIDParts, username)

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	groupIDParts = append(groupIDParts, hostName)

	groupIDParts = append(groupIDParts, fmt.Sprintf("%d", time.Now().Unix()))

	return strings.Join(groupIDParts, "-")
}
