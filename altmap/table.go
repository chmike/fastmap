package altmap

import (
	"iter"
	"unsafe"
)

// nItems is the number of items in a group.
const nItems = uintByteLength

// tableSizeLog2 is the log base 2 of the number of groups in a table.
const tableSizeLog2 = 8
const tableSize = 1 << tableSizeLog2

// sizeGroup is the byte size of a group.
const sizeGroup = uint32(unsafe.Sizeof(Group{}))

// sizeGroups is the byte size of all groups.
const sizeGroups = sizeGroup * tableSize

const tableItems = tableSize * nItems

// number of hash bits used by a table
const tableHashBits = topHashBits + tableSizeLog2

// maxUsed is the minimum number of free slots triggering a table split.
const maxUsed = (tableItems * 90) / 100

// maxTombstones is the maximum number of tombstones a table should contain.
const maxTombstones = (tableItems * 15) / 100

type Item struct {
	key   string
	value int
}

type Group struct {
	header Hdr
	item   [nItems]Item
}

type table struct {
	groups      [tableSize]Group // array of groups
	nItems      uint16           // number of items (used only to measure table occupancy)
	nTombstones uint16           // number of tombstones
	depth       byte             // depth of table in the directory
}

// newTable returns a new table of the given depth.
func newTable(depth byte) *table {
	return &table{depth: depth}
}

// len returns the number of items stored in the table.
func (t *table) len() int {
	return int(t.nItems)
}

// cap returns the maximum capacity in items of the table.
func (t *table) cap() int {
	return tableItems
}

// occupancy returns the occupancy of the table.
func (t *table) occupancy() int {
	return (t.len() * 100) / t.cap()
}

func makeOffset(h1 uint) uint32 {
	return (uint32(h1) & (tableSize - 1)) * sizeGroup
}

// get returns the value associated to key if found in the table. hash is the hash value of key.
// Returns false and the default value if not found.
func (t *table) get(key *string, hash uint) (value int, ok bool) {
	pattern := MakePattern(H2(hash))
	var pos uint32
	offset := makeOffset(H1(hash))
	basePtr := unsafe.Pointer(unsafe.SliceData(t.groups[:]))
	for {
		g := (*Group)(unsafe.Add(basePtr, offset))
		for set := g.header.Find(pattern); !set.Empty(); set = set.Next() {
			if item := &g.item[set.Pos()&(nItems-1)]; item.key == *key {
				return item.value, true
			}
		}
		if g.header.HasFreeSlots() {
			return
		}
		// to avoid a product by groupSize or a modulo
		// pos never reach sizeGroups
		pos += sizeGroup
		if offset += pos; offset >= sizeGroups {
			offset -= sizeGroups
		}
	}
}

// swap swaps the value associated with the key if found in the table. hash is the hash of key.
// Returns the default value and false of the key is not found in the table.
func (t *table) swap(key string, value int, hash uint) (oldValue int, ok bool) {
	pattern := MakePattern(H2(hash))
	var pos uint32
	offset := makeOffset(H1(hash))
	basePtr := unsafe.Pointer(unsafe.SliceData(t.groups[:]))
	for {
		g := (*Group)(unsafe.Add(basePtr, offset))
		for set := g.header.Find(pattern); !set.Empty(); set = set.Next() {
			if item := &g.item[set.Pos()]; item.key == key {
				oldValue, item.value = item.value, value
				return oldValue, true
			}
		}
		if g.header.HasFreeSlots() {
			return
		}
		// to avoid a product by groupSize or a modulo
		// pos never reach sizeGroups
		pos += sizeGroup
		if offset += pos; offset >= sizeGroups {
			offset -= sizeGroups
		}
	}
}

// add adds the key and value to the table. Requires that the key is not in the
// table. Returns true if succeeded, and false if the table is full.
func (t *table) add(key *string, value int, hash uint) bool {
	if int(t.nItems)+int(t.nTombstones) > maxUsed {
		return false
	}
	var pos uint32
	offset := makeOffset(H1(hash))
	basePtr := unsafe.Pointer(unsafe.SliceData(t.groups[:]))
	for {
		g := (*Group)(unsafe.Add(basePtr, offset))
		if set := g.header.FindUnused(); !set.Empty() {
			// pick first unused slot in header
			i := set.Pos()
			g.header = g.header.Set(i, H2(hash))
			g.item[i] = Item{key: *key, value: value}
			t.nItems++
			return true
		}
		// to avoid a product by groupSize or a modulo
		// pos never reach sizeGroups
		pos += sizeGroup
		if offset += pos; offset >= sizeGroups {
			offset -= sizeGroups
		}
	}
}

func (t *table) items() iter.Seq2[string, int] {
	return func(yield func(string, int) bool) {
		for i := range tableSize {
			g := &t.groups[i]
			h := g.header
			for j := range nItems {
				if byte(h)&0x7F != 0 {
					if !yield(g.item[j].key, g.item[j].value) {
						return
					}
				}
				h >>= 8
			}
		}
	}
}

func (t *table) split(bit uint, seed Seed) (t1, t2 *table) {
	bit <<= tableHashBits
	t1, t2 = newTable(t.depth+1), newTable(t.depth+1)
	for k, v := range t.items() {
		hash := seed.Hash(k)
		if hash&bit == 0 {
			if !t1.add(&k, v, hash) {
				panic("failed to split")
			}
		} else {
			if !t2.add(&k, v, hash) {
				panic("failed to split")
			}
		}
	}
	return t1, t2
}

// rehash rehashes table to remove all tombstones.
func (t *table) rehash(seed Seed) *table {
	t2 := newTable(t.depth)
	for k, v := range t.items() {
		if !t2.add(&k, v, seed.Hash(k)) {
			panic("failed rehashing")
		}
	}
	return t2
}

// del deletes the item with the given key. Returns true if the number of tombstones
// exceeds a threshold.
func (t *table) del(key *string, hash uint) (rehash bool, ok bool) {
	pattern := MakePattern(H2(hash))
	var pos uint32
	offset := makeOffset(H1(hash))
	basePtr := unsafe.Pointer(unsafe.SliceData(t.groups[:]))
	for {
		g := (*Group)(unsafe.Add(basePtr, offset))
		for set := g.header.Find(pattern); !set.Empty(); set = set.Next() {
			i := set.Pos()
			if item := &g.item[i]; item.key == *key {
				*item = Item{}
				g.header = g.header.Set(i, tombstone)
				t.nTombstones++
				t.nItems--
				return int(t.nTombstones) > maxTombstones, true
			}
		}
		if g.header.HasFreeSlots() {
			return false, false
		}
		// to avoid a product by groupSize or a modulo
		// pos never reach sizeGroups
		pos += sizeGroup
		if offset += pos; offset >= sizeGroups {
			offset -= sizeGroups
		}
	}
}
