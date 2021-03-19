package spinlock

import "testing"

func TestSpinLock(t *testing.T) {
	var mu SpinLock
	mu.Lock()
	defer mu.Unlock()
}

func BenchmarkSpinLock(b *testing.B) {
	var mu SpinLock

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			mu.Unlock()
		}
	})
}
