// 用于判断两个浮点数是否相等

package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(AlmostEqual(1.0, 1.0, 0.01))  // true
	fmt.Println(AlmostEqual(1.0, 1.01, 0.01)) // true
}

var (
	// MinNormal is the smallest positive normal value of type float64.
	MinNormal = math.Float64frombits(0x0010000000000000)
	// MinNormal32 is the smallest positive normal value of type float32.
	MinNormal32 = math.Float32frombits(0x00800000)
)

// AlmostEqual returns true if a and b are equal within a relative error of
// ε. See http://floating-point-gui.de/errors/comparison/ for the details of the
// applied method.
func AlmostEqual(a, b, ε float64) bool {
	if a == b {
		return true
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || absA+absB < MinNormal {
		return diff < ε*MinNormal
	}
	return diff/math.Min(absA+absB, math.MaxFloat64) < ε
}

// AlmostEqual32 returns true if a and b are equal within a relative error of
// ε. See http://floating-point-gui.de/errors/comparison/ for the details of the
// applied method.
func AlmostEqual32(a, b, ε float32) bool {
	if a == b {
		return true
	}
	absA := Abs32(a)
	absB := Abs32(b)
	diff := Abs32(a - b)
	if a == 0 || b == 0 || absA+absB < MinNormal32 {
		return diff < ε*MinNormal32
	}
	return diff/Min32(absA+absB, math.MaxFloat32) < ε
}

// Abs32 works like math.Abs, but for float32.
func Abs32(x float32) float32 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}

// Min32 works like math.Min, but for float32.
func Min32(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}
