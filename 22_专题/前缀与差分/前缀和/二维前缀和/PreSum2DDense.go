// 广义的二维前缀和

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://yukicoder.me/problems/no/1141
// 每次操作将第r行和第c列所有格子涂黑
// 求剩下的格子里所有数的乘积模1e9+7
// ! 注意幺元1和0不兼容,需要记录0的个数
func solve(grid [][]int, ops [][2]int) []int {
	ROW, COL := len(grid), len(grid[0])
	leaves := make([][]E, ROW)
	for i := 0; i < ROW; i++ {
		leaves[i] = make([]E, COL)
		for j := 0; j < COL; j++ {
			if grid[i][j] == 0 {
				leaves[i][j] = E{1, 1}
			} else {
				leaves[i][j] = E{grid[i][j], 0}
			}
		}
	}

	P := NewPreSum2DDense(leaves)
	res := make([]int, 0, len(ops))
	for _, op := range ops {
		r, c := op[0], op[1]
		res1 := P.Query(0, r, 0, c)
		res2 := P.Query(r+1, ROW, 0, c)
		res3 := P.Query(0, r, c+1, COL)
		res4 := P.Query(r+1, ROW, c+1, COL)
		tmp := P.op(P.op(res1, res2), P.op(res3, res4))
		if tmp.zero > 0 {
			res = append(res, 0)
		} else {
			res = append(res, tmp.mul)
		}
	}
	return res
}

const MOD int = 1e9 + 7

type E = struct{ mul, zero int }

func (*PreSum2DDense) e() E        { return E{1, 0} }
func (*PreSum2DDense) op(a, b E) E { return E{a.mul * b.mul % MOD, a.zero + b.zero} }
func (*PreSum2DDense) inv(a E) E   { return E{pow(a.mul, MOD-2, MOD), -a.zero} }
func pow(base, exp, mod int) int {
	base %= mod
	res := 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return res
}

type PreSum2DDense struct {
	row, col int
	data     []E
}

func NewPreSum2DDense(matrix [][]E) *PreSum2DDense {
	res := &PreSum2DDense{}
	row := len(matrix)
	col := 0
	if row > 0 {
		col = len(matrix[0])
	}
	data := make([]E, row*col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			data[i*col+j] = res.e()
		}
	}

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			k := i*col + j
			if j == 0 {
				data[k] = matrix[i][j]
			} else {
				data[k] = res.op(data[k-1], matrix[i][j]) // 行
			}
		}
	}

	for i := col; i < row*col; i++ {
		data[i] = res.op(data[i-col], data[i]) // 列
	}

	res.row = row
	res.col = col
	res.data = data
	return res
}

// [x1,x2) x [y1,y2)
//
//	0 <= x1 <= x2 <= row
//	0 <= y1 <= y2 <= col
func (p *PreSum2DDense) Query(x1, x2, y1, y2 int) E {
	if x2 == 0 || y2 == 0 {
		return p.e()
	}
	x1, x2, y1, y2 = x1-1, x2-1, y1-1, y2-1
	var a, b, c, d E
	if x1 >= 0 && y1 >= 0 {
		a = p.data[x1*p.col+y1]
	} else {
		a = p.e()
	}
	if x1 >= 0 && y2 >= 0 {
		b = p.data[x1*p.col+y2]
	} else {
		b = p.e()
	}
	if x2 >= 0 && y1 >= 0 {
		c = p.data[x2*p.col+y1]
	} else {
		c = p.e()
	}
	if x2 >= 0 && y2 >= 0 {
		d = p.data[x2*p.col+y2]
	} else {
		d = p.e()
	}
	return p.op(p.op(a, d), p.inv(p.op(b, c)))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var row, col int
	fmt.Fscan(in, &row, &col)
	grid := make([][]int, row)
	for i := 0; i < row; i++ {
		grid[i] = make([]int, col)
		for j := 0; j < col; j++ {
			fmt.Fscan(in, &grid[i][j])
		}
	}

	var q int
	fmt.Fscan(in, &q)
	ops := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &ops[i][0], &ops[i][1])
		ops[i][0]--
		ops[i][1]--
	}

	res := solve(grid, ops)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}
