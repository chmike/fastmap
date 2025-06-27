package altmapint

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/bits"
	"unsafe"
)

/*
This header uses 8 bit top hashes and tombstones, and doesn't move
items.

The byte 0x00 signals a free slots and 0x80 a tombstone as it avoids
one op when locating unused slots.

When deleting, the top hash is replaced with a tombstone and the item
of the group is cleared to avoid memory leaks.

When inserting, the first unused slot (free slot or tombstone) is
picked as insertion location. The consequence is that free slots are
packed at the end of the group. The presence of at least one free slot
is a criteria to stop the probings.
*/

// uintByteLength is the byte length of an uint. It is 8 on 64 bit
// cpu and 4 on 32bit cpu.
const uintByteLength = int(unsafe.Sizeof(uint(0)))
const topHashBits = 8

const (
	freeSlot  byte = 0x00
	tombstone byte = 0x80
)

type Seed uint64

func MakeSeed() Seed {
	var buf [8]byte
	rand.Read(buf[:])
	return Seed(binary.LittleEndian.Uint64(buf[:]))
}

func (s Seed) Hash(key int) uint {
	return uint(HashUint64(uint64(key), uint64(s)))
}

func H1(hash uint) uint {
	return hash >> topHashBits
}

func H2(hash uint) byte {
	h2 := byte(hash)
	if h2&0x7f == 0 {
		h2 |= 0x01
	}
	return h2
}

type Pattern uint

// MakePattern returns a Pattern to be used with FindByte.
func MakePattern(b byte) Pattern {
	return Pattern(b) * 0x0101_0101_0101_0101
}

// Hdr is a uint packing 8 bits top hashes (byte). Bytes are numbered
// from the less significant to the most significant as 0 to uintByteLength.
type Hdr uint

// HasFreeSlots returns true if h has at least one free slot.
func (h Hdr) HasFreeSlots() bool {
	return uint(h)&^uint(^uint(0)>>8) == 0
}

// Find returns the Set of bytes in header matching the pattern
// created with MakePattern. There are no false positive.
func (h Hdr) Find(pattern Pattern) Set {
	return (h ^ Hdr(pattern)).findZeros()
}

// findZeros returns the set of zeros.
func (h Hdr) findZeros() Set {
	v := uint(h)
	v &= 0x7f7f_7f7f_7f7f_7f7f
	v += 0x7f7f_7f7f_7f7f_7f7f
	v |= uint(h) | 0x7f7f_7f7f_7f7f_7f7f
	return Set(^v)
}

// FindUnused returns the set of unused slots.
func (h Hdr) FindUnused() Set {
	v := uint(h)&0x7f7f_7f7f_7f7f_7f7f + 0x7f7f_7f7f_7f7f_7f7f
	v |= 0x7f7f_7f7f_7f7f_7f7f
	return Set(^v)
}

// FindUsed returns the set of used slots.
func (h Hdr) FindUsed() Set {
	v := uint(h)&0x7f7f_7f7f_7f7f_7f7f + 0x7f7f_7f7f_7f7f_7f7f
	return Set(v & 0x8080_8080_8080_8080)
}

// FirstFree returns the index of the first free slot.
func (h Hdr) FirstFree() int {
	return uintByteLength - bits.LeadingZeros(uint(h))>>3
}

// Set sets byte i in h to b. Requires i is smaller than byteLength.
func (h Hdr) Set(i int, b byte) Hdr {
	i = (i * 8) & 63 // the &63 is to avoid an overflow test for panic
	b ^= byte(h >> i)
	h ^= Hdr(b) << i
	return h
}

// Check returns an error if h is invalid.
func (h Hdr) Check() error {
	// make sure that all free slots are at the end.
	zeros := h.findZeros()
	check := Set(0x0101_0101_0101_0101) << bits.TrailingZeros(uint(zeros))
	if check != zeros {
		return fmt.Errorf("header %016x has non terminal free slots", h)
	}
	return nil
}
