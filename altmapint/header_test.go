package altmapint

import (
	"testing"
)

func TestHdrHasFreeSlots(t *testing.T) {
	tests := []struct {
		hdr Hdr
		out bool
	}{
		// 0
		{hdr: 0x0000_0000_0000_0000, out: true},
		{hdr: 0x0000_002a_6702_0267, out: true},
		{hdr: 0x5d02_6767_0267_8005, out: false},
		{hdr: 0x005d_6767_0267_097f, out: true},
		{hdr: 0x805d_6767_0267_097f, out: false},
		// 5
		{hdr: 0x135d_8067_8267_817f, out: false},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		out := test.hdr.HasFreeSlots()
		if out != test.out {
			t.Errorf("%d for Hdr %016x expect %v, got %v", i, test.hdr, test.out, out)
		}
	}
}

func TestHdrFind(t *testing.T) {
	tests := []struct {
		hdr  Hdr
		set  Set
		hash byte
	}{
		// 0
		{hdr: 0x0002_157f_2a15_095d, hash: 0x02, set: 0x0080_0000_0000_0000},
		{hdr: 0x0017_157f_2a15_235d, hash: 0x15, set: 0x0000_8000_0080_0000},
		{hdr: 0x5d05_1767_0267_0502, hash: 0x05, set: 0x0080_0000_0000_8000},
		{hdr: 0x0000_0000_0000_0000, hash: 0x67, set: 0x0000_0000_0000_0000},
		{hdr: 0x0000_002a_6702_0267, hash: 0x67, set: 0x0000_0000_8000_0080},
		// 5
		{hdr: 0x5151_6767_0267_0502, hash: 0x67, set: 0x0000_8080_0080_0000},
		{hdr: 0x8051_e7e7_0267_0503, hash: 0xe7, set: 0x0000_8080_0000_0000},
		{hdr: 0x0080_8080_0267_0502, hash: 0x02, set: 0x0000_0000_8000_0080},
		{hdr: 0x5151_8080_0267_0502, hash: 0x02, set: 0x0000_0000_8000_0080},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		if err := test.set.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		set := test.hdr.Find(MakePattern(test.hash))
		if set != test.set {
			t.Errorf("%d for Hdr %016x expect set %016x, got %016x", i, test.hdr, test.set, set)
		}
	}
}

func TestHdrFindUnused(t *testing.T) {
	tests := []struct {
		hdr Hdr
		set Set
	}{
		// 0
		{hdr: 0x0002_157f_2a15_095d, set: 0x8000_0000_0000_0000},
		{hdr: 0x8017_157f_2a15_235d, set: 0x8000_0000_0000_0000},
		{hdr: 0x5d05_1767_0267_0502, set: 0x0000_0000_0000_0000},
		{hdr: 0x0000_0000_0000_0000, set: 0x8080_8080_8080_8080},
		{hdr: 0x0000_002a_6703_0267, set: 0x8080_8000_0000_0000},
		// 5
		{hdr: 0x0051_6767_0367_0502, set: 0x8000_0000_0000_0000},
		{hdr: 0x0051_e7e7_8067_0503, set: 0x8000_0000_8000_0000},
		{hdr: 0x0101_8080_0267_0503, set: 0x0000_8080_0000_0000},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		if err := test.set.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		set := test.hdr.FindUnused()
		if set != test.set {
			t.Errorf("%d for Hdr %016x expect set %016x, got %016x", i, test.hdr, test.set, set)
		}
	}
}

func TestHdrFirstFree(t *testing.T) {
	tests := []struct {
		hdr Hdr
		pos int
	}{
		// 0
		{hdr: 0x0002_157f_2a15_095d, pos: 7},
		{hdr: 0x8017_157f_2a15_235d, pos: 8},
		{hdr: 0x5d05_1767_0267_0502, pos: 8},
		{hdr: 0x0000_0000_0000_0000, pos: 0},
		{hdr: 0x0000_002a_6703_0267, pos: 5},
		// 5
		{hdr: 0x0051_6767_0367_0502, pos: 7},
		{hdr: 0x0051_e7e7_8067_0503, pos: 7},
		{hdr: 0x0101_8080_0267_0503, pos: 8},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		pos := test.hdr.FirstFree()
		if pos != test.pos {
			t.Errorf("%d for Hdr %016x expect pos %v, got %v", i, test.hdr, test.pos, pos)
		}
	}
}

func TestHdrFindUsed(t *testing.T) {
	tests := []struct {
		hdr Hdr
		set Set
	}{
		// 0
		{hdr: 0x0002_157f_2a15_095d, set: 0x0080_8080_8080_8080},
		{hdr: 0x8017_157f_2a15_235d, set: 0x0080_8080_8080_8080},
		{hdr: 0x5d05_1767_0267_0502, set: 0x8080_8080_8080_8080},
		{hdr: 0x0000_0000_0000_0000, set: 0x0000_0000_0000_0000},
		{hdr: 0x0000_002a_6703_0267, set: 0x0000_0080_8080_8080},
		// 5
		{hdr: 0x0051_6767_0367_0502, set: 0x0080_8080_8080_8080},
		{hdr: 0x0051_e7e7_8067_0503, set: 0x0080_8080_0080_8080},
		{hdr: 0x0101_8080_0267_0503, set: 0x8080_0000_8080_8080},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		if err := test.set.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		set := test.hdr.FindUsed()
		if set != test.set {
			t.Errorf("%d for Hdr %016x expect set %016x, got %016x", i, test.hdr, test.set, set)
		}
	}
}

func TestHdrSet(t *testing.T) {
	tests := []struct {
		hdr, out Hdr
		i        int
		b        byte
	}{
		{hdr: 0x0000_0000_0000_0000, i: 0, b: 0x55, out: 0x0000_0000_0000_0055},
		{hdr: 0x0000_002a_6705_0267, i: 5, b: 0x2a, out: 0x0000_2a2a_6705_0267},
		{hdr: 0x0000_6767_1867_027f, i: 3, b: 0x02, out: 0x0000_6767_0267_027f},
	}
	for i, test := range tests {
		if err := test.hdr.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
		if err := test.out.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}

		out := test.hdr.Set(test.i, test.b)
		if out != test.out {
			t.Errorf("%d for Hdr %016x, i %d and byte %02x expect Hdr %016x, got %016x", i, test.hdr, test.i, test.b, test.out, out)
		}
		if err := out.Check(); err != nil {
			t.Errorf("%d check: %v", i, err)
		}
	}
}
