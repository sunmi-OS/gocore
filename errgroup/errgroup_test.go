package errgroup

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormal(t *testing.T) {
	m := make(map[int]int)
	for i := range 4 {
		m[i] = i
	}
	eg := WithContext(context.Background())
	eg.Go(func(context.Context) (err error) {
		m[1]++
		return
	})
	eg.Go(func(context.Context) (err error) {
		m[2]++
		return
	})
	if err := eg.Wait(); err != nil {
		t.Log(err)
	}
	t.Log(m)
}

func sleep1s(context.Context) error {
	time.Sleep(time.Second)
	return nil
}

func TestGOMAXPROCS(t *testing.T) {
	ctx := context.Background()

	// 没有并发数限制
	eg := WithContext(ctx)
	now := time.Now()
	eg.Go(sleep1s)
	eg.Go(sleep1s)
	eg.Go(sleep1s)
	eg.Go(sleep1s)
	err := eg.Wait()
	assert.Nil(t, err)
	sec := math.Round(time.Since(now).Seconds())
	if sec != 1 {
		t.FailNow()
	}

	// 限制并发数
	eg2 := WithContext(ctx)
	eg2.GOMAXPROCS(2)
	now = time.Now()
	eg2.Go(sleep1s)
	eg2.Go(sleep1s)
	eg2.Go(sleep1s)
	eg2.Go(sleep1s)
	err = eg2.Wait()
	assert.Nil(t, err)
	sec = math.Round(time.Since(now).Seconds())
	if sec != 2 {
		t.FailNow()
	}

	// context canceled
	eg3 := WithContext(ctx)
	eg3.GOMAXPROCS(2)
	eg3.Go(func(context.Context) error {
		return errors.New("error for testing errgroup context")
	})
	eg3.Go(func(ctx context.Context) error {
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			t.Log("caused by", context.Cause(ctx))
		default:
		}
		return nil
	})
	err = eg3.Wait()
	assert.NotNil(t, err)
	t.Log(err)
}

func TestRecover(t *testing.T) {
	eg := WithContext(context.Background())
	eg.Go(func(context.Context) (err error) {
		panic("oh my god!")
	})
	if err := eg.Wait(); err != nil {
		t.Log(err)
		return
	}
	t.FailNow()
}

func TestZeroGroup(t *testing.T) {
	err1 := errors.New("errgroup_test: 1")
	err2 := errors.New("errgroup_test: 2")

	cases := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{nil}},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		eg := WithContext(context.Background())

		var firstErr error
		for i, err := range tc.errs {
			err := err
			eg.Go(func(context.Context) error { return err })

			if firstErr == nil && err != nil {
				firstErr = err
			}

			if gErr := eg.Wait(); gErr != firstErr {
				t.Errorf("after g.Go(func() error { return err }) for err in %v\n"+
					"g.Wait() = %v; want %v", tc.errs[:i+1], err, firstErr)
			}
		}
	}
}
