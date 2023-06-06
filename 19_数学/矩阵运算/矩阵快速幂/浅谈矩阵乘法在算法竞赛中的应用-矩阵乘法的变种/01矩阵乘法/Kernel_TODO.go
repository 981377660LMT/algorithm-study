package main

import (
	"fmt"
	"math/bits"
)

func main() {
	mat1 := [][]uint8{{1, 0, 1}, {0, 1, 0}, {1, 0, 1}}
	mat2 := [][]uint8{{1, 0, 1}, {0, 1, 0}, {1, 0, 1}}
	res := ZeroOneSquareMatrix(mat1, mat2)

	// 不对,应该是 [[2 0 2] [0 1 0] [2 0 2]]
	fmt.Println(res) // [[2 0 2] [0 0 0] [0 0 0]]

}

// pool
const (
	_MAX_N       = 2880
	_BLOCK_COUNT = (_MAX_N + 63) >> 6
	_KERNEL_W    = 16
	_KERNEL_H    = 16
)

var (
	_BLOCK_1 = [_MAX_N * _BLOCK_COUNT]uint64{}
	_BLOCK_2 = [_MAX_N * _BLOCK_COUNT]uint64{}
	_RES     = [_MAX_N * _MAX_N]uint16{}
)

func ZeroOneSquareMatrix(mat1, mat2 [][]uint8) [][]uint16 {
	for i := range _BLOCK_1 {
		_BLOCK_1[i] = 0
		_BLOCK_2[i] = 0
	}

	n := len(mat1)
	for i := 0; i < n; i++ {
		cache1, cache2 := mat1[i], mat2[i]
		for j := 0; j < n; j++ {
			_BLOCK_1[i*_BLOCK_COUNT+(j>>6)] |= uint64(cache1[j]) << (j & 63)
			_BLOCK_2[(i>>6)*_MAX_N+j] |= uint64(cache2[j]) << (i & 63)
		}
	}

	processKernel := func(x, y, l, r int) {
		tmp := [_KERNEL_W][_KERNEL_H]uint16{}
		for k := l; k < r; k++ {
			offset := k*_MAX_N + y
			for i := 0; i < _KERNEL_W; i++ {
				mask := _BLOCK_1[(x+i)*_BLOCK_COUNT+k]
				for j := 0; j < _KERNEL_H; j++ {
					tmp[i][j] += uint16(bits.OnesCount64(mask & _BLOCK_2[offset+j]))
				}
			}
		}

		for i := 0; i < _KERNEL_W; i++ {
			for j := 0; j < _KERNEL_H; j++ {
				_RES[(x+i)*_MAX_N+y+j] += tmp[i][j]
			}
		}
	}

	const s3 int = 256
	const s2 int = 16
	const s1 int = 16
	for i3 := 0; i3 < n; i3 += s3 {
		for i2 := 0; i2 < n; i2 += s2 {
			for i1 := 0; i1 < _BLOCK_COUNT; i1 += s1 {
				for x := i2; x < min(i2+s2, n); x += _KERNEL_W {
					for y := i3; y < min(i3+s3, n); y += _KERNEL_H {
						processKernel(x, y, i1, min(i1+s1, _BLOCK_COUNT))
					}
				}
			}
		}
	}

	res := make([][]uint16, n)
	for i := 0; i < n; i++ {
		row := make([]uint16, n)
		res[i] = row
		for j := 0; j < n; j++ {
			row[j] = _RES[i*n+j]
		}
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
