package main

import "fmt"

func main() {

}

const INF int = 1e18

type sign int8

const (
	neg  sign = -1
	zero sign = 0
	pos  sign = 1
)

type Fraction struct {
	a, b  int  // a/b
	pos   sign // -1, 0, 1
	isInf bool
}

// Namespace
type FractionModule struct{}

func (*FractionModule) New(a, b int) Fraction {
	if b == 0 {
		panic(fmt.Sprintf("ZeroDivisionError: Fraction(%d, %d)", a, b))
	}

	if a == 0 {
		return Fraction{0, 1, zero, false}
	}

	isMinus := (a > 0) != (b > 0)
	if a == -INF || a == INF {
		if isMinus {
			return Fraction{INF, 1, neg, true}
		} else {
			return Fraction{INF, 1, pos, true}
		}
	}

	if isMinus {
		return Fraction{abs(a), abs(b), neg, false}
	}
	return Fraction{abs(a), abs(b), pos, false}
}

func (*FractionModule) Add(f1, f2 Fraction) Fraction {
	if f1.isInf || f2.isInf {
		if f1.pos != f2.pos {
			panic("nan")
		}
		return f1
	}
	sum := f1.a*f2.b + f2.a*f1.b
	mul := f1.b * f2.b
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// 有时也可以用这个简化分数大小比较 (不超过1e9时)
//
// 分母不为0的分数比较大小
//  a1/b1 < a2/b2
func less(a1, b1, a2, b2 int) bool {
	if a1 == INF || a1 == -INF || a2 == INF || a2 == -INF { // 有一个是+-INF
		return a1/b1 < a2/b2
	}
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	return diff^mul < 0
}
