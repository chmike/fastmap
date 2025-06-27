package altmapint

import "testing"

func TestSetEmpty(t *testing.T) {
	tests := []struct {
		set Set
		out bool
	}{
		{set: 0x0000_0000_0000_0000, out: true},
		{set: 0x8080_8000_8080_8080, out: false},
		{set: 0x0080_8000_8080_8080, out: false},
	}
	for i, test := range tests {
		out := test.set.Empty()
		if out != test.out {
			t.Errorf("%d for bit set %016x expect %v, got %v", i, test.set, test.out, out)
		}
	}
}

func TestSetLen(t *testing.T) {
	tests := []struct {
		set Set
		out int
	}{
		{set: 0x0000_0000_0000_0000, out: 0},
		{set: 0x8080_8000_8080_8080, out: 7},
		{set: 0x0080_8000_8080_8080, out: 6},
	}
	for i, test := range tests {
		out := test.set.Len()
		if out != test.out {
			t.Errorf("%d for bit set %016x expect %v, got %v", i, test.set, test.out, out)
		}
	}
}

func TestSetNext(t *testing.T) {
	tests := []struct {
		set Set
		out Set
	}{
		{set: 0x0000_0000_0000_0000, out: 0x0000_0000_0000_0000},
		{set: 0x8080_8000_8080_8080, out: 0x8080_8000_8080_8000},
		{set: 0x0080_8000_8080_0000, out: 0x0080_8000_8000_0000},
	}
	for i, test := range tests {
		out := test.set.Next()
		if out != test.out {
			t.Errorf("%d for bit set %016x expect %v, got %v", i, test.set, test.out, out)
		}
	}
}

func TestSetPos(t *testing.T) {
	tests := []struct {
		set Set
		out int
	}{
		{set: 0x0000_0000_0000_0000, out: 8},
		{set: 0x8080_8000_8080_8080, out: 0},
		{set: 0x0080_8000_8080_0000, out: 2},
	}
	for i, test := range tests {
		out := test.set.Pos()
		if out != test.out {
			t.Errorf("%d for bit set %016x expect %v, got %v", i, test.set, test.out, out)
		}
	}
}
