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
//
// BooleanSquareMatrixDense
//
//
// !这里是bitset+four russians mathod的实现.对于稠密的矩阵,5000*5000的矩阵乘法需要500ms左右.
// complex: O(n^3 / wlogn)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// yuki1340()
	test()
	bs := NewBitset(70)
	bs.Set(10)
	bs.Set(15)
	fmt.Println(bs.Has(10))
	fmt.Println(bs._HasRange(10, 16))
	bs.Set(63)
	bs.Set(64)
	bs.Set(65)
	fmt.Println(bs._HasRange(62, 66))
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
	mat := NewBooleanSquareMatrix(n)
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
	mat := NewBooleanSquareMatrix(numCourses)
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
	// 5000*5000的矩阵乘法:275.2395ms
	// 2000*2000的传递闭包:356.363ms
	// ====================
	// 测试稀疏矩阵
	// 5000*5000的矩阵乘法:235.1503ms
	// 2000*2000的传递闭包:359.9047ms
	// ====================
	// 测试稠密矩阵
	// 5000*5000的矩阵乘法:230.2989ms
	// 2000*2000的传递闭包:377.9874ms

	mat := NewBooleanSquareMatrix(3)
	mat.Set(0, 0, true)
	mat.Set(0, 1, true)
	mat.Set(1, 2, true)
	mat.Set(1, 0, true)

	testRandom := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试随机矩阵")
		// !随机01矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		mat := NewBooleanSquareMatrix(N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				if rand.Intn(2) == 0 {
					mat.Set(i, j, true)
				}
			}
		}
		time1 := time.Now()
		Mul(mat, mat)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		mat = NewBooleanSquareMatrix(N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				if rand.Intn(2) == 0 {
					mat.Set(i, j, true)
				}
			}
		}
		time3 := time.Now()
		mat.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testSparse := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试稀疏矩阵")
		// !稀疏矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		mat := NewBooleanSquareMatrix(N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				if rand.Intn(10) == 0 {
					mat.Set(i, j, true)
				}
			}
		}
		time1 := time.Now()
		Mul(mat, mat)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		mat = NewBooleanSquareMatrix(N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				if rand.Intn(10) == 0 {
					mat.Set(i, j, true)
				}
			}
		}
		time3 := time.Now()
		mat.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testDense := func() {
		fmt.Println(strings.Repeat("=", 20))
		fmt.Println("测试稠密矩阵")
		// !稠密矩阵
		// 5000*5000的矩阵乘法
		N_5000 := 5000
		mat := NewBooleanSquareMatrix(N_5000)
		for i := 0; i < N_5000; i++ {
			for j := 0; j < N_5000; j++ {
				mat.Set(i, j, true)
			}
		}
		time1 := time.Now()
		Mul(mat, mat)
		time2 := time.Now()
		fmt.Println(fmt.Sprintf("5000*5000的矩阵乘法:%v", time2.Sub(time1)))

		// 2000*2000的传递闭包
		N_2000 := 2000
		mat = NewBooleanSquareMatrix(N_2000)
		for i := 0; i < N_2000; i++ {
			for j := 0; j < N_2000; j++ {
				mat.Set(i, j, true)
			}
		}
		time3 := time.Now()
		mat.TransitiveClosure()
		time4 := time.Now()
		fmt.Println(fmt.Sprintf("2000*2000的传递闭包:%v", time4.Sub(time3)))
	}

	testRandom()
	testSparse()
	testDense()
}

// trailing zero table
var _BSF [1e4 + 10]int

func init() {
	for i := range _BSF {
		_BSF[i] = bits.TrailingZeros(uint(i))
	}
}

// 布尔方阵.
type BooleanSquareMatrix struct {
	N  int
	bs []BitSet64
	dp []BitSet64 // 在计算矩阵乘法时用到
}

// n<=1e4.
func NewBooleanSquareMatrix(n int) *BooleanSquareMatrix {
	bs := make([]BitSet64, n)
	for i := range bs {
		bs[i] = NewBitset(n)
	}
	return &BooleanSquareMatrix{N: n, bs: bs}
}

// n<=1e4.
func Eye(n int) *BooleanSquareMatrix {
	res := NewBooleanSquareMatrix(n)
	for i := 0; i < n; i++ {
		res.bs[i].Set(i)
	}
	return res
}

