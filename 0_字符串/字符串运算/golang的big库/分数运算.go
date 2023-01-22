// Golang实现Decimal分数计算
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta

package main

import (
	"fmt"
	"math/big"
)

func main() {
	num := NewDecimal(3)
	deno := NewDecimal(4)
	res := num.div(deno)
	fmt.Println(res)
}

type Decimal struct {
	*big.Rat
}

func NewDecimal(value int) Decimal {
	var e big.Rat
	return Decimal{e.SetInt64(int64(value))}
}

// i开头表示in-place操作，即修改自身
func (a Decimal) iAdd(b Decimal) Decimal { a.Add(a.Rat, b.Rat); return a }
func (a Decimal) iSub(b Decimal) Decimal { a.Sub(a.Rat, b.Rat); return a }
func (a Decimal) iMul(b Decimal) Decimal { a.Mul(a.Rat, b.Rat); return a }
func (a Decimal) iDiv(b Decimal) Decimal { a.Quo(a.Rat, b.Rat); return a }
func (a Decimal) iNeg(b Decimal) Decimal { a.Neg(a.Rat); return a }

func (a Decimal) add(b Decimal) Decimal { return Decimal{new(big.Rat).Add(a.Rat, b.Rat)} }
func (a Decimal) sub(b Decimal) Decimal { return Decimal{new(big.Rat).Sub(a.Rat, b.Rat)} }
func (a Decimal) mul(b Decimal) Decimal { return Decimal{new(big.Rat).Mul(a.Rat, b.Rat)} }
func (a Decimal) div(b Decimal) Decimal { return Decimal{new(big.Rat).Quo(a.Rat, b.Rat)} }
func (a Decimal) neg(b Decimal) Decimal { return Decimal{new(big.Rat).Neg(a.Rat)} }
func (a Decimal) cmp(b Decimal) int     { return a.Cmp(b.Rat) }
