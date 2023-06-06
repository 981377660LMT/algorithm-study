package main

import (
	"fmt"
	"math/bits"
	"time"
)

func main() {
	mat1 := [][]uint8{{1, 0, 1}, {0, 1, 0}, {1, 0, 1}}
	mat2 := [][]uint8{{1, 0, 1}, {0, 1, 0}, {1, 0, 1}}
	res := ZeroOneSquareMatrix(mat1, mat2)
	fmt.Println(res)

	for _, n := range [5]int{1000, 2000, 3000, 4000, 5000} {
		fmt.Println(n)

		mat1 = make([][]uint8, n)
		mat2 = make([][]uint8, n)
		for i := 0; i < n; i++ {
			mat1[i] = make([]uint8, n)
			mat2[i] = make([]uint8, n)
			for j := 0; j < n; j++ {
				mat1[i][j] = uint8(i & 1)
				mat2[i][j] = uint8(j & 1)
			}
		}

		time1 := time.Now()
		res = ZeroOneSquareMatrix(mat1, mat2)
		fmt.Println(time.Since(time1))
	}
}

// !01矩阵乘法(方阵).C = A * B.
// 输入矩阵中的元素只包含0或1,然后进行正常的矩阵乘法.
// 也就是说输出矩阵元素值的范围为0~n.
// 一个直观意义是, C[i][j]代表A的第i行和B的第j列的`公共元素个数`.
//  mat1, mat2: 方阵A, B.边长n<=5000.
//  res: 方阵C.每个元素的值范围为0~n(在uint16范围内).
//  时间复杂度为 `O(n^3/64)`.
//  `1000*1000 => 35ms`
//  `2000*2000 => 200ms`
//  `3000*3000 => 650ms
//  `4000*4000 => 1.5s`
//  `5000*5000 => 2.9s`.
func ZeroOneSquareMatrix(mat1, mat2 [][]uint8) (res [][]uint16) {
	if len(mat1) == 0 || len(mat2) == 0 {
		return
	}
	if len(mat1) != len(mat2) {
		panic("mat1.length != mat2.length")
	}

	n := len(mat1)
	blockCount := (n + 63) >> 6
	block1 := make([]uint64, blockCount*n)
	block2 := make([]uint64, blockCount*n)
	for i := 0; i < n; i++ {
		cache1, cache2 := mat1[i], mat2[i]
		for j := 0; j < n; j++ {
			block1[i*blockCount+(j>>6)] |= uint64(cache1[j]) << (j & 63)
			block2[j*blockCount+(i>>6)] |= uint64(cache2[j]) << (i & 63)
		}
	}

	res = make([][]uint16, n)
	for i := 0; i < n; i++ {
		row := make([]uint16, n)
		res[i] = row
		for j := 0; j < n; j++ {
			sum := 0
			for k := 0; k < blockCount; k++ {
				sum += bits.OnesCount64(block1[i*blockCount+k] & block2[j*blockCount+k])
			}
			row[j] = uint16(sum)
		}
	}

	return
}
