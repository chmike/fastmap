package altmapint

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func TestCacheAddGet(t *testing.T) {
	var c Cache
	c.Init()
	ss := []int{}
	for i := range 5000 {
		//t.Log(i)
		// if i == 7 {
		// 	print()
		// }
		_, ok := c.Add(i, i)
		if exp, got := false, ok; exp != got {
			t.Fatalf("%3d expect %v, got %v", i, exp, got)
		}

		ss = append(ss, i)

		v, ok := c.Get(i)
		if exp, got := true, ok; exp != got {
			t.Fatalf("%3d expect %v, got %v", i, exp, got)
		}
		if exp, got := i, v; exp != got {
			t.Fatalf("%3d expect %v, got %v", i, exp, got)
		}

		for j, s := range ss {
			if j == 7 {
				print()
			}
			_, ok := c.Get(s)
			if exp, got := true, ok; exp != got {
				t.Fatalf("%3d.%d for key %q expect %v, got %v", i, j, s, exp, got)
			}
		}
	}
}

func TestCacheAddDel(t *testing.T) {
	var c Cache
	c.Init()
	ss := []int{}
	for i := range 5000 {
		//t.Log(i)
		// if i == 7 {
		// 	print()
		// }
		_, ok := c.Add(i, i)
		if exp, got := false, ok; exp != got {
			t.Fatalf("%3d expect %v, got %v", i, exp, got)
		}

		ss = append(ss, i)
	}

	rand.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})

	for len(ss) > 0 {
		key := ss[len(ss)-1]
		ss = ss[:len(ss)-1]
		c.Del(key)

		for _, key := range ss {
			if _, ok := c.Get(key); !ok {
				t.Fatalf("failed to find key %q", key)
			}
		}
	}
	if c.Len() != 0 {
		t.Fatalf("expect empty, got %d", c.Len())
	}
}

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
			var c Cache
			c.Init()
			for i := range size {
				c.Add(ss[i], i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := c.Get(us[idx])
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
			var c Cache
			c.Init()
			for i := range size {
				c.Add(ss[i], i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				idx := i % size
				_, found := c.Get(us[idx])
				if found {
					b.Fatalf("Key %v should not be found", us[idx])
				}
			}
		})
	}
}
