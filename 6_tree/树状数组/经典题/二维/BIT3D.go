// 三维树状数组
// 单点修改，区间查询

package main

import "fmt"

func main() {
	bit3d := NewBIT3D(10, 10, 10)
	bit3d.Add(1, 1, 1, 1)
	fmt.Println(bit3d.QueryRange(0, 0, 0, 10, 10, 10))
}

type BIT3D struct {
	data    [][][]int
	x, y, z int
}

func NewBIT3D(x, y, z int) *BIT3D {
	data := make([][][]int, x+1)
	for i := 0; i <= x; i++ {
		data[i] = make([][]int, y+1)
		for j := 0; j <= y; j++ {
			data[i][j] = make([]int, z+1)
		}
	}
	return &BIT3D{data: data, x: x + 1, y: y + 1, z: z + 1}
}

// 0<=x<X, 0<=y<Y, 0<=z<Z
func (b *BIT3D) Add(x, y, z, v int) {
	for i := x; i < b.x; i |= i + 1 {
		for j := y; j < b.y; j |= j + 1 {
			for k := z; k < b.z; k |= k + 1 {
				b.data[i][j][k] += v
			}
		}
	}
}

// 0<=x<X, 0<=y<Y, 0<=z<Z
func (b *BIT3D) Query(x, y, z int) int {
	res := 0
	x--
	y--
	z--
	for i := x; i >= 0; i = (i & (i + 1)) - 1 {
		for j := y; j >= 0; j = (j & (j + 1)) - 1 {
			for k := z; k >= 0; k = (k & (k + 1)) - 1 {
				res += b.data[i][j][k]
			}
		}
	}
	return res
}

// 0<=x1<=x2<X, 0<=y1<=y2<Y, 0<=z1<=z2<Z
func (b *BIT3D) QueryRange(x1, y1, z1, x2, y2, z2 int) int {
	return b.Query(x2, y2, z2) - b.Query(x1, y2, z2) -
		b.Query(x2, y1, z2) - b.Query(x2, y2, z1) +
		b.Query(x1, y1, z2) + b.Query(x1, y2, z1) +
		b.Query(x2, y1, z1) - b.Query(x1, y1, z1)
}
