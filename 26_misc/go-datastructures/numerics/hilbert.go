/*
Package Hilbert is designed to allow consumers to find the Hilbert
distance on the Hilbert curve if given a 2 dimensional coordinate.
This could be useful for hashing or constructing a Hilbert R-Tree.
Algorithm taken from here:

http://en.wikipedia.org/wiki/Hilbert_curve

This expects coordinates in the range [0, 0] to [MaxInt32, MaxInt32].
Using negative values for x and y will have undefinied behavior.

Benchmarks:
BenchmarkEncode-8	10000000	       181 ns/op
BenchmarkDecode-8	10000000	       191 ns/op
*/
package main

import "fmt"

// **希尔伯特曲线**是一种空间填充曲线，可以把 2D 坐标与 1D 距离一一对应地映射
func main() {
	// 例：对 (x, y) = (1000, 2000) 进行 Hilbert 编码
	x, y := int32(1000), int32(2000)
	distance := Encode(x, y)
	fmt.Println("Hilbert distance =", distance)

	// 然后 Decode 回来
	x2, y2 := Decode(distance)
	fmt.Printf("Decoded back to: (%d, %d)\n", x2, y2)

	if x2 == x && y2 == y {
		fmt.Println("Test passed: decode(encode(x,y)) = (x,y).")
	} else {
		fmt.Println("Test failed: mismatch.")
	}
}

// n defines the maximum power of 2 that can define a bound,
// this is the value for 2-d space if you want to support
// all hilbert ids with a single integer variable
const n = 1 << 31

// Encode will encode the provided x and y coordinates into a Hilbert
// distance.
func Encode(x, y int32) int64 {
	var rx, ry int32
	var d int64
	for s := int32(n / 2); s > 0; s /= 2 {
		rx = boolToInt(x&s > 0)
		ry = boolToInt(y&s > 0)
		d += int64(int64(s) * int64(s) * int64(((3 * rx) ^ ry)))
		rotate(s, rx, ry, &x, &y)
	}

	return d
}

// Decode will decode the provided Hilbert distance into a corresponding
// x and y value, respectively.
func Decode(h int64) (int32, int32) {
	var ry, rx int64
	var x, y int32
	t := h

	for s := int64(1); s < int64(n); s *= 2 {
		rx = 1 & (t / 2)
		ry = 1 & (t ^ rx)
		rotate(int32(s), int32(rx), int32(ry), &x, &y)
		x += int32(s * rx)
		y += int32(s * ry)
		t /= 4
	}

	return x, y
}

func boolToInt(value bool) int32 {
	if value {
		return int32(1)
	}

	return int32(0)
}

func rotate(n, rx, ry int32, x, y *int32) {
	if ry == 0 {
		if rx == 1 {
			*x = n - 1 - *x
			*y = n - 1 - *y
		}

		t := *x
		*x = *y
		*y = t
	}
}
