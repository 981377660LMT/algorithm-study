// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=0350
// 小数到分数

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var decimal string
	fmt.Fscan(in, &decimal)
	num, deno := DecimalToFraction(decimal)
	res := fmt.Sprintf("%d/%d", num, deno)
	fmt.Fprintln(out, res)
}

// eg: 输入0.1(6) 输出 1/6
//     输入5.2(143), 输出 52091/9990
func DecimalToFraction(decimal string) (numerator, denominator int) {
	if decimal[len(decimal)-1] != ')' {
		decimal += "(0)"
	}
	pow := [17]int{1} // 循环长度不超过17
	for i := 0; i < 16; i++ {
		pow[i+1] = 10 * pow[i]
	}
	p := strings.Index(decimal, ".")
	a := decimal[:p]
	b := decimal[p+1:]
	res := _FT(0, 1)
	for i := 0; i < len(a); i++ {
		x := int(a[i] - '0')
		x *= pow[len(a)-1-i]
		res.IAdd(_FT(x, 1))
	}
	p = strings.Index(b, "(")
	c := b[p+1:]
	b = b[:p]
	for i := 0; i < len(b); i++ {
		b_ := int(b[i] - '0')
		res.IAdd(_FT(b_, pow[i+1]))
	}
	cf := _FT(1, pow[len(b)])
	c = c[:len(c)-1]
	num, _ := strconv.Atoi(c)
	res.IAdd(cf.Mul(_FT(num, pow[len(c)]-1)))
	return res.a, res.b
}

type F struct {
	a int // 分子, numerator
	b int // 分母, denominator，在内部恒为正数
}

func _FT(a, b int) F {
	if b == 0 {
		panic(fmt.Sprintf("ZeroDivisionError: Fraction(%d, 0)", a))
	}
	if b != 1 {
		gcd_ := gcd(a, b)
		if gcd_ != 1 {
			a /= gcd_
			b /= gcd_
		}
		if b < 0 {
			a = -a
			b = -b
		}
	}
	return F{a, b}
}

func (this *F) Add(other F) F {
	return _FT(this.a*other.b+this.b*other.a, this.b*other.b)
}

func (this *F) IAdd(other F) *F {
	*this = this.Add(other)
	return this
}

func (this *F) Mul(other F) F {
	return _FT(this.a*other.a, this.b*other.b)
}

func (this *F) IMul(other F) *F {
	*this = this.Mul(other)
	return this
}

// binary gcd
func gcd(a, b int) int {
	x, y := a, b
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	if x == 0 || y == 0 {
		return x + y
	}
	n := bits.TrailingZeros(uint(x))
	m := bits.TrailingZeros(uint(y))
	x >>= uint(n)
	y >>= uint(m)
	for x != y {
		d := bits.TrailingZeros(uint(x - y))
		f := x > y
		var c int
		if f {
			c = x
		} else {
			c = y
		}
		if !f {
			y = x
		}
		x = (c - y) >> uint(d)
	}
	return x << uint(min(n, m))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
