// https://blog.csdn.net/stormlovetao/article/details/7048481
// https://github.com/shenwei356/bwt
// https://github.com/rossmerr/fm-index/tree/77e6c665a79e
// https://hc1023.github.io/2020/03/17/Short-Read-Alignment/
//
// !FMIndex 的朴素实现.
//
// func (fmi *FMIndex) Transform(s []byte) ([]byte, error)
// func (fmi *FMIndex) Locate(query []byte, mismatches int) ([]int, error)
// func (fmi *FMIndex) Match(query []byte, mismatches int) (bool, error)
// func (fmi *FMIndex) Last2First(i int) int
// func (fmi *FMIndex) String() string
//
// 搜索原理：
// 从后向前找。
// 使用 LF 映射通过 C 表和 Occ 表将当前的搜索范围 [match.start, match.end] 映射到 F 列中的新范围 [start, end]，
// 以便在下一步迭代中继续缩小搜索范围

package main

import (
	"bytes"
	"errors"
	"fmt"
	"index/suffixarray"
	"reflect"
	"slices"
	"strings"
)

// TLE.
// https://leetcode.cn/problems/find-the-occurrence-of-first-almost-equal-substring/description/
func minStartingIndex(s string, pattern string) int {
	fmi := NewFMIndex()
	fmi.Transform([]byte(s))
	locations, err := fmi.Locate([]byte(pattern), 1)
	if err != nil {
		return -1
	}
	if len(locations) == 0 {
		return -1
	}
	return locations[0]
}

func demo() {
	fmi := NewFMIndex()
	bwt, err := fmi.Transform([]byte("banana"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bwt)) //annb$aa

	locations, err := fmi.Locate([]byte("ana"), 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(locations) //[1 3]

	ok, err := fmi.Match([]byte("ana"), 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ok) //true

	fmt.Println(fmi)
}

// ErrEmptySeq means a empty sequence is given
var ErrEmptySeq = errors.New("bwt: empty sequence")

// ErrInvalidSuffixArray means length of sa is not equal to 1+len(s)
var ErrInvalidSuffixArray = errors.New("bwt: invalid suffix array")

// FMIndex is Burrows-Wheeler Index.
// !两个大数组：SuffixArray 和 Occ 在工程场景中需要优化存储.
type FMIndex struct {
	EndSymbol byte

	// !工业实现中不开这个大数组，而是压缩的后缀数组(或者使用checkpoint方法减少内存占用).
	SuffixArray []int

	// BWT 重新排列了 s 的字符，以提升压缩效率并支持高效搜索。
	BWT []byte

	// BWM（布尔洛斯-维尔纳矩阵）的第一列，通常在文献中称为 F。
	F []byte

	Alphabet []byte

	// 用于存储每个字符在 BWT 中的出现次数。使用切片而非映射以提升性能，使用 ASCII 值（0-127）作为索引
	CountOfLetters []int

	// 前缀和数组，其中 C[c] 表示文本中所有字典顺序小于字符 c 的字符总数
	C []int

	// !Occ[c][k] 表示在 BWT 的前缀 L[1..k] 中字符 c 的出现次数
	// !工业实现中不开这个大的数组，只在某些位置设置 checkpoint.
	// 在bowtie里每个448行设置一个checkpoint(为什么不临接表二分?)
	Occ []*[]int32
}

// NewFMIndex is constructor of FMIndex
func NewFMIndex() *FMIndex {
	fmi := new(FMIndex)
	fmi.EndSymbol = byte(0)
	return fmi
}

// Transform return Burrows-Wheeler-Transform of s
func (fmi *FMIndex) Transform(s []byte) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptySeq
	}
	var err error

	sa := suffixArray(s)
	fmi.SuffixArray = sa

	fmi.BWT, err = fromSuffixArray(s, fmi.SuffixArray, fmi.EndSymbol)
	if err != nil {
		return nil, err
	}

	F := make([]byte, len(s)+1)
	F[0] = fmi.EndSymbol
	for i := 1; i <= len(s); i++ {
		F[i] = s[sa[i]]
	}
	fmi.F = F

	count := make([]int, 128)
	for _, b := range fmi.BWT {
		count[b]++
	}
	count[fmi.EndSymbol] = 0
	fmi.CountOfLetters = count

	alphabet := make([]byte, 0, 128)
	for b, c := range count {
		if c > 0 {
			alphabet = append(alphabet, byte(b))
		}
	}
	fmi.Alphabet = alphabet

	fmi.C = computeC(fmi.F)

	fmi.Occ = computeOccurrence(fmi.BWT, fmi.Alphabet)

	return fmi.BWT, nil
}

