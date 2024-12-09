package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	now := time.Now()
	err := Retry(context.Background(), func(ctx context.Context) error {
		fmt.Println("Retry...")
		return errors.New("something wrong")
	}, 3, time.Second)
	assert.NotNil(t, err)
	assert.Equal(t, 2, int(time.Since(now).Seconds()))
}
