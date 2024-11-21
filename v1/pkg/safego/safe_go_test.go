package safego

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSafeGoWithLoggerLoopCapture(t *testing.T) {
	var res []int
	var wg sync.WaitGroup
	numberOfTimes := 500
	wg.Add(numberOfTimes)
	var mu sync.Mutex
	for i := 0; i < numberOfTimes; i++ {
		i := i
		SafeGo(func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			mu.Lock()
			res = append(res, i)
			mu.Unlock()
		})
	}
	wg.Wait()
	assert.Contains(t, res, 0)
	assert.Contains(t, res, 1)
	assert.Contains(t, res, 2)
	assert.Contains(t, res, 4)
	assert.Equal(t, numberOfTimes, len(res))
}
func panicUtil() {
	panic("panicking here...")
}
func TestIfHandlingPanics(t *testing.T) {
	var res []int
	var wg sync.WaitGroup
	wg.Add(1)
	SafeGo(func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		panicUtil()
		res = append(res, 1)
	})
	wg.Wait()
	assert.Equal(t, len(res), 0)
}
func TestSafeGoWithLoggerLoopCapturePanic(t *testing.T) {
	var res []int
	var wg sync.WaitGroup
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		SafeGo(func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			if i == 4 {
				panic("")
			}
			ch <- i
		})
	}
	wg.Wait()
	close(ch)
	for i := range ch {
		res = append(res, i)
	}
	assert.Contains(t, res, 0)
	assert.Contains(t, res, 1)
	assert.Contains(t, res, 2)
	assert.Contains(t, res, 3)
}