// Last2First mapping
// 该映射允许从 BWT 矩阵的最后一列 (L) O(1) 时间内找到对应的第一列 (F) 中的字符。
func (fmi *FMIndex) Last2First(i int) int {
	c := fmi.BWT[i]
	return fmi.C[c] + int((*fmi.Occ[c])[i])
}

func (fmi *FMIndex) nextLetterInAlphabet(c byte) byte {
	var nextLetter byte
	for i, letter := range fmi.Alphabet {
		if letter == c {
			if i < len(fmi.Alphabet)-1 {
				nextLetter = fmi.Alphabet[i+1]
			} else {
				nextLetter = fmi.Alphabet[i]
			}
			break
		}
	}
	return nextLetter
}

// Locate locates the pattern
func (fmi *FMIndex) Locate(query []byte, mismatches int) ([]int, error) {
	if len(query) == 0 {
		return []int{}, nil
	}
	var locations []int

	locationsSet := make(map[int]struct{})

	if mismatches == 0 {
		count := make([]int, 128)
		for _, b := range query {
			if count[b] == 0 {
				count[b]++
			}
		}
		letters := make([]byte, 0, 128)
		for b, c := range count {
			if c > 0 {
				letters = append(letters, byte(b))
			}
		}

		for _, letter := range letters { // query having letter not in alphabet
			if fmi.CountOfLetters[letter] == 0 {
				return locations, nil
			}
		}
	}

	n := len(fmi.BWT)
	var matches Stack // 用于管理中间搜索状态的栈

	matches.Push(sMatch{query: query, start: 0, end: n - 1, mismatches: mismatches})
	var match sMatch
	var last, c byte
	var start, end int
	var m int

	var letters []byte

	for !matches.Empty() {
		match = matches.Pop()
		query = match.query[0 : len(match.query)-1]
		last = match.query[len(match.query)-1]

		// 如果不允许误差（mismatches == 0），则仅考虑精确匹配的字符。
		// 否则，考虑所有字母表中的字符以允许近似匹配。
		if match.mismatches == 0 {
			letters = []byte{last}
		} else {
			letters = fmi.Alphabet
		}

		for _, c = range letters {
			// letter not in alphabet
			if fmi.CountOfLetters[c] == 0 {
				continue
			}

			if match.start == 0 {
				start = fmi.C[c] + 0
			} else {
				start = fmi.C[c] + int((*fmi.Occ[c])[match.start-1])
			}
			end = fmi.C[c] + int((*fmi.Occ[c])[match.end]-1)

			if start > end {
				continue
			}

			if len(query) == 0 {
				for _, i := range fmi.SuffixArray[start : end+1] {
					locationsSet[i] = struct{}{}
				}
			} else {
				m = match.mismatches
				if c != last {
					if match.mismatches > 1 {
						m = match.mismatches - 1
					} else {
						m = 0
					}
				}

				matches.Push(sMatch{query: query, start: start, end: end, mismatches: m})
			}
		}
	}

	i := 0
	locations = make([]int, len(locationsSet))
	for loc := range locationsSet {
		locations[i] = loc
		i++
	}
	slices.Sort(locations)

	return locations, nil
}

// Match is a simple version of Locate, which returns immediately for a match.
func (fmi *FMIndex) Match(query []byte, mismatches int) (bool, error) {
	if len(query) == 0 {
		return false, nil
	}
	if mismatches == 0 {
		count := make([]int, 128)
		for _, b := range query {
			if count[b] == 0 {
				count[b]++
			}
		}
		letters := make([]byte, 0, 128)
		for b, c := range count {
			if c > 0 {
				letters = append(letters, byte(b))
			}
		}

		for _, letter := range letters { // query having letter not in alphabet
			if fmi.CountOfLetters[letter] == 0 {
				return false, nil
			}
		}
	}

	n := len(fmi.BWT)
	var matches Stack

	matches.Push(sMatch{query: query, start: 0, end: n - 1, mismatches: mismatches})

	var match sMatch
	var last, c byte
	var start, end int
	var m int

	var letters []byte

	for !matches.Empty() {
		match = matches.Pop()
		query = match.query[0 : len(match.query)-1]
		last = match.query[len(match.query)-1]
		if match.mismatches == 0 {
			letters = []byte{last}
		} else {
			letters = fmi.Alphabet
		}

		for _, c = range letters {
			if fmi.CountOfLetters[c] == 0 {
				continue
			}

			if match.start == 0 {
				start = fmi.C[c] + 0
			} else {
				start = fmi.C[c] + int((*fmi.Occ[c])[match.start-1])
			}
			end = fmi.C[c] + int((*fmi.Occ[c])[match.end]-1)

			if start > end {
				continue
			}

			if len(query) == 0 {
				return true, nil
			} else {
				m = match.mismatches
				if c != last {
					if match.mismatches > 1 {
						m = match.mismatches - 1
					} else {
						m = 0
					}
				}

				matches.Push(sMatch{query: query, start: start, end: end, mismatches: m})
			}
		}
	}

	return false, nil
}

