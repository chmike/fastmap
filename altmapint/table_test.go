package altmapint

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
)

func TestTableAddGet(t *testing.T) {
	seed := MakeSeed()
	//seed = 0 // for debugging

	c := newTable(0)
	var i int
	for i = range 8192 {
		// t.Log(i)
		// if i == 283 {
		// 	print()
		// }
		key := i
		hash := seed.Hash(key)
		if !c.add(key, i, hash) {
			break
		}
		if _, ok := c.get(key, hash); !ok {
			t.Fatalf("%3d failed to find key %v", i, key)
		}

		for j := range c.nItems {
			key := i
			_, ok := c.get(key, seed.Hash(key))
			if !ok {
				t.Fatalf("%3d.%d could not find key %q", i, j, key)
			}
		}
	}
	t.Logf("reached %d, %d%%", i, c.occupancy())
	var buf strings.Builder
	buf.WriteString("\n")
	for i := range c.groups[:tableSize] {
		count := c.groups[i].header.FindUsed().Len()
		occupancy := float64(count*100) / 8.
		buf.WriteString(fmt.Sprintf("%4d: %d %5.1f%%\n", i, count, occupancy))
	}
	t.Log(buf.String())

	// check that the iterator see all keys once
	m := make(map[int]struct{}, c.nItems)
	var count int
	for _, v := range c.items() {
		m[v] = struct{}{}
		count++
	}
	if len(m) != int(c.nItems) || count != int(c.nItems) {
		t.Fatalf("expect %d, found %d different keys and count %d", c.nItems, len(m), count)
	}
}

func TestTableAddDel(t *testing.T) {
	seed := MakeSeed()
	//seed = 0 // for debugging

	c := newTable(0)
	var keys []int
	for i := range tableItems {
		key := i
		if !c.add(key, i, seed.Hash(key)) {
			break
		}
		keys = append(keys, key)
	}
	nItems := c.len()
	rand.Shuffle(nItems, func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, key := range keys {
		rehash, ok := c.del(key, seed.Hash(key))
		if !ok {
			t.Fatalf("failed to delete key %q", key)
		}
		if rehash {
			c = c.rehash(seed)
		}
	}
	if c.len() != 0 {
		t.Fatalf("failed to erase all items, remains %d", c.len())
	}
}

var sizes2 = []int{1, 200, 400, 600, 800, 1000}

func BenchmarkTable2Hit(b *testing.B) {
	seed := MakeSeed()
	var size int
	for _, v := range sizes2 {
		if v > size {
			size = v
		}
	}
	ss := make([]int, size)
	us := make([]int, size)
	for i := range size {
		ss[i] = i
		us[i] = i
	}

	for _, size := range sizes2 {
		b.Log(size)
		rand.Shuffle(size, func(i, j int) {
			us[i], us[j] = us[j], us[i]
		})

		b.Run(fmt.Sprintf("tbl2 %3d", size), func(b *testing.B) {
			c := newTable(0)
			for i := range size {
				key := ss[i]
				c.add(key, i, seed.Hash(key))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := us[i%size]
				hash := seed.Hash(key)
				_, found := c.get(key, hash)
				if !found {
					b.Fatalf("Key %q should be found", key)
				}
			}
		})

		b.Run(fmt.Sprintf("map  %3d", size), func(b *testing.B) {
			m := make(map[int]int)
			for i := range size {
				m[ss[i]] = i
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := us[i%size]
				_, found := m[key]
				if !found {
					b.Fatalf("Key %q should be found", key)
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	seed := MakeSeed()
	size := maxUsed
	c := newTable(0)
	ss := make([]int, size)
	for i := range size {
		key := i
		c.add(key, i, seed.Hash(key))
		ss[i] = i
	}
	rand.Shuffle(len(ss), func(i, j int) {
		ss[i], ss[j] = ss[j], ss[i]
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := ss[i%size]
		hash := seed.Hash(key)
		_, found := c.get(key, hash)
		if !found {
			b.Fatalf("Key %q should be found", key)
		}
	}
}
