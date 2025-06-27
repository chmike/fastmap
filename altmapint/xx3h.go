package altmapint

import "math/bits"

// HashUint64 returns the xhh3 hash of v using the given seed. Use math.Float64bits
// to convert a float64 to uint64.
func HashUint64(v, seed uint64) uint64 {
	const k0 uint64 = 0x1cad21f72c81017c ^ 0xdb979083e96dd4de // secret[8:]^secret[16:]
	const k1 uint64 = 0x9FB21C651E98DF25                      // prime_MX2

	// standard xxh3 modifies the seed, but we don't as we assume the seed is
	// really random and generated with cipher rand.
	//modifiedSeed := seed ^ (uint64(bits.ReverseBytes32(uint32(seed))) << 32)
	modifiedSeed := seed
	h := (k0 - modifiedSeed) ^ bits.RotateLeft64(v, 32)
	h ^= bits.RotateLeft64(h, 49) ^ bits.RotateLeft64(h, 24)
	h *= k1
	h ^= (h >> 35) + 8
	h *= k1
	h ^= h >> 28
	return h
}

// HashUint32 returns the xhh3 hash of v using the given seed. Use math.Float32bits
// to convert a float32 to uint32.
func HashUint32(v uint32, seed uint64) uint64 {
	const k0 uint64 = 0x1cad21f72c81017c ^ 0xdb979083e96dd4de
	const k1 uint64 = 0x9FB21C651E98DF25

	// standard xxh3 modifies the seed, but we don't as we assume the seed is
	// really random and generated with cipher rand.
	//modifiedSeed := seed ^ (uint64(bits.ReverseBytes32(uint32(seed))) << 32)
	modifiedSeed := seed
	h := (k0 - modifiedSeed) ^ (uint64(v) | uint64(v)<<32)
	h ^= bits.RotateLeft64(h, 49) ^ bits.RotateLeft64(h, 24)
	h *= k1
	h ^= (h >> 35) + 4
	h *= k1
	h ^= h >> 28
	return h
}
