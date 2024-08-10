package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

// N 頂点
// M 辺の単純無向グラフが与えられます。
// i 番目の辺は頂点
// u
// i
// ​
//   と
// v
// i
// ​
//   を双方向に結びます。

// このグラフの各頂点に
// 1 以上
// 2
// 60
//   未満の整数を書き込む方法であって、次の条件を満たすものが存在するか判定し、存在するならば一つ構築してください。

// 次数が
// 1 以上のすべての頂点
// v について、隣接する頂点 (
// v 自身は含まない) に書き込まれている数の総 XOR が
// 0 となる
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, m := io.NextInt(), io.NextInt()
	adjList := make([]BitSet64, n)
	for i := 0; i < n; i++ {
		adjList[i] = NewBitset(n)
	}
	for i := 0; i < m; i++ {
		u, v := io.NextInt(), io.NextInt()
		u, v = u-1, v-1
		adjList[u].Set(v)
		adjList[v].Set(u)
	}

	res1, res2, res3 := SystemOfLinearEquations(adjList, make([]bool, n), n)
	fmt.Println(res1, res2, res3)

}

// 高斯-约当消元法.将矩阵原地修改为最简行阶梯矩阵,即对角线上全是1，其他地方的值都为0的对角阵(diagonal matrix).
//
//	即将Ax=b变为Ux=c.
//	col: 矩阵的列数
//	时间复杂度: O(HW + HW rank(M) / 64)
//	Verified: abc276_h (2000 x 8000)
func GaussJordan(mat []BitSet64, col int) []BitSet64 {
	row := len(mat)
	c := 0
	for h := 0; h < row && c < col; h, c = h+1, c+1 {
		pivot := -1
		for j := h; j < row; j++ {
			if mat[j].Has(c) {
				pivot = j
				break
			}
		}
		if pivot == -1 {
			h--
			continue
		}
		mat[pivot], mat[h] = mat[h], mat[pivot]
		for hh := 0; hh < row; hh++ {
			if hh != h && mat[hh].Has(c) {
				mat[hh].IXOr(mat[h])
			}
		}
	}
	return mat
}

// 01矩阵的秩.
func CalRank(mat []BitSet64, col int) int {
	for h := len(mat) - 1; h >= 0; h-- {
		j := 0
		for j < col && !mat[h].Has(j) {
			j++
		}
		if j < col {
			return h + 1
		}
	}
	return 0
}

// F2矩阵乘法.
func MatMul(mat1, mat2 []BitSet64) []BitSet64 {
	row, col := len(mat1), len(mat2)
	res := make([]BitSet64, row)
	for i := 0; i < row; i++ {
		res[i] = NewBitset(col)
		for j := 0; j < col; j++ {
			if mat1[i].Has(j) {
				res[i].IXOr(mat2[j])
			}
		}
	}
	return res
}

// F2矩阵快速幂.返回一个新的矩阵.
func MatPow(mat []BitSet64, k int) []BitSet64 {
	n := len(mat)
	res := make([]BitSet64, n)
	for i := 0; i < n; i++ {
		res[i] = NewBitset(n)
		res[i].Set(i)
	}
	matCopy := make([]BitSet64, n)
	for i := 0; i < n; i++ {
		matCopy[i] = mat[i].Copy()
	}
	mat = matCopy
	for ; k > 0; k >>= 1 {
		if k&1 == 1 {
			res = MatMul(res, mat)
		}
		mat = MatMul(mat, mat)
	}
	return res
}

// 求解f2上的线性方程组Ax=b.
//
//	返回值: (一组解, 方程组的自由变量, 是否有解)
//	时间复杂度: O(HW + HW rank(A) / 64 + W^2 len(freedoms))
func SystemOfLinearEquations(A []BitSet64, b []bool, col int) (solution0 BitSet64, freedoms []BitSet64, ok bool) {
	row := len(A)
	if row != len(b) {
		panic("len(A) != len(b)")
	}

	M := make([]BitSet64, row)
	for i := 0; i < row; i++ {
		bs := NewBitset(col + 1)
		for j := 0; j < col; j++ {
			if A[i].Has(j) {
				bs.Set(j)
			}
		}
		if b[i] {
			bs.Set(col)
		}
		M[i] = bs
	}

	M = GaussJordan(M, col+1)
	ss := make([]int, col)
	for i := range ss {
		ss[i] = -1
	}
	var ssNonnegJs []int
	for i := 0; i < row; i++ {
		j := 0
		for j <= col && !M[i].Has(j) {
			j++
		}
		if j == col {
			return
		}
		if j < col {
			ssNonnegJs = append(ssNonnegJs, j)
			ss[j] = i
		}
	}

	ok = true
	solution0 = NewBitset(col)
	for j := 0; j < col; j++ {
		if ss[j] == -1 {
			// This part may require W^2 space complexity in output
			d := NewBitset(col)
			d.Set(j)
			for _, jj := range ssNonnegJs {
				if M[ss[jj]].Has(j) {
					d.Set(jj)
				}
			}
			freedoms = append(freedoms, d)
		} else {
			if M[ss[j]].Has(col) {
				solution0.Set(j)
			}
		}
	}

	return
}

type BitSet64 []uint64

func NewBitset(n int) BitSet64 { return make(BitSet64, n>>6+1) } // (n+_w-1)>>6

func (b BitSet64) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b BitSet64) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b BitSet64) Set(p int)      { b[p>>6] |= 1 << (p & 63) }
func (b BitSet64) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) }

func (b BitSet64) Copy() BitSet64 {
	res := make(BitSet64, len(b))
	copy(res, b)
	return res
}

// !f2上的加法
func (b BitSet64) IXOr(c BitSet64) {
	for i, v := range c {
		b[i] ^= v
	}
}

func Xor(a, b BitSet64) BitSet64 {
	res := make(BitSet64, len(a))
	for i, v := range a {
		res[i] = v ^ b[i]
	}
	return res
}
