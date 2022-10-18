package consumer

import (
	"regexp"
	"testing"
)

func TestGenerateConsumerGroupID(t *testing.T) {
	groupID := generateConsumerGroupID()
	// Looking for a string like `kpm-user-hostname-1777070123`
	match, _ := regexp.MatchString("kpm\\-.+\\-.+\\-\\d{10,}", groupID)
	if !match {
		t.Errorf("Generated consumer ID didn't matched the pattern: %v", groupID)
	}
}
