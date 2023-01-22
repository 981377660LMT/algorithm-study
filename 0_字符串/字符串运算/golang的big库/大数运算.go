// Golang实现BigInt大数计算
// https://cloud.tencent.com/developer/article/2029133
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta

package main

import (
	"math/big"
)

type BigInt struct {
	*big.Int
}

func NewBigInt(s string) BigInt {
	var e big.Int
	newValue, ok := e.SetString(s, 10)
	if !ok {
		return BigInt{big.NewInt(0)}
	}
	return BigInt{newValue}
}

// i开头表示in-place操作，即修改自身
func (a BigInt) iAdd(b BigInt) BigInt { a.Add(a.Int, b.Int); return a }
func (a BigInt) iSub(b BigInt) BigInt { a.Sub(a.Int, b.Int); return a }
func (a BigInt) iMul(b BigInt) BigInt { a.Mul(a.Int, b.Int); return a }
func (a BigInt) iDiv(b BigInt) BigInt { a.Div(a.Int, b.Int); return a }
func (a BigInt) iMod(b BigInt) BigInt { a.Mod(a.Int, b.Int); return a }
func (a BigInt) iNeg(b BigInt) BigInt { a.Neg(a.Int); return a }

func (a BigInt) add(b BigInt) BigInt { return BigInt{new(big.Int).Add(a.Int, b.Int)} }
func (a BigInt) sub(b BigInt) BigInt { return BigInt{new(big.Int).Sub(a.Int, b.Int)} }
func (a BigInt) mul(b BigInt) BigInt { return BigInt{new(big.Int).Mul(a.Int, b.Int)} }
func (a BigInt) div(b BigInt) BigInt { return BigInt{new(big.Int).Div(a.Int, b.Int)} }
func (a BigInt) mod(b BigInt) BigInt { return BigInt{new(big.Int).Mod(a.Int, b.Int)} }
func (a BigInt) neg(b BigInt) BigInt { return BigInt{new(big.Int).Neg(a.Int)} }
func (a BigInt) cmp(b BigInt) int    { return a.Cmp(b.Int) }
