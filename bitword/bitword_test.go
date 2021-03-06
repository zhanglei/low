package bitword

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToAndFromStr(t *testing.T) {
	cases := []struct {
		src  string
		n    int
		want []byte
	}{
		{"", 1,
			[]byte{}},
		{"", 2,
			[]byte{}},
		{"", 4,
			[]byte{}},
		{"", 8,
			[]byte{}},
		{"a", 1,
			[]byte{0, 1, 1, 0, 0, 0, 0, 1}},
		{"a", 2,
			[]byte{0x1, 0x2, 0x0, 0x1}},
		{"a", 4,
			[]byte{0x6, 0x1}},
		{"a", 8,
			[]byte{0x61}},
		{"\x00", 4,
			[]byte{0, 0}},
		{"\x01\x02\xff", 1,
			[]byte{
				0, 0, 0, 0, 0, 0, 0, 1,
				0, 0, 0, 0, 0, 0, 1, 0,
				1, 1, 1, 1, 1, 1, 1, 1}},
		{"\x01\x02\xff", 2,
			[]byte{0, 0, 0, 1, 0, 0, 0, 2, 3, 3, 3, 3}},
		{"\x01\x02\xff", 4,
			[]byte{0, 1, 0, 2, 0xf, 0xf}},
		{"\x01\x02\xff", 8,
			[]byte{1, 2, 0xff}},
		{"我", 1,
			[]byte{
				1, 1, 1, 0, 0, 1, 1, 0,
				1, 0, 0, 0, 1, 0, 0, 0,
				1, 0, 0, 1, 0, 0, 0, 1,
			},
		},
		{"我", 2,
			[]byte{
				3, 2, 1, 2,
				2, 0, 2, 0,
				2, 1, 0, 1,
			},
		},
		{"我", 4,
			[]byte{0xe, 0x6, 0x8, 0x8, 0x9, 0x1}},
		{"我", 8,
			[]byte{0xe6, 0x88, 0x91},
		},
	}

	for i, c := range cases {
		res := BitWord[c.n].FromStr(c.src)

		if !reflect.DeepEqual(res, c.want) {
			t.Errorf("test %d: got %#v, want %#v",
				i+1, res, c.want)
		}

		str := BitWord[c.n].ToStr(res)
		if str != c.src {
			t.Fatalf(" expect: %v; but: %v", c.src, str)
		}
	}
}

func TestToStrIncomplete(t *testing.T) {
	ta := require.New(t)
	got := BitWord[4].ToStr([]byte{1, 2, 3})
	want := "\x12\x30"
	ta.Equal(want, got)
}

func TestToAndFromStrs(t *testing.T) {

	cases := []struct {
		input []string
		n     int
		want  [][]byte
	}{
		{[]string{"a", "bc", "d"}, 4,
			[][]byte{
				{6, 1},
				{6, 2, 6, 3},
				{6, 4},
			},
		},
		{[]string{"a", "bc", "d"}, 2,
			[][]byte{
				{1, 2, 0, 1},
				{1, 2, 0, 2, 1, 2, 0, 3},
				{1, 2, 1, 0},
			},
		},
	}

	for i, c := range cases {
		rst := BitWord[c.n].FromStrs(c.input)
		if !reflect.DeepEqual(c.want, rst) {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}

		strs := BitWord[c.n].ToStrs(rst)
		if !reflect.DeepEqual(c.input, strs) {
			t.Fatalf("%d-th expect: %v; but: %v", i+1, c.input, strs)
		}
	}
}

func TestGet(t *testing.T) {

	ta := require.New(t)

	type getInput struct {
		s   string
		n   int
		ith int
	}

	cases := []struct {
		input getInput
		want  byte
	}{
		{getInput{"a", 1, 0}, 0},
		{getInput{"a", 1, 1}, 1},
		{getInput{"a", 1, 2}, 1},
		{getInput{"a", 1, 3}, 0},
		{getInput{"a", 1, 4}, 0},
		{getInput{"a", 1, 5}, 0},
		{getInput{"a", 1, 6}, 0},
		{getInput{"a", 1, 7}, 1},

		{getInput{"a", 2, 0}, 1},
		{getInput{"a", 2, 1}, 2},
		{getInput{"a", 2, 2}, 0},
		{getInput{"a", 2, 3}, 1},

		{getInput{"a", 4, 0}, 6},
		{getInput{"a", 4, 1}, 1},

		{getInput{"a", 8, 0}, 0x61},

		{getInput{"abc", 4, 0}, 6},
		{getInput{"abc", 4, 1}, 1},
		{getInput{"abc", 4, 2}, 6},
		{getInput{"abc", 4, 3}, 2},
		{getInput{"abc", 4, 4}, 6},
		{getInput{"abc", 4, 5}, 3},
	}

	for i, c := range cases {
		got := BitWord[c.input.n].Get(c.input.s, c.input.ith)
		ta.Equal(c.want, got,
			"%d-th: input: %#v; want: %#v; got: %#v",
			i+1, c.input, c.want, got)
	}
}

func TestGet_panic(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		s   string
		n   int
		ith int
	}{
		{"", 1, 0},
		{"", 2, 0},
		{"", 3, 0},
		{"", 3, 1},
		{"a", 1, 8},
		{"a", 1, 9},
		{"ab", 1, 16},

		{"a", 2, 4},
		{"a", 4, 2},
		{"a", 8, 1},
		{"a", 8, 2},
	}

	for i, c := range cases {
		ta.Panics(func() { BitWord[c.n].Get(c.s, c.ith) }, "%d-th: %v", i+1, c)
	}
}

func TestFirstDiff(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		a, b      string
		n         int
		from, end int
		want      int
	}{
		{"", "", 1, 0, 0, 0},

		// 0x61 0x62
		{"a", "b", 1, 0, 0, 0},
		{"a", "b", 1, 0, 1, 1},
		{"a", "b", 1, 0, 2, 2},
		{"a", "b", 1, 0, 6, 6},
		{"a", "b", 1, 0, 7, 6},

		{"aa", "ab", 1, 5, 7, 7},
		{"aa", "ab", 1, 5, 14, 14},
		{"aa", "ab", 1, 5, 15, 14},

		{"aa", "ab", 1, 15, 15, 15},
		{"aa", "ab", 2, 2, 7, 7},
		{"aa", "ab", 4, 0, 4, 3},
		{"aa", "ab", 8, 0, 1000, 1},

		{"aa", "aa", 4, 0, 4, 4},

		{"aac", "aa", 4, 0, 4, 4},
		{"aac", "ab", 4, 0, 4, 3},
		{"aac", "ab", 4, 0, 100, 3},

		// bug in 0.1.2
		{"aaa", "aaa", 4, 0, -1, 6},
	}

	for i, c := range cases {
		got := BitWord[c.n].FirstDiff(c.a, c.b, c.from, c.end)
		ta.Equal(c.want, got, "%d-th: case: %+v", i+1, c)
	}
}
