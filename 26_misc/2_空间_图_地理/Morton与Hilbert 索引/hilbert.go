// https://github.com/google/hilbert
// 用于将值映射到空间填充曲线（如希尔伯特曲线和皮亚诺曲线）及其反向映射的 Go 包。
// google s2 希尔伯特曲线 https://www.shenyanchao.cn/blog/2020/01/16/geo_google_s2/

package main

import (
	"errors"
	"fmt"
)

func main() {
	// Create a new Hilbert space
	h, err := NewHilbert(8)
	if err != nil {
		panic(err)
	}

	// Map a value to the curve
	x, y, err := h.Map(42)
	if err != nil {
		panic(err)
	}

	// Map the coordinates back to the value
	t, err := h.MapInverse(x, y)
	if err != nil {
		panic(err)
	}

	// Output the result
	fmt.Println(x, y, t)
}

// #region hilbert

// Hilbert represents a 2D Hilbert space of order N for mapping to and from.
// Implements SpaceFilling interface.
type Hilbert struct {
	N int
}

// NewHilbert returns a Hilbert space which maps integers to and from the curve.
// n must be a power of two.
func NewHilbert(n int) (*Hilbert, error) {
	if n <= 0 {
		return nil, ErrNotPositive
	}

	// Test if power of two
	if (n & (n - 1)) != 0 {
		return nil, ErrNotPowerOfTwo
	}

	return &Hilbert{
		N: n,
	}, nil
}

// GetDimensions returns the width and height of the 2D space.
func (s *Hilbert) GetDimensions() (int, int) {
	return s.N, s.N
}

// Map transforms a one dimension value, t, in the range [0, n^2-1] to coordinates on the Hilbert
// curve in the two-dimension space, where x and y are within [0,n-1].
func (s *Hilbert) Map(t int) (x, y int, err error) {
	if t < 0 || t >= s.N*s.N {
		return -1, -1, ErrOutOfRange
	}

	for i := 1; i < s.N; i = i * 2 {
		rx := t&2 == 2
		ry := t&1 == 1
		if rx {
			ry = !ry
		}

		x, y = s.rotate(i, x, y, rx, ry)

		if rx {
			x = x + i
		}
		if ry {
			y = y + i
		}

		t /= 4
	}

	return
}

// MapInverse transform coordinates on Hilbert curve from (x,y) to t.
func (s *Hilbert) MapInverse(x, y int) (t int, err error) {
	if x < 0 || x >= s.N || y < 0 || y >= s.N {
		return -1, ErrOutOfRange
	}

	for i := s.N / 2; i > 0; i = i / 2 {
		rx := (x & i) > 0
		ry := (y & i) > 0

		a := 0
		if rx {
			a = 3
		}
		t += i * i * (a ^ b2i(ry))

		x, y = s.rotate(i, x, y, rx, ry)
	}

	return
}

// rotate rotates and flips the quadrant appropriately.
func (s *Hilbert) rotate(n, x, y int, rx, ry bool) (int, int) {
	if !ry {
		if rx {
			x = n - 1 - x
			y = n - 1 - y
		}

		x, y = y, x
	}
	return x, y
}

// #endregion

// Errors returned when validating input.
var (
	ErrNotPositive     = errors.New("N must be greater than zero")
	ErrNotPowerOfTwo   = errors.New("N must be a power of two")
	ErrNotPowerOfThree = errors.New("N must be a power of three")
	ErrOutOfRange      = errors.New("value is out of range")
)

// SpaceFilling represents a space-filling curve that can map points from one dimensions to two.
type SpaceFilling interface {
	// Map transforms a one dimension value, t, in the range [0, n^2-1] to coordinates on the
	// curve in the two-dimension space, where x and y are within [0,n-1].
	Map(t int) (x, y int, err error)

	// MapInverse transform coordinates on the curve from (x,y) to t.
	MapInverse(x, y int) (t int, err error)

	// GetDimensions returns the width and height of the 2D space.
	GetDimensions() (x, y int)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
