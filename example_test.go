package fastnum_test

import (
	"fmt"

	"github.com/saromanov/fastnum"
)

func ExampleAppendInt64Ones() {
	dst := fastnum.AppendInt64Ones(nil, 4)
	fmt.Println(dst)
	// Output: [1 1 1 1]
}

func ExampleAppendFloat64Zeros() {
	dst := []float64{1.5}
	dst = fastnum.AppendFloat64Zeros(dst, 3)
	fmt.Println(dst)
	// Output: [1.5 0 0 0]
}

func ExampleIsInt64Zeros() {
	fmt.Println(fastnum.IsInt64Zeros([]int64{0, 0, 0}))
	fmt.Println(fastnum.IsInt64Zeros([]int64{0, 1, 0}))
	// Output:
	// true
	// false
}

func ExampleIsFloat64Ones() {
	fmt.Println(fastnum.IsFloat64Ones([]float64{1, 1, 1}))
	fmt.Println(fastnum.IsFloat64Ones([]float64{1, -1}))
	// Output:
	// true
	// false
}
