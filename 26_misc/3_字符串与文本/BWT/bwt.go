// Burrows–Wheeler Transform (BWT) 是一种在数据压缩领域广受关注的可逆变换。
// 将输入序列（通常为字符串）的字符重新排列，使得相同或相似的字符在变换后更趋于聚集，
// 从而有利于后续使用游程编码（Run-Length Encoding, RLE）或其它编码算法（例如 MTF、Huffman 等）来进一步压缩。
//
// !BWT 的重要性质：
// last-to-first mapping: https://curiouscoding.nl/posts/bwt/
// !以字母x开头的旋转顺序与以字母x结尾的旋转顺序相同
// 证明：将字符串旋转一次，将第一个字符移动到最后，因此，最后一列中 G 的顺序与第一列中的顺序相同.
package main

import (
	"errors"
	"fmt"
	"index/suffixarray"
	"reflect"
)

func main() {
	s := "banana"
	bwt, err := Transform([]byte(s), '$')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bwt)) //annb$aa

	origin := InverseTransform(bwt, '$')
	fmt.Println(string(origin)) //banana
}

// CheckEndSymbol is a global variable for checking end symbol before Burrows–Wheeler transform
var CheckEndSymbol = true

// ErrEndSymbolExisted means you should choose another EndSymbol
var ErrEndSymbolExisted = errors.New("bwt: end-symbol existed in string")

// ErrEmptySeq means a empty sequence is given
var ErrEmptySeq = errors.New("bwt: empty sequence")

// ErrInvalidSuffixArray means length of sa is not equal to 1+len(s)
var ErrInvalidSuffixArray = errors.New("bwt: invalid suffix array")

// Transform returns Burrows–Wheeler transform of a byte slice.
// See https://en.wikipedia.org/wiki/Burrows%E2%80%93Wheeler_transform
func Transform(s []byte, es byte) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptySeq
	}
	if CheckEndSymbol {
		for _, c := range s {
			if c == es {
				return nil, ErrEndSymbolExisted
			}
		}
	}
	sa := suffixArray(s)
	bwt, err := fromSuffixArray(s, sa, es)
	return bwt, err
}

// LF(Last to First)-mapping O(n) 计数排序还原.
// 将BWT矩阵最后一列的某个字符x映射到第一列中，看这个x在第一列中哪一行.
// !BWT矩阵性质：最后一列第i个出现的字符x，与第一列中第i个出现的x是同一个.
func InverseTransform(t []byte, es byte) []byte {
	n := len(t)
	if n == 0 {
		return nil
	}

	var freq [256]int
	for _, c := range t {
		freq[c]++
	}

	// 构造前缀和数组 C，C[c] = 字符 c 在排序后 F 中的起始下标
	// 其中 F = sort(t)，即 BWT 矩阵的第一列
	var C [256]int
	sum := 0
	for i := 0; i < 256; i++ {
		C[i] = sum
		sum += freq[i]
	}

	// 构造 next 数组
	// next[i] = C[t[i]] + rank   (其中 rank 为字符 t[i] 在 t[:i] 中出现的次数)
	next := make([]int, n)
	// occ[c] 记录「当前已经处理」多少次字符 c
	var occ [256]int
	for i := 0; i < n; i++ {
		c := t[i]
		next[i] = C[c] + occ[c]
		occ[c]++
	}

	// 找到含有终止符 es 的行下标 rowWithEndSymbol
	rowWithEndSymbol := 0
	for i := 0; i < n; i++ {
		if t[i] == es {
			rowWithEndSymbol = i
			break
		}
	}

	// 逆推：从 rowWithEndSymbol 开始，根据 next 数组不断往回追溯
	// 注意：由于 BWT 包含终止符 es，原串实际长度为 n-1
	// 逆推得到的是倒序结果，需要最后做一次翻转
	res := make([]byte, 0, n-1)
	curRow := rowWithEndSymbol
	for i := 0; i < n-1; i++ {
		curRow = next[curRow]
		res = append(res, t[curRow])
	}

	reverseBytes(res)
	return res
}

func reverseBytes(a []byte) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func suffixArray(s []byte) []int {
	_sa := suffixarray.New(s)
	tmp := reflect.ValueOf(_sa).Elem().FieldByName("sa").FieldByIndex([]int{0})
	var sa []int = make([]int, len(s)+1)
	sa[0] = len(s)
	for i := 0; i < len(s); i++ {
		sa[i+1] = int(tmp.Index(i).Int())
	}
	return sa
}

// fromSuffixArray compute BWT from sa
func fromSuffixArray(s []byte, sa []int, es byte) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptySeq
	}
	if len(s)+1 != len(sa) || sa[0] != len(s) {
		return nil, ErrInvalidSuffixArray
	}
	bwt := make([]byte, len(sa))
	bwt[0] = s[len(s)-1]
	for i := 1; i < len(sa); i++ {
		if sa[i] == 0 {
			bwt[i] = es
		} else {
			bwt[i] = s[sa[i]-1]
		}
	}
	return bwt, nil
}
