package domain

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func BenchmarkEngine_LLeftPush(b *testing.B) {
	engine := NewEngine()
	key := "a"
	value := []string{"1"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = engine.LLeftPush(key, value)
	}
}

func TestEngine_LLeftPushConcurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("concurrency-related")
	}

	a := assert.New(t)

	engine := NewEngine()
	key := "a"
	value := []string{"1"}

	numGoroutines := 10
	numIterations := 1000

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				_, _ = engine.LLeftPush(key, value)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	values, err := engine.LRange(key, 0, -1)
	if err != nil {
		a.Errorf(err, "LRange failed")
	}

	wantValuesLen := numGoroutines * numIterations

	a.Equal(wantValuesLen, len(values))
}