func Pow(mat *BooleanSquareMatrix, k int) *BooleanSquareMatrix {
	return mat.Copy().IPow(k)
}

func Mul(mat1, mat2 *BooleanSquareMatrix) *BooleanSquareMatrix {
	return mat1.Copy().IMul(mat2)
}

func Add(mat1, mat2 *BooleanSquareMatrix) *BooleanSquareMatrix {
	return mat1.Copy().IAdd(mat2)
}

// (A + I)^n 是传递闭包.
//  Deprecated: 建议使用`O(n^3/64)`的Floyd-Warshall算法.
func (bm *BooleanSquareMatrix) TransitiveClosure() *BooleanSquareMatrix {
	n := bm.N
	newMat := Eye(n).IAdd(bm)
	newMat.IPow(n)
	return newMat
}

func (bm *BooleanSquareMatrix) IPow(k int) *BooleanSquareMatrix {
	res := Eye(bm.N)
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

// O(n^3/wlogn),这里logn指的是分块的大小.
func (bm *BooleanSquareMatrix) IMul(mat *BooleanSquareMatrix) *BooleanSquareMatrix {
	n := mat.N
	res := NewBooleanSquareMatrix(n)
	step := 8 // !理论最优是logn,实际取8效果最好(n为5000时)
	bm._initDpIfAbsent(step, n)
	dp := bm.dp
	bmBs := bm.bs
	matBs := mat.bs

	for l, r := 0, step; l != n; l, r = r, r+step {
		if r > n {
			r = n
		}

		for s := 1; s < (1 << step); s++ {
			bsf := _BSF[s]
			if l+bsf < n {
				dp[s] = Or(dp[s^(1<<bsf)], matBs[l+bsf]) // Xor => f2矩阵乘法
			} else {
				dp[s] = dp[s^(1<<bsf)]
			}
		}

		for i := 0; i != n; i++ {
			res.bs[i].IOr(dp[bmBs[i]._HasRange(l, r)]) // IXor => f2矩阵乘法
		}
	}

	bm.bs, res.bs = res.bs, bm.bs
	return res
}

func (bm *BooleanSquareMatrix) IAdd(mat *BooleanSquareMatrix) *BooleanSquareMatrix {
	for i := 0; i < bm.N; i++ {
		bm.bs[i].IOr(mat.bs[i])
	}
	return bm
}

func (bm *BooleanSquareMatrix) Copy() *BooleanSquareMatrix {
	bs := make([]BitSet64, bm.N)
	for i := range bs {
		bs[i] = bm.bs[i].Copy()
	}
	return &BooleanSquareMatrix{N: bm.N, bs: bs, dp: bm.dp}
}

func (bm *BooleanSquareMatrix) Get(row, col int) bool {
	return bm.bs[row].Has(col)
}

func (bm *BooleanSquareMatrix) Set(row, col int, b bool) {
	if b {
		bm.bs[row].Set(col)
	} else {
		bm.bs[row].Reset(col)
	}
}

// To 2D grid.
func (mat *BooleanSquareMatrix) String() string {
	n := mat.N
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if mat.Get(i, j) {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
			}
		}
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("BooleanSquareMatrix(%d,%d)\n", n, n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d ", grid[i][j]))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (mat *BooleanSquareMatrix) _initDpIfAbsent(step int, n int) {
	if mat.dp == nil {
		dp := make([]BitSet64, 1<<step)
		for i := range dp {
			dp[i] = NewBitset(n)
		}
		mat.dp = dp
	}
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

// ![l,r) 范围内数与bitset相交的数.r-l<64.
// eg: l=10, r=16
//    [15,14,13,12,11,10] & Bitset(10,11,13,15) => 101011
func (b BitSet64) _HasRange(l, r int) uint64 {
	posL, shiftL := l>>6, l&63
	posR, shiftR := r>>6, r&63
	maskL, maskR := ^(^uint64(0) << shiftL), ^(^uint64(0) << shiftR) // 低位全1
	if posL == posR {
		return (b[posL] & maskR) >> shiftL
	}
	// posL+1 == posR
	return (b[posL] & ^maskL)>>shiftL | (b[posR]&maskR)<<(64-shiftL)
}
