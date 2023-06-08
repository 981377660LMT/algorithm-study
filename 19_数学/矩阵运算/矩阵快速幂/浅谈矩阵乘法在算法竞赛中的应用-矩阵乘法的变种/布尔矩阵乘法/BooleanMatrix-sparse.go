// !布尔矩阵乘法(Boolean Matrix Multiplication, BMM)
// 输入和输出矩阵的元素均为布尔值。
// !按矩阵乘法的公式运算时，可以把“乘”看成and，把“加”看成or
// 对矩阵乘法 C[i][j] |= A[i][k] & B[k][j], 它的一个直观意义是把A的行和B的列看成集合，
// A的第i行包含元素k当且仅当A[i][k]=1。
// B的第j列包含元素k当且仅当B[k][j]=1。
// !那么C[i][j]代表A的第i行和B的第j列是否包含公共元素。
//
// 一个应用是传递闭包(Transitive Closure)的加速计算。
//
// https://zhuanlan.zhihu.com/p/631804105
// !https://blog.csdn.net/qq_42101694/article/details/121227383
// https://github.com/spaghetti-source/algorithm/blob/4fdac8202e26def25c1baf9127aaaed6a2c9f7c7/math/matrix_bool.cc#L4
//
// BooleanMatrixSparse
//
// !这里是bitset的实现.当输入矩阵比较`稀疏`时可以跑得非常快(5000*5000 => 200ms).
//
// complexity:O(n^3/w)

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// yuki1340()
	test()
}

// https://yukicoder.me/problems/no/1340
// 给定一个n个点m条边的有向图，求t步后可能所在的顶点个数(每一步必须移动到一个相邻点).
// n<=100 m<=1e4 t<=1e18
func yuki1340() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, t int
	fmt.Fscan(in, &n, &m, &t)
	mat := NewBooleanMatrix(n, n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		mat.Set(a, b, true)
	}
	mat.IPow(t)
	res := 0
	for i := 0; i < n; i++ {
		if mat.Get(0, i) {
			res++
		}
	}
	fmt.Fprintln(out, res)
}

// https://leetcode.cn/problems/course-schedule-iv/
func checkIfPrerequisite(numCourses int, prerequisites [][]int, queries [][]int) []bool {
	mat := NewBooleanMatrix(numCourses, numCourses)
	for _, p := range prerequisites {
		mat.Set(p[0], p[1], true)
	}
	trans := mat.TransitiveClosure()
	res := make([]bool, len(queries))
	for i, q := range queries {
		res[i] = trans.Get(q[0], q[1])
	}
	return res
}

func test() {
	// ====================
	// 测试随机矩阵
	// 5000*5000的矩阵乘法:828.835ms
	// 2000*2000的传递闭包:1.3121475s
	// ====================
	// 测试稀疏矩阵
	// 5000*5000的矩阵乘法:202.4407ms
	// 2000*2000的传递闭包:1.21563s
	// ====================
	// 测试稠密矩阵
	// 5000*5000的矩阵乘法:1.1976128s
	// 2000*2000的传递闭包:1.2487846s

	mat := NewBooleanMatrix(3, 3)
	mat.Set(0, 0, true)
	mat.Set(0, 1, true)
	mat.Set(1, 2, true)
	mat.Set(1, 0, true)
	fmt.Println(mat)
	fmt.Println(Mul(Mul(mat, mat), mat), Pow(mat, 3))
	fmt.Println(Eye(8))

	testRandom := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试随机矩阵")
		// !随机01矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		eye := Eye(N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				if rand.Intn(2) == 0 {
					eye.Set(i, j, true)
				}
			}
		}
		time1 := time.Now()
		Mul(eye, eye)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		eye = Eye(N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				if rand.Intn(2) == 0 {
					eye.Set(i, j, true)
				}
			}
		}
		time3 := time.Now()
		eye.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testSparse := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试稀疏矩阵")
		// !稀疏矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		eye := NewBooleanMatrix(N_5000, N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				if rand.Intn(10) == 0 {
					eye.Set(i, j, true)
				}
			}
		}
		time1 := time.Now()
		Mul(eye, eye)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		eye = NewBooleanMatrix(N_2000, N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				if rand.Intn(10) == 0 {
					eye.Set(i, j, true)
				}
			}
		}
		time3 := time.Now()
		eye.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testDense := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试稠密矩阵")
		// !稠密矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		eye := NewBooleanMatrix(N_5000, N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				eye.Set(i, j, true)
			}
		}
		time1 := time.Now()
		Mul(eye, eye)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		eye = NewBooleanMatrix(N_2000, N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				eye.Set(i, j, true)
			}
		}
		time3 := time.Now()
		eye.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testRandom()
	testSparse()
	testDense()

}

type BooleanMatrix struct {
	ROW, COL int
	bs       []BitSet64
}

