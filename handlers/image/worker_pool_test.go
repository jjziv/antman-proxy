package handlers

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool_Behavior(t *testing.T) {
	t.Parallel()

	t.Run("basic worker pool functionality", func(t *testing.T) {
		pool := NewWorkerPool(5)
		var counter int32

		// Submit 10 jobs
		for i := 0; i < 10; i++ {
			pool.Submit(func() {
				atomic.AddInt32(&counter, 1)
				time.Sleep(10 * time.Millisecond) // Simulate work
			})
		}

		pool.Wait()
		assert.Equal(t, int32(10), counter, "All jobs should be completed")
	})

	t.Run("concurrent job submission", func(t *testing.T) {
		pool := NewWorkerPool(3)
		var counter int32
		var wg sync.WaitGroup

		// Submit jobs concurrently
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pool.Submit(func() {
					atomic.AddInt32(&counter, 1)
					time.Sleep(10 * time.Millisecond)
				})
			}()
		}

		wg.Wait()
		pool.Wait()
		assert.Equal(t, int32(5), counter, "All concurrent jobs should be completed")
	})

	t.Run("worker pool under heavy load", func(t *testing.T) {
		pool := NewWorkerPool(2)
		var counter int32
		jobCount := 100

		start := time.Now()
		for i := 0; i < jobCount; i++ {
			pool.Submit(func() {
				atomic.AddInt32(&counter, 1)
				time.Sleep(1 * time.Millisecond)
			})
		}

		pool.Wait()
		duration := time.Since(start)

		assert.Equal(t, int32(jobCount), counter, "All jobs should be completed")
		// With 2 workers and 100 jobs taking 1ms each, it should take ~50ms
		// Add some buffer for scheduling overhead
		assert.Less(t, duration, 200*time.Millisecond, "Should complete within reasonable time")
	})

	t.Run("error handling in jobs", func(t *testing.T) {
		pool := NewWorkerPool(2)
		var errorCount int32
		var successCount int32

		// Submit jobs that may panic
		for i := 0; i < 10; i++ {
			i := i // Capture loop variable
			pool.Submit(func() {
				defer func() {
					if r := recover(); r != nil {
						atomic.AddInt32(&errorCount, 1)
					}
				}()

				if i%2 == 0 {
					panic("intentional panic")
				}
				atomic.AddInt32(&successCount, 1)
			})
		}

		pool.Wait()
		assert.Equal(t, int32(5), errorCount, "Half of the jobs should have panicked")
		assert.Equal(t, int32(5), successCount, "Half of the jobs should have succeeded")
	})

	t.Run("zero workers", func(t *testing.T) {
		assert.NotPanicsf(t, func() {
			NewWorkerPool(0)
		}, "Should not panic when creating pool with zero workers")
	})

	t.Run("worker pool shutdown", func(t *testing.T) {
		pool := NewWorkerPool(2)
		var counter int32
		completedChan := make(chan struct{})

		// Submit a long-running job
		pool.Submit(func() {
			atomic.AddInt32(&counter, 1)
			time.Sleep(100 * time.Millisecond)
			close(completedChan)
		})

		// Wait for job completion
		select {
		case <-completedChan:
			assert.Equal(t, int32(1), counter, "Job should complete")
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Job didn't complete in time")
		}
	})

	t.Run("sequential job ordering", func(t *testing.T) {
		pool := NewWorkerPool(1) // Single worker to guarantee order
		var results []int
		var mu sync.Mutex

		for i := 0; i < 5; i++ {
			i := i
			pool.Submit(func() {
				mu.Lock()
				results = append(results, i)
				mu.Unlock()
			})
		}

		pool.Wait()
		assert.Len(t, results, 5, "Should complete all jobs")
		for i := 0; i < len(results)-1; i++ {
			assert.Less(t, results[i], results[i+1], "Jobs should complete in order")
		}
	})
}

func BenchmarkWorkerPool(b *testing.B) {
	b.Run("small jobs", func(b *testing.B) {
		pool := NewWorkerPool(4)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			pool.Submit(func() {
				time.Sleep(1 * time.Microsecond)
			})
		}
		pool.Wait()
	})

	b.Run("medium jobs", func(b *testing.B) {
		pool := NewWorkerPool(4)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			pool.Submit(func() {
				time.Sleep(1 * time.Millisecond)
			})
		}
		pool.Wait()
	})

	b.Run("concurrent submission", func(b *testing.B) {
		pool := NewWorkerPool(4)
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				pool.Submit(func() {
					time.Sleep(100 * time.Microsecond)
				})
			}
		})
		pool.Wait()
	})
}
