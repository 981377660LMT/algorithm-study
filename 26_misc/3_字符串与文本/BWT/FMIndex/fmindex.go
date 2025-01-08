// https://blog.csdn.net/stormlovetao/article/details/7048481
// https://github.com/shenwei356/bwt
// https://github.com/rossmerr/fm-index/tree/77e6c665a79e
// https://hc1023.github.io/2020/03/17/Short-Read-Alignment/
//
// read alignment
// 使用FM索引和BWT快速将reads与参考基因组对齐

package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

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

// Put puts element to stack
func (s *Stack) Put(i sMatch) {
	(*s) = append((*s), i)
}

// Pop pops element from the stack
func (s *Stack) Pop() sMatch {
	d := (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
	return d
}

// FMIndex is Burrows-Wheeler Index
type FMIndex struct {
	// EndSymbol
	EndSymbol byte

	// SuffixArray
	SuffixArray []int

	// Burrows-Wheeler Transform
	BWT []byte

	// First column of BWM
	F []byte

	// Alphabet in the BWT
	Alphabet []byte

	// Count of Letters in Alphabet.
	// CountOfLetters map[byte]int
	CountOfLetters []int // slice is faster han map

	// C[c] is a table that, for each character c in the alphabet,
	// contains the number of occurrences of lexically smaller characters
	// in the text.
	// C map[byte]int
	C []int // slice is faster han map

	// Occ(c, k) is the number of occurrences of character c in the
	// prefix L[1..k], k is 0-based.
	// Occ map[byte]*[]int32
	Occ []*[]int32 // slice is faster han map
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

	// fmi.Alphabet = byteutil.AlphabetFromCountOfByte(fmi.CountOfLetters)
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

	locationsMap := make(map[int]struct{})

	if mismatches == 0 {
		// letters := byteutil.Alphabet(query)
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
			// if _, ok := fmi.CountOfLetters[letter]; !ok {
			if fmi.CountOfLetters[letter] == 0 {
				return locations, nil
			}
		}
	}

	n := len(fmi.BWT)
	var matches Stack

	// start and end are 0-based
	matches.Put(sMatch{query: query, start: 0, end: n - 1, mismatches: mismatches})
	// fmt.Printf("====%s====\n", query)
	// fmt.Println(fmi)
	var match sMatch
	var last, c byte
	var start, end int
	var m int

	var letters []byte

	// var ok bool
	for !matches.Empty() {
		match = matches.Pop()
		query = match.query[0 : len(match.query)-1]
		last = match.query[len(match.query)-1]
		if match.mismatches == 0 {
			letters = []byte{last}
		} else {
			letters = fmi.Alphabet
		}

		// fmt.Println("\n--------------------------------------------")
		// fmt.Printf("%s, %s, %c\n", match.query, query, last)
		// fmt.Printf("query: %s, last: %c\n", query, last)
		for _, c = range letters {
			// if _, ok = fmi.CountOfLetters[c]; !ok { //  letter not in alphabet
			if fmi.CountOfLetters[c] == 0 {
				continue
			}

			// fmt.Printf("letter: %c, start: %d, end: %d, mismatches: %d\n", c, match.start, match.end, match.mismatches)
			if match.start == 0 {
				start = fmi.C[c] + 0
			} else {
				start = fmi.C[c] + int((*fmi.Occ[c])[match.start-1])
			}
			end = fmi.C[c] + int((*fmi.Occ[c])[match.end]-1)
			// fmt.Printf("    s: %d, e: %d\n", start, end)

			if start > end {
				continue
			}

			if len(query) == 0 {
				for _, i := range fmi.SuffixArray[start : end+1] {
					// fmt.Printf("    >>> found: %d\n", i)
					locationsMap[i] = struct{}{}
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

				// fmt.Printf("    >>> candidate: query: %s, start: %d, end: %d, m: %d\n", query, start, end, m)
				matches.Put(sMatch{query: query, start: start, end: end, mismatches: m})
			}
		}
	}

	i := 0
	locations = make([]int, len(locationsMap))
	for loc := range locationsMap {
		locations[i] = loc
		i++
	}
	sort.Ints(locations)

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

	matches.Put(sMatch{query: query, start: 0, end: n - 1, mismatches: mismatches})

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

				matches.Put(sMatch{query: query, start: start, end: end, mismatches: m})
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
