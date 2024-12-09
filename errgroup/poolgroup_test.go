package errgroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/sunmi-OS/gocore/v2/worker"
)

func TestPoolNormal(t *testing.T) {
	m := make(map[int]int)
	for i := 0; i < 4; i++ {
		m[i] = i
	}
	eg := WithPoolContext(context.Background(), worker.P())
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

func TestPoolRecover(t *testing.T) {
	eg := WithPoolContext(context.Background(), worker.P())
	eg.Go(func(context.Context) (err error) {
		panic("oh my god!")
	})
	if err := eg.Wait(); err != nil {
		t.Log(err)
		return
	}
	t.FailNow()
}

// JustErrors illustrates the use of a Group in place of a sync.WaitGroup to
// simplify goroutine counting and error handling. This example is derived from
// the sync.WaitGroup example at https://golang.org/pkg/sync/#example_WaitGroup.
func ExamplePoolGroup_justErrors() {
	eg := WithPoolContext(context.Background(), worker.P())
	urls := []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, _url := range urls {
		// Launch a goroutine to fetch the URL.
		url := _url // https://golang.org/doc/faq#closures_and_goroutines
		eg.Go(func(context.Context) error {
			// Fetch the URL.
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}
	// Wait for all HTTP fetches to complete.
	if err := eg.Wait(); err == nil {
		fmt.Println("Successfully fetched all URLs.")
	}
}

// Parallel illustrates the use of a Group for synchronizing a simple parallel
// task: the "Google Search 2.0" function from
// https://talks.golang.org/2012/concurrency.slide#46, augmented with a Context
// and error-handling.
func ExamplePoolGroup_parallel() {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		eg := WithPoolContext(ctx, worker.P())

		searches := []Search{Web, Image, Video}
		results := make([]Result, len(searches))
		for _i, _search := range searches {
			i, search := _i, _search // https://golang.org/doc/faq#closures_and_goroutines
			eg.Go(func(context.Context) error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := eg.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}

	// Output:
	// web result for "golang"
	// image result for "golang"
	// video result for "golang"
}

func TestPoolZeroGroup(t *testing.T) {
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
		eg := WithPoolContext(context.Background(), worker.P())

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

func TestPoolWithCancel(t *testing.T) {
	eg := WithPoolContext(context.Background(), worker.P())
	eg.Go(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		return fmt.Errorf("boom")
	})
	var doneErr error
	eg.Go(func(ctx context.Context) error {
		<-ctx.Done()
		doneErr = ctx.Err()
		return doneErr
	})
	_ = eg.Wait()
	if doneErr != context.Canceled {
		t.Error("error should be Canceled")
	}
}
