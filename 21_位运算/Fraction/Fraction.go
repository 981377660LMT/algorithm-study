// Rational/Fraction
// https://nyaannyaan.github.io/library/math/rational.hpp

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	a, b := Fraction(4, 3), Fraction(2, 3)
	assert := func(value interface{}, expected interface{}) {
		if value != expected {
			panic(fmt.Sprintf("assert failed: %v != %v", value, expected))
		}
	}

	assert(a.Add(b), Fraction(2, 1))
	assert(a.Sub(b), Fraction(2, 3))
	assert(b.Sub(a), Fraction(-2, 3))
	assert(a.Mul(b), Fraction(8, 9))
	assert(a.Div(b), Fraction(2, 1))
	assert(a.Inv(), Fraction(3, 4))
	assert(a.Pow(3), Fraction(64, 27))
	assert(a.Gt(b), true)
	assert(a.Ge(b), true)
	assert(a.Lt(b), false)
	assert(a.Le(b), false)
}

// !中途所有数不超过 2e9

type F struct {
	a int // 分子, numerator
	b int // 分母, denominator，在内部恒为正数
}

func Fraction(a, b int) F {
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
	return Fraction(this.a*other.b+this.b*other.a, this.b*other.b)
}

func (this *F) IAdd(other F) *F {
	*this = this.Add(other)
	return this
}

func (this *F) Sub(other F) F {
	return Fraction(this.a*other.b-this.b*other.a, this.b*other.b)
}

func (this *F) ISub(other F) *F {
	*this = this.Sub(other)
	return this
}

func (this *F) Mul(other F) F {
	return Fraction(this.a*other.a, this.b*other.b)
}

func (this *F) IMul(other F) *F {
	*this = this.Mul(other)
	return this
}

func (this *F) Div(other F) F {
	return Fraction(this.a*other.b, this.b*other.a)
}

func (this *F) IDiv(other F) *F {
	*this = this.Div(other)
	return this
}

// 负数
func (this *F) Neg() F { return F{-this.a, this.b} }

// 倒数
func (this *F) Inv() F {
	if this.a == 0 {
		panic("Not invertible")
	}
	res := F{this.b, this.a}
	if res.b < 0 {
		res.a = -res.a
		res.b = -res.b
	}
	return res
}

func (this *F) Pow(exp int) F {
	res := F{1, 1}
	base := *this
	for exp > 0 {
		if exp&1 == 1 {
			res.IMul(base)
		}
		base.IMul(base)
		exp >>= 1
	}
	return res
}

// this == other
func (this *F) Eq(other F) bool { return this.a == other.a && this.b == other.b }

// this != other
func (this *F) Neq(other F) bool { return this.a != other.a || this.b != other.b }

// this < other
func (this *F) Lt(other F) bool { return this.a*other.b < other.a*this.b }

// this <= other
func (this *F) Le(other F) bool { return this.a*other.b <= other.a*this.b }

// this > other
func (this *F) Gt(other F) bool { return this.a*other.b > other.a*this.b }

// this >= other
func (this *F) Ge(other F) bool { return this.a*other.b >= other.a*this.b }

func (this *F) String() string {
	sb := []string{}
	sb = append(sb, fmt.Sprintf("%d", this.a))
	if this.a != 0 && this.b != 1 {
		sb = append(sb, fmt.Sprintf("/%d", this.b))
	}
	return strings.Join(sb, "")
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
	x >>= n
	y >>= m
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
		x = (c - y) >> d
	}
	return x << min(n, m)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//
//
// 一般的分数比较大小
// !分母不为0且所有数绝对值不超过1e9

// a1/b1 == a2/b2
func eq(a1, b1, a2, b2 int) bool {
	return a1*b2 == a2*b1
}

// a1/b1 < a2/b2
func lt(a1, b1, a2, b2 int) bool {
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	if diff == 0 {
		return false
	}
	return (diff > 0) != (mul > 0)
}

// a1/b1 <= a2/b2
func le(a1, b1, a2, b2 int) bool {
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	if diff == 0 {
		return true
	}
	return (diff > 0) != (mul > 0)
}
