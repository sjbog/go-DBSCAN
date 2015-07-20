package dbscan

import (
	"sync"
	"testing"
)

func Test_Queue(t *testing.T) {
	t.Parallel()
	var (
		q              = NewConcurrentQueue_InsertOnly()
		waitGroup      = new(sync.WaitGroup)
		numSize        = uint(100)
		goroutinesSize = uint(10)
	)

	waitGroup.Add(int(goroutinesSize))

	for gi := uint(0); gi < goroutinesSize; gi += 1 {
		go func() {
			for i := uint(0); i < numSize; i += 1 {
				q.Add(i)
			}
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()

	if q.Size != uint64(numSize*goroutinesSize) {
		t.Fatalf("Didn't get the expected queue size : %v, result : %v", numSize*goroutinesSize, q.Size)
	}

	var (
		slice   = q.Slice()
		counter = make(map[uint]uint)
	)

	// t.Logf("len=%d cap=%d %v", len(slice), cap(slice), slice)

	if len(slice) != int(numSize*goroutinesSize) {
		t.Fatalf("Didn't get the expected slice size : %v, result : %v", numSize*goroutinesSize, len(slice))
	}

	for _, x := range slice {
		if _, ok := counter[x]; !ok {
			counter[x] = 1
		} else {
			counter[x] += 1
		}
	}
	for k, v := range counter {
		if v != goroutinesSize {
			t.Fatalf("Didn't get the expected count of \"%v\"'s, expected %v, got %v", k, numSize, v)
		}
	}
}
