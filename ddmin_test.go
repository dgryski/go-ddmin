package ddmin

import (
	"bytes"
	"testing"
)

func TestMinimize(t *testing.T) {

	// the example from the paper
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}

	f := func(d []byte) bool {
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
		return !(seen1 && seen7 && seen8)
	}

	want := []byte{1, 7, 8}

	m := Minimize(data, f)

	if !bytes.Equal(m, want) {
		t.Errorf("Minimize()=% 02x, want % 02x\n", m, want)
	}
}