func NewBooleanMatrix(row, col int) *BooleanMatrix {
	bs := make([]BitSet64, row)
	for i := range bs {
		bs[i] = NewBitset(col)
	}
	return &BooleanMatrix{ROW: row, COL: col, bs: bs}
}

func Eye(n int) *BooleanMatrix {
	res := NewBooleanMatrix(n, n)
	for i := 0; i < n; i++ {
		res.bs[i].Set(i)
	}
	return res
}

func Pow(mat *BooleanMatrix, k int) *BooleanMatrix {
	return mat.Copy().IPow(k)
}

func Mul(mat1, mat2 *BooleanMatrix) *BooleanMatrix {
	return mat1.Copy().IMul(mat2)
}

func Add(mat1, mat2 *BooleanMatrix) *BooleanMatrix {
	return mat1.Copy().IAdd(mat2)
}

// (A + I)^n 是传递闭包.
//  Deprecated: 建议使用`O(n^3/64)`的Floyd-Warshall算法.
func (bm *BooleanMatrix) TransitiveClosure() *BooleanMatrix {
	if bm.ROW != bm.COL {
		panic("Not a square matrix")
	}
	n := bm.ROW
	newMat := Eye(n).IAdd(bm)
	newMat.IPow(n)
	return newMat
}

func (bm *BooleanMatrix) IPow(k int) *BooleanMatrix {
	res := Eye(bm.ROW)
	for k > 0 {
		if k&1 == 1 {
			res.IMul(bm)
		}
		bm.IMul(bm)
		k >>= 1
	}
	res.bs, bm.bs = bm.bs, res.bs
	return bm
}

// !O(n^3)/w
func (bm *BooleanMatrix) IMul(mat *BooleanMatrix) *BooleanMatrix {
	row, col := bm.ROW, mat.COL
	res := NewBooleanMatrix(row, col)
	mbs := mat.bs
	for i := 0; i < row; i++ {
		rowBs := bm.bs[i]
		resBs := res.bs[i]
		for j := 0; j < bm.COL; j++ {
			if rowBs.Has(j) {
				resBs.IOr(mbs[j])
				// resBs.IXOr(mbs[j])  // f2上的矩阵乘法
			}
		}
	}
	bm.bs, res.bs = res.bs, bm.bs
	return bm
}

func (bm *BooleanMatrix) IAdd(mat *BooleanMatrix) *BooleanMatrix {
	for i := 0; i < bm.ROW; i++ {
		bm.bs[i].IOr(mat.bs[i])
	}
	return bm
}

func (bm *BooleanMatrix) Copy() *BooleanMatrix {
	bs := make([]BitSet64, bm.ROW)
	for i := range bs {
		bs[i] = bm.bs[i].Copy()
	}
	return &BooleanMatrix{ROW: bm.ROW, COL: bm.COL, bs: bs}
}

func (bm *BooleanMatrix) Get(row, col int) bool {
	return bm.bs[row].Has(col)
}

func (bm *BooleanMatrix) Set(row, col int, b bool) {
	if b {
		bm.bs[row].Set(col)
	} else {
		bm.bs[row].Reset(col)
	}
}

// To 2D grid.
func (mat *BooleanMatrix) String() string {
	grid := make([][]int, mat.ROW)
	for i := 0; i < mat.ROW; i++ {
		grid[i] = make([]int, mat.COL)
		for j := 0; j < mat.COL; j++ {
			if mat.Get(i, j) {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
			}
		}
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("BooleanMatrix(%d,%d)\n", mat.ROW, mat.COL))
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			sb.WriteString(fmt.Sprintf("%d ", grid[i][j]))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

type BitSet64 []uint64

func NewBitset(n int) BitSet64 { return make(BitSet64, n>>6+1) } // (n+_w-1)>>6

func (b BitSet64) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b BitSet64) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b BitSet64) Set(p int)      { b[p>>6] |= 1 << (p & 63) }
func (b BitSet64) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) }

func (b BitSet64) Get(p int) int { return int(b[p>>6] >> (p & 63) & 1) }

func (b BitSet64) Copy() BitSet64 {
	res := make(BitSet64, len(b))
	copy(res, b)
	return res
}

// 将 c 的元素合并进 b
func (b BitSet64) IOr(c BitSet64) BitSet64 {
	for i, v := range c {
		b[i] |= v
	}
	return b
}

// !f2上的加法
func (b BitSet64) IXOr(c BitSet64) {
	for i, v := range c {
		b[i] ^= v
	}
}

func Or(a, b BitSet64) BitSet64 {
	res := make(BitSet64, len(a))
	for i, v := range a {
		res[i] = v | b[i]
	}
	return res
}

func Xor(a, b BitSet64) BitSet64 {
	res := make(BitSet64, len(a))
	for i, v := range a {
		res[i] = v ^ b[i]
	}
	return res
}
