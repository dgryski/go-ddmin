package ddmin

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"testing/quick"
)

func ExampleMinimize() {

	// the example from the paper
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	f := func(d []byte) Result {
		var seen1, seen7, seen8 bool
		for _, v := range d {
			if v == 1 {
				seen1 = true
			}
			if v == 7 {
				seen7 = true
			}
			if v == 8 {
				seen8 = true
			}

		}

		if seen1 && seen7 && seen8 {
			return Fail
		} else {
			return Pass
		}
	}

	m := Minimize(data, f)

	fmt.Println(m)

	// Output: [1 7 8]
}

func TestSplit(t *testing.T) {

	f := func(i byte) bool {
		if i < 2 {
			return true
		}

		n := rand.Intn(int(i))
		if n < 2 {
			return true
		}

		i16 := uint16(int(i)*rand.Intn(10) + rand.Intn(int(i)))
		b := make([]byte, i16)
		for i := range b {
			b[i] = byte(rand.Intn(256))
		}
		s := makeSubsets(b, n)
		// make sure we got the correct number of subsets
		if len(s) != n {
			return false
		}

		// make sure we have only one or two different slice lengths
		m := make(map[int]int)
		// make sure we got back the entire slice we sent in
		got := make([]byte, 0, len(b))
		for _, v := range s {
			m[len(v)]++
			got = append(got, v...)
		}
		if len(m) != 1 && len(m) != 2 {
			return false
		}
		if !bytes.Equal(got, b) {
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
