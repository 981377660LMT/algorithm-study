// https://judge.yosupo.jp/problem/convolution_mod_1000000007
// 给定一个n次多项式F(x)，和一个m次多项式G(x)。
// 请求出 F(x)和G(x)的卷积。
// 1≤N,M≤524,288,nums[i]<=1e9+7

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
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
func (io *Iost) Input() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Input()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Input()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Input()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, m := io.NextInt(), io.NextInt()

	poly1 := make([]int, n) // 从低到高表示F(x)的系数
	poly2 := make([]int, m) // 从低到高表示G(x)的系数
	for i := 0; i < n; i++ {
		poly1[i] = io.NextInt()
	}
	for i := 0; i < m; i++ {
		poly2[i] = io.NextInt()
	}

	const MOD int = 1e9 + 7
	conv := Convolution(poly1, poly2, MOD)
	for i := 0; i < n+m-1; i++ {
		fmt.Fprint(out, conv[i], " ")
	}
}

const (
	MOD1, ROOT1 = 167772161, 3
	MOD2, ROOT2 = 469762049, 3
	MOD3, ROOT3 = 1224736769, 3
)

func Convolution(a, b []int, mod int) []int {
	x := nttConvolve(a, b, MOD1, ROOT1)
	y := nttConvolve(a, b, MOD2, ROOT2)
	z := nttConvolve(a, b, MOD3, ROOT3)

	inv1_2 := pow(MOD1, MOD2-2, MOD2)
	inv12_3 := pow(MOD1*MOD2, MOD3-2, MOD3)
	mod12 := MOD1 * MOD2 % mod
	res := make([]int, len(x))
	for i := range x {
		v1 := ((y[i]-x[i])*inv1_2%MOD2 + MOD2) % MOD2
		v2 := ((z[i]-(x[i]+MOD1*v1)%MOD3)*inv12_3%MOD3 + MOD3) % MOD3
		res[i] = (x[i] + MOD1*v1 + mod12*v2) % mod
	}
	return res[:len(a)+len(b)-1]
}

func nttConvolve(a, b []int, mod, root int) []int {
	n := 1 << (bits.Len(uint(len(a) + len(b) - 1)))
	h := bits.Len(uint(n)) - 1
	a = append(a, make([]int, n-len(a))...)
	b = append(b, make([]int, n-len(b))...)
	ntt(a, h, mod, root)
	ntt(b, h, mod, root)
	for i := range a {
		a[i] = a[i] * b[i] % mod
	}
	intt(a, h, mod, root)
	return a
}

func ntt(a []int, h, mod, root int) {
	roots := make([]int, h+1)
	for i := range roots {
		roots[i] = pow(root, (mod-1)>>i, mod)
	}

	for i := 0; i < h; i++ {
		m := 1 << (h - i - 1)
		for j := 0; j < (1 << i); j++ {
			w := 1
			s := 2 * m * j
			for k := 0; k < m; k++ {
				a[s+k], a[s+k+m] = (a[s+k]+a[s+k+m])%mod, ((a[s+k]-a[s+k+m])*w%mod+mod)%mod
				w = w * roots[h-i] % mod
			}
		}
	}
}

func intt(a []int, h, mod, root int) {
	iRoots := make([]int, h+1)
	for i := range iRoots {
		iRoots[i] = pow(pow(root, (mod-1)>>i, mod), mod-2, mod)
	}
	for i := 0; i < h; i++ {
		m := 1 << i
		for j := 0; j < 1<<(h-i-1); j++ {
			w := 1
			s := 2 * m * j
			for k := 0; k < m; k++ {
				a[s+k], a[s+k+m] = (a[s+k]+a[s+k+m]*w)%mod, ((a[s+k]-a[s+k+m]*w)%mod+mod)%mod
				w = w * iRoots[i+1] % mod
			}
		}
	}
	inv := pow(1<<h, mod-2, mod)
	for i := range a {
		a[i] = a[i] * inv % mod
	}
}

func pow(base int, exp int, mod int) (res int) {
	base %= mod
	res = 1
	for ; exp > 0; exp >>= 1 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
	}
	return
}
