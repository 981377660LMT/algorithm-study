// 二维仿射变换.
// api:
//
//  1. New() -> [3][3]T
//  2. Shift(mapping [3][3]T, shiftX, shiftY T) -> [3][3]T
//  3. Expand(mapping [3][3]T, ratioX, ratioY T) -> [3][3]T
//  4. Rotate90Clockwise(mapping [3][3]T) -> [3][3]T
//  5. Rotate90AntiClockwise(mapping [3][3]T) -> [3][3]T
//  6. RotateClockwise(mapping [3][3]T, degree T) -> [3][3]T
//  7. RotateAntiClockwise(mapping [3][3]T, degree T) -> [3][3]T
//  8. XSymmetricalMove(mapping [3][3]T, x T) -> [3][3]T
//  9. YSymmetricalMove(mapping [3][3]T, y T) -> [3][3]T
//  10. Get(mapping [3][3]T, x, y T) -> (T, T)

package main

import (
	"fmt"
	"math"
)

func main() {
	M := NewAffineMapping[int]()
	m0 := M.New()
	m1 := M.Rotate90AntiClockwise(m0)
	a, b := M.Get(m1, 1, 1)
	fmt.Println(a, b)
	m2 := M.Rotate90Clockwise(m0)
	a, b = M.Get(m2, 1, 1)
	fmt.Println(a, b)

	m3 := M.RotateAntiClockwise(m0, 90)
	a, b = M.Get(m3, 1, 1)
	fmt.Println(a, b)
	m4 := M.RotateClockwise(m0, 90)
	a, b = M.Get(m4, 1, 1)
	fmt.Println(a, b)

	m5 := M.XSymmetricalMove(m0, 2)
	a, b = M.Get(m5, 1, 1)
	fmt.Println(a, b)

	m6 := M.YSymmetricalMove(m0, 2)
	a, b = M.Get(m6, 1, 1)
	fmt.Println(a, b)

	m7 := M.Shift(m0, 2, 2)
	a, b = M.Get(m7, 1, 1)
	fmt.Println(a, b)

	m8 := M.Expand(m0, 2, 2)
	a, b = M.Get(m8, 1, 1)
	fmt.Println(a, b)
}

type Num interface {
	int | int8 | int16 | int32 | int64 | float64
}

type AffineMapping[T Num] struct{}

func NewAffineMapping[T Num]() *AffineMapping[T] { return &AffineMapping[T]{} }

func (am *AffineMapping[T]) New() [3][3]T {
	return [3][3]T{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
}

// 平移变换.
func (am *AffineMapping[T]) Shift(mapping [3][3]T, shiftX, shiftY T) [3][3]T {
	b := [3][3]T{{1, 0, shiftX}, {0, 1, shiftY}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 伸缩变换.
func (am *AffineMapping[T]) Expand(mapping [3][3]T, ratioX, ratioY T) [3][3]T {
	b := [3][3]T{{ratioX, 0, 0}, {0, ratioY, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 🔁 顺时针旋转变换.
func (am *AffineMapping[T]) Rotate90Clockwise(mapping [3][3]T) [3][3]T {
	b := [3][3]T{{0, 1, 0}, {-1, 0, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 🔄 逆时针旋转变换.
func (am *AffineMapping[T]) Rotate90AntiClockwise(mapping [3][3]T) [3][3]T {
	b := [3][3]T{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 顺时针旋转变换.
func (am *AffineMapping[T]) RotateClockwise(mapping [3][3]T, degree T) [3][3]T {
	radian := math.Pi / 180 * float64(degree)
	cos := T(math.Cos(radian))
	sin := T(math.Sin(radian))
	b := [3][3]T{{cos, sin, 0}, {-sin, cos, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 逆时针旋转变换.
func (am *AffineMapping[T]) RotateAntiClockwise(mapping [3][3]T, degree T) [3][3]T {
	radian := math.Pi / 180 * float64(degree)
	cos := T(math.Cos(radian))
	sin := T(math.Sin(radian))
	b := [3][3]T{{cos, -sin, 0}, {sin, cos, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// x轴对称移动(x轴对称).
func (am *AffineMapping[T]) XSymmetricalMove(mapping [3][3]T, x T) [3][3]T {
	b := [3][3]T{{-1, 0, 2 * x}, {0, 1, 0}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// y轴对称移动(y轴对称).
func (am *AffineMapping[T]) YSymmetricalMove(mapping [3][3]T, y T) [3][3]T {
	b := [3][3]T{{1, 0, 0}, {0, -1, 2 * y}, {0, 0, 1}}
	return am._matmul3(b, mapping)
}

// 获取变换后的坐标.
func (am *AffineMapping[T]) Get(mapping [3][3]T, x, y T) (T, T) {
	a0 := mapping[0]
	a1 := mapping[1]
	return a0[0]*x + a0[1]*y + a0[2], a1[0]*x + a1[1]*y + a1[2]
}

func (am *AffineMapping[T]) _matmul3(a, b [3][3]T) [3][3]T {
	res := [3][3]T{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			for j := 0; j < 3; j++ {
				res[i][j] += b[k][j] * a[i][k]
			}
		}
	}
	return res
}
