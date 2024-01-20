package main

import "fmt"

func main() {
	mat := make([][][]int, 2)
	mat[0] = [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	mat[1] = [][]int{
		{7, 8, 9},
		{10, 11, 12},
	}
	S := NewPreSum3D(mat)
	fmt.Println(S.Query(0, 0, 0, 1, 1, 1))
}

type E = int

func (*PreSum3D) e() E        { return 0 }
func (*PreSum3D) op(a, b E) E { return a + b }
func (*PreSum3D) inv(a E) E   { return -a }

// 三维前缀和.
type PreSum3D struct {
	xSize, ySize, zSize int
	preSum              [][][]E
}

func NewPreSum3D(mat [][][]E) *PreSum3D {
	res := &PreSum3D{}
	xSize, ySize, zSize := len(mat), len(mat[0]), len(mat[0][0])
	preSum := make([][][]E, xSize+1)
	for x := range preSum {
		preSum[x] = make([][]E, ySize+1)
		for y := range preSum[x] {
			row := make([]E, zSize+1)
			preSum[x][y] = row
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = mat[x-1][y-1][z-1]
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x-1][y][z])
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x][y-1][z])
			}
		}
	}

	for x := 1; x <= xSize; x++ {
		for y := 1; y <= ySize; y++ {
			for z := 1; z <= zSize; z++ {
				preSum[x][y][z] = res.op(preSum[x][y][z], preSum[x][y][z-1])
			}
		}
	}

	res.xSize, res.ySize, res.zSize = xSize, ySize, zSize
	res.preSum = preSum
	return res
}

// 查询 sum(A[x1:x2+1][y1:y2+1][z1:z2+1])的值
// 0 <= x1 <= x2 < xSize
// 0 <= y1 <= y2 < ySize
// 0 <= z1 <= z2 < zSize
func (ps *PreSum3D) Query(x1, y1, z1, x2, y2, z2 int) E {
	res := ps.preSum[x2+1][y2+1][z2+1]
	res = ps.op(res, ps.inv(ps.preSum[x1][y2+1][z2+1]))
	res = ps.op(res, ps.inv(ps.preSum[x2+1][y1][z2+1]))
	res = ps.op(res, ps.inv(ps.preSum[x2+1][y2+1][z1]))
	res = ps.op(res, ps.preSum[x1][y1][z2+1])
	res = ps.op(res, ps.preSum[x1][y2+1][z1])
	res = ps.op(res, ps.preSum[x2+1][y1][z1])
	res = ps.op(res, ps.inv(ps.preSum[x1][y1][z1]))
	return res
}
