package puremapint

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

const fixedSeed1 = 12345
const fixedSeed2 = 76890

var cacheSizes = []int{1, 10, 100, 1000, 10000, 100000, 1000000, 10000000}

func BenchmarkCache2Hit(b *testing.B) {
	size := cacheSizes[len(cacheSizes)-1]
	ss := make([]int, size)
	us := make([]int, size)
	for i := range size {
		ss[i] = i
		us[i] = i
	}
	for _, size := range cacheSizes {
		rng := rand.New(rand.NewPCG(fixedSeed1, fixedSeed2))
		rng.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("%8d", size), func(b *testing.B) {
			m := make(map[int]int, size)
			//m := make(map[string]int, size)
			for i := range size {
				m[ss[i]] = i
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := m[us[idx]]
				if !found {
					b.Fatalf("Key %v should be found", us[idx])
				}
			}
		})
	}
}

func BenchmarkCache2Miss(b *testing.B) {
	size := cacheSizes[len(cacheSizes)-1]
	ss := make([]int, size)
	us := make([]int, size)
	for i := range size {
		ss[i] = i
		us[i] = i + size
	}
	for _, size := range cacheSizes {
		rng := rand.New(rand.NewPCG(fixedSeed1, fixedSeed2))
		rng.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("%8d", size), func(b *testing.B) {
			m := make(map[int]int, size)
			for i := range size {
				m[ss[i]] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := m[us[idx]]
				if found {
					b.Fatalf("Key %v should not be found", us[idx])
				}
			}
		})
	}
}
