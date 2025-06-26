package altmap

import (
	"unsafe"
)

// Cache is a map using an extensible directory to tables of tableSize groups.
// The table uses 8bit top hashes with tombstones and doesn't move items.
// A table is split when it contains more than maxItems.
type Cache struct {
	tables  []*table // directory of tables
	seed    Seed     // hash seed
	nItems  int      // number of stored items
	depth   byte     // depth of the directory
	mask    uint     // mask for hash
	basePtr **table  // pointer on first entry in tables
}

func (c *Cache) Init() {
	c.tables = []*table{newTable(0)}
	c.seed = MakeSeed()
	c.nItems = 0
	c.depth = 0
	c.mask = 0
	c.basePtr = unsafe.SliceData(c.tables)
}

// Len returns the number of items stored in the cache.
func (c *Cache) Len() int {
	return c.nItems
}

// Cap returns the number of item slots in the cache.
func (c *Cache) Cap() int {
	return len(c.tables) * tableItems
}

// H0 returns the hash used for the directory.
func H0(hash uint) uint {
	return hash >> tableHashBits
}

// table return pointer on the table corresponding to the given hash value.
func (c *Cache) table(hash uint) *table {
	offset := (hash >> (tableHashBits - 3)) & c.mask
	return *(**table)(unsafe.Add(unsafe.Pointer(c.basePtr), offset))
}

// Get returns the value associated to key and true if it is found.
func (c *Cache) Get(key string) (value int, ok bool) {
	hash := c.seed.Hash(key)
	t := c.table(hash)
	pattern := MakePattern(H2(hash))
	var pos uint32
	idx := H1(hash) & (tableSize - 1)
	offset := uint32(idx) * sizeGroup
	basePtr := unsafe.Pointer(unsafe.SliceData(t.groups[:]))
	for {
		g := (*Group)(unsafe.Add(basePtr, offset))
		for set := g.header.Find(pattern); !set.Empty(); set = set.Next() {
			if item := &g.item[set.Pos()&(nItems-1)]; item.key == key {
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

// Add swaps the value and return true if the key is found in the cache,
// otherwise it adds the key and value and returns false.
func (c *Cache) Add(key string, value int) (oldValue int, ok bool) {
	hash := c.seed.Hash(key)
	t := c.table(hash)
	if oldValue, ok = t.swap(key, value, hash); ok {
		return
	}

	for !t.add(&key, value, hash) {

		// the table is full, it must be split
		l := uint(len(c.tables))
		if t.depth == c.depth {
			// grow the directory
			tables := c.tables
			l2 := l * 2
			c.tables = make([]*table, l2)
			copy(c.tables, tables)
			copy(c.tables[l:], tables)
			c.depth++
			c.mask = (l2 - 1) * 8 // pre multiply mask by pointer byte size
			c.basePtr = unsafe.SliceData(c.tables)
			l = l2
		}

		step := uint(1 << t.depth)    // interval between pointers to the table
		tIdx := H0(hash) & (step - 1) // index to the first table pointer in the table
		t1, t2 := t.split(step, c.seed)

		for tIdx < l {
			c.tables[tIdx] = t1
			tIdx += step
			c.tables[tIdx] = t2
			tIdx += step
		}

		t = c.table(hash)
	}
	c.nItems++
	return
}

// Del deletes key from the cache.
func (c *Cache) Del(key string) {
	hash := c.seed.Hash(key)
	t := c.table(hash)
	rehash, ok := t.del(&key, hash)
	if ok {
		c.nItems--
		if rehash {
			t2 := t.rehash(c.seed)
			step := uint(1 << t.depth) // interval between pointers to the table
			for tIdx, l := H0(hash)&(step-1), uint(len(c.tables)); tIdx < l; tIdx += step {
				c.tables[tIdx] = t2
			}
		}
	}
}
