package puremap

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func str(i int) string {
	return fmt.Sprintf("%7d ", i)
}

func strB(i int) string {
	return fmt.Sprintf("%7d-", i)
}

const fixedSeed1 = 12345
const fixedSeed2 = 76890

var cacheSizes = []int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000}

func BenchmarkCache2Hit(b *testing.B) {
	size := cacheSizes[len(cacheSizes)-1]
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
	size := cacheSizes[len(cacheSizes)-1]
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
