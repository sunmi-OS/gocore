package robots

import (
	"context"
	"testing"
)

func TestRobot(t *testing.T) {
	robot := NewWithUrl("webhook-url")
	ctx := context.Background()
	err := robot.SendMarkdown(ctx, "title", "**hello**\nworld")
	if err != nil {
		t.Error(err)
	}
}
