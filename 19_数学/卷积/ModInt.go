package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

var iost *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp io.Reader, wfp io.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (i *Iost) Text() string {
	if !i.Scanner.Scan() {
		panic("scan failed")
	}
	return i.Scanner.Text()
}
func (i *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (i *Iost) GetNextInt() int                   { return i.Atoi(i.Text()) }
func (i *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (i *Iost) GetNextInt64() int64               { return i.Atoi64(i.Text()) }
func (i *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (i *Iost) GetNextFloat64() float64           { return i.Atof64(i.Text()) }
func (i *Iost) Print(x ...interface{})            { fmt.Fprint(i.Writer, x...) }
func (i *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(i.Writer, s, x...) }
func (i *Iost) Println(x ...interface{})          { fmt.Fprintln(i.Writer, x...) }

func main() {
	fp := os.Stdin
	wfp := os.Stdout

	iost = NewIost(fp, wfp)
	defer func() {
		iost.Writer.Flush()
	}()
	solve()
}

func solve() {
	// https://atcoder.jp/contests/abc291/tasks/abc291_d
	// dp[i][0/1] : 当前看到第i个点, 之前那张牌是否翻转

	n := iost.GetNextInt()
	cards := make([][2]int, n)
	for i := 0; i < n; i++ {
		cards[i][0] = iost.GetNextInt()
		cards[i][1] = iost.GetNextInt()
	}

	dp := make([][2]ModInt, n)
	dp[0][0] = 1
	dp[0][1] = 1
	for i := 1; i < n; i++ {
		for pre := 0; pre < 2; pre++ {
			for cur := 0; cur < 2; cur++ {
				if cards[i-1][pre] != cards[i][cur] {
					dp[i][cur].IAdd(dp[i-1][pre])
				}
			}
		}
	}

	iost.Println(dp[n-1][0].Add(dp[n-1][1]))
}

const MOD = 998244353

type ModInt int64

func (m ModInt) Add(x ModInt) ModInt {
	return (m + x).mod()
}
func (m *ModInt) IAdd(x ModInt) {
	*m = m.Add(x)
}
func (m ModInt) Sub(x ModInt) ModInt {
	return (m - x).mod()
}
func (m *ModInt) ISub(x ModInt) {
	*m = m.Sub(x)
}
func (m ModInt) Mul(x ModInt) ModInt {
	return (m * x).mod()
}
func (m *ModInt) IMul(x ModInt) {
	*m = m.Mul(x)
}
func (m ModInt) Div(x ModInt) ModInt {
	return m.Mul(x.Inv())
}
func (m *ModInt) IDiv(x ModInt) {
	*m = m.Div(x)
}
func (m ModInt) Pow(n ModInt) ModInt {
	p := ModInt(1)
	for n > 0 {
		if n&1 == 1 {
			p.IMul(m)
		}
		m.IMul(m)
		n >>= 1
	}
	return p
}
func (m ModInt) Inv() ModInt {
	return m.Pow(ModInt(0).Sub(2))
}
func (m ModInt) mod() ModInt {
	m %= MOD
	if m < 0 {
		m += MOD
	}
	return m
}
