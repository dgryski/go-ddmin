package ddmin

import (
	"fmt"
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
