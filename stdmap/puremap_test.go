package puremap

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

// func str(i int) string {
// 	return fmt.Sprintf("test string %d", i)
// }

// func strB(i int) string {
// 	return fmt.Sprintf("test-string %d", i)
// }

func str(i int) string {
	return fmt.Sprintf("%7d ", i)
}

func strB(i int) string {
	return fmt.Sprintf("%7d-", i)
}

func BenchmarkCacheMiss(b *testing.B) {
	for _, loadFactor := range []float64{0.2, 0.6, 1} {
		for _, size := range []int{1000, 10000, 100000, 1000000} {
			b.Run(fmt.Sprintf("Size=%d_Load=%.1f", size, loadFactor), func(b *testing.B) {
				numItems := int(float64(size) * loadFactor)
				cache := make(map[string]int, numItems)
				hitKeys := make([]string, numItems)

				for i := range numItems {
					cache[str(i)] = i
					hitKeys[i] = strB(i)
				}

				rand.Shuffle(len(hitKeys), func(i, j int) {
					hitKeys[i], hitKeys[j] = hitKeys[j], hitKeys[i]
				})

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					idx := i % numItems
					key := hitKeys[idx]
					_, found := cache[key]
					if found {
						b.Fatalf("Key %s should not be found", key)
					}
				}
			})
		}
	}
}

var sum uint64

func BenchmarkCacheHit(b *testing.B) {
	for _, loadFactor := range []float64{0.2, 0.6, 1} {
		for _, size := range []int{1000, 10000, 100000, 1000000} {
			b.Run(fmt.Sprintf("Size=%d_Load=%.1f", size, loadFactor), func(b *testing.B) {
				numItems := int(float64(size) * loadFactor)
				cache := make(map[string]int, numItems)
				hitKeys := make([]string, numItems)

				for i := range numItems {
					cache[str(i)] = i
					hitKeys[i] = str(i)
				}

				rand.Shuffle(len(hitKeys), func(i, j int) {
					hitKeys[i], hitKeys[j] = hitKeys[j], hitKeys[i]
				})

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					idx := i % numItems
					key := hitKeys[idx]
					_, found := cache[key]
					if !found {
						b.Fatalf("Key %s should be found", key)
					}
				}
			})
		}
	}
}

func BenchmarkCacheAddDel(b *testing.B) {
	for _, loadFactor := range []float64{0.2, 0.6, 1} {
		for _, size := range []int{1000, 10000, 100000, 1000000} {
			b.Run(fmt.Sprintf("Size=%d_Load=%.1f", size, loadFactor), func(b *testing.B) {
				numItems := int(float64(size) * loadFactor)
				cache := make(map[string]int, numItems)
				hitKeys := make([]string, numItems)

				for i := range numItems {
					cache[str(i)] = i
					hitKeys[i] = str(i)
				}

				rand.Shuffle(len(hitKeys), func(i, j int) {
					hitKeys[i], hitKeys[j] = hitKeys[j], hitKeys[i]
				})

				key := str(numItems)
				val := numItems

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					cache[key] = val
					delete(cache, key)
				}
			})
		}
	}
}

const fixedSeed1 = 12345
const fixedSeed2 = 76890

var cacheSizes = []int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000}

func BenchmarkCache2Hit(b *testing.B) {
	var size int
	for i := range cacheSizes {
		if cacheSizes[i] > size {
			size = cacheSizes[i]
		}
	}
	ss := make([]string, size)
	us := make([]string, size)
	for i := range size {
		ss[i] = str(i)
		us[i] = str(i)
	}
	for _, size := range cacheSizes {
		rng := rand.New(rand.NewPCG(fixedSeed1, fixedSeed2))
		rng.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("%8d", size), func(b *testing.B) {
			m := make(map[string]int, size)
			//m := make(map[string]int, size)
			for i := range size {
				m[ss[i]] = i
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := m[us[idx]]
				if !found {
					b.Fatalf("Key %s should be found", us[idx])
				}
			}
		})
	}
}

func BenchmarkCache2Miss(b *testing.B) {
	var size int
	for i := range cacheSizes {
		if cacheSizes[i] > size {
			size = cacheSizes[i]
		}
	}
	ss := make([]string, size)
	us := make([]string, size)
	for i := range size {
		ss[i] = str(i)
		us[i] = strB(i)
	}
	for _, size := range cacheSizes {
		rng := rand.New(rand.NewPCG(fixedSeed1, fixedSeed2))
		rng.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("%8d", size), func(b *testing.B) {
			m := make(map[string]int, size)
			for i := range size {
				m[ss[i]] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := m[us[idx]]
				if found {
					b.Fatalf("Key %s should not be found", us[idx])
				}
			}
		})
	}
}

func BenchmarkCache3Hit(b *testing.B) {
	var size int
	for i := range cacheSizes {
		if cacheSizes[i] > size {
			size = cacheSizes[i]
		}
	}
	ss := make([]string, size)
	us := make([]string, size)
	for i := range size {
		ss[i] = str(i)
		us[i] = str(i)
	}
	for _, size := range cacheSizes {
		rng := rand.New(rand.NewPCG(fixedSeed1, fixedSeed2))
		rng.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("%8d map", size), func(b *testing.B) {
			m := make(map[string]int, size)
			for i := range size {
				m[ss[i]] = i
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := m[us[idx]]
				if !found {
					b.Fatalf("Key %s should be found", us[idx])
				}
			}
		})
	}
}
