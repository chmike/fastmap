package altmap

import (
	"fmt"
	"math/bits"
)

// A Set uses one byte per slot with the value 0x80 when the slot is member of the set
// and 0x00 whe it is not. A Set is invalid if it contains a different byte value.
type Set uint

// Empty returns true if the set is empty.
func (b Set) Empty() bool {
	return b == 0
}

// Len returns the number of members in the set.
func (b Set) Len() int {
	return bits.OnesCount(uint(b))
}

// Next returns a bitSet where the less significant (first) member of the set
// is removed. Requires the Set is not empty.
func (b Set) Next() Set {
	return b & (b - 1)
}

// Pos returns the index of the less significant (first) member of the set.
// Returns uintByteLength if the set is empty.
func (b Set) Pos() int {
	return bits.TrailingZeros(uint(b)) >> 3
}

// Check returns an error if the set is invalid.
func (b Set) Check() error {
	if b&Set(0x7f7f_7f7f_7f7F_7f7f) != 0 {
		return fmt.Errorf("invalid set %016x", b)
	}
	return nil
}

// PSet is a packed set where each bit represent the status of a member.
// The bit is set if the slot is member of the set.
// N.B. it is provided as documentation only as it isn't used.
type PSet byte

// Pack returns the packed set of the set b.
func (b Set) Pack() PSet {
	normalized := b >> 7
	gathered := normalized * 0x0102040810204080
	return PSet(gathered >> 56)
}

// Empty returns true if the packed set is empty.
func (b PSet) Empty() bool {
	return b == 0
}

// Pos returns the index of the first member of the set.
func (b PSet) Pos() int {
	return bits.TrailingZeros8(uint8(b))
}

// Next returns a PSet with its first member removed. Requires
// the PSet is not empty.
func (b PSet) Next() PSet {
	return b & (b - 1)
}