func (fmi *FMIndex) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("EndSymbol: %c\n", fmi.EndSymbol))
	buffer.WriteString(fmt.Sprintf("BWT: %s\n", string(fmi.BWT)))
	buffer.WriteString(fmt.Sprintf("Alphabet: %s\n", string(fmi.Alphabet)))
	buffer.WriteString("F:\n")
	buffer.WriteString(string(fmi.F) + "\n")
	buffer.WriteString("C:\n")
	for _, letter := range fmi.Alphabet {
		buffer.WriteString(fmt.Sprintf("  %c: %d\n", letter, fmi.C[letter]))
	}
	buffer.WriteString("Occ:\n")
	buffer.WriteString(fmt.Sprintf("  BWT[%s]\n", strings.Join(strings.Split(string(fmi.BWT), ""), " ")))
	for _, letter := range fmi.Alphabet {
		buffer.WriteString(fmt.Sprintf("  %c: %v\n", letter, fmi.Occ[letter]))
	}

	buffer.WriteString("SA:\n")
	buffer.WriteString(fmt.Sprintf("  %d\n", fmi.SuffixArray))

	return buffer.String()
}

// ComputeC computes C.
// C[c] is a table that, for each character c in the alphabet,
// contains the number of occurrences of lexically smaller characters
// in the text.
//
//	func ComputeC(L []byte, alphabet []byte) map[byte]int {
//		if alphabet == nil {
//			alphabet = byteutil.Alphabet(L)
//		}
//		C := make(map[byte]int, len(alphabet))
//		count := 0
//		for _, c := range L {
//			if _, ok := C[c]; !ok {
//				C[c] = count
//			}
//			count++
//		}
//		return C
//	}
func computeC(L []byte) []int {
	C := make([]int, 128)
	count := 0
	for _, c := range L {
		if C[c] == 0 {
			C[c] = count
		}
		count++
	}
	return C
}

// ComputeOccurrence returns occurrence information.
// Occ(c, k) is the number of occurrences of character c in the prefix L[1..k]
//
//	func ComputeOccurrence(bwt []byte, letters []byte) map[byte]*[]int32 {
//		if letters == nil {
//			letters = byteutil.Alphabet(bwt)
//		}
//		occ := make(map[byte]*[]int32, len(letters)-1)
//		for _, letter := range letters {
//			t := make([]int32, 1, len(bwt))
//			t[0] = 0
//			occ[letter] = &t
//		}
//		t := make([]int32, 1, len(bwt))
//		t[0] = 1
//		occ[bwt[0]] = &t
//		var letter, k byte
//		var v *[]int32
//		for _, letter = range bwt[1:] {
//			for k, v = range occ {
//				if k == letter {
//					*v = append(*v, (*v)[len(*v)-1]+1)
//				} else {
//					*v = append(*v, (*v)[len(*v)-1])
//				}
//			}
//		}
//		return occ
//	}
func computeOccurrence(bwt []byte, letters []byte) []*[]int32 {
	if letters == nil {
		count := make([]int, 128)
		for _, b := range bwt {
			if count[b] == 0 {
				count[b]++
			}
		}

		letters = make([]byte, 0, 128)
		for b, c := range count {
			if c > 0 {
				letters = append(letters, byte(b))
			}
		}
	}

	occ := make([]*[]int32, 128)
	for _, letter := range letters {
		t := make([]int32, 1, len(bwt))
		t[0] = 0
		occ[letter] = &t
	}
	t := make([]int32, 1, len(bwt))
	t[0] = 1
	occ[bwt[0]] = &t
	var letter byte
	var k, letterInt int
	var v *[]int32
	for _, letter = range bwt[1:] {
		letterInt = int(letter)
		for k, v = range occ {
			if v == nil {
				continue
			}

			if k == letterInt {
				*v = append(*v, (*v)[len(*v)-1]+1)
			} else {
				*v = append(*v, (*v)[len(*v)-1])
			}
		}
	}
	return occ
}

type sMatch struct {
	query      []byte
	start, end int
	mismatches int
}

// Stack struct
type Stack []sMatch

// Empty tell if it is empty
func (s Stack) Empty() bool {
	return len(s) == 0
}

// Peek return the last element
func (s Stack) Peek() sMatch {
	return s[len(s)-1]
}

// Push puts element to stack
func (s *Stack) Push(i sMatch) {
	(*s) = append((*s), i)
}

// Pop pops element from the stack
func (s *Stack) Pop() sMatch {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
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
