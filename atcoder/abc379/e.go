package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const MOD = 998244353

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp *os.File, wfp *os.File) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}

func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}

func (io *Iost) Atoi(s string) int {
	x, _ := strconv.Atoi(s)
	return x
}

func (io *Iost) Println(x ...interface{}) {
	fmt.Fprintln(io.Writer, x...)
}

func generateValidRows(options [][]byte, W int) []string {
	valid := []string{}
	current := make([]byte, W)
	var backtrack func(int)
	backtrack = func(pos int) {
		if pos == W {
			valid = append(valid, string(current))
			return
		}
		for _, c := range options[pos] {
			if pos > 0 && c == current[pos-1] {
				continue
			}
			current[pos] = c
			backtrack(pos + 1)
		}
	}
	backtrack(0)
	return valid
}

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer io.Writer.Flush()

	H := io.Atoi(io.Text())
	W := io.Atoi(io.Text())
	S := make([]string, H)
	for i := 0; i < H; i++ {
		S[i] = io.Text()
	}

	// 生成所有可能的合法行
	validRows := make([][]string, H)
	for i := 0; i < H; i++ {
		options := make([][]byte, W)
		for j := 0; j < W; j++ {
			if S[i][j] == '?' {
				options[j] = []byte{'1', '2', '3'}
			} else {
				options[j] = []byte{S[i][j]}
			}
		}
		validRows[i] = generateValidRows(options, W)
	}

	// 映射行到ID
	rowToID := make([]map[string]int, H)
	for i := 0; i < H; i++ {
		rowToID[i] = make(map[string]int)
		for idx, row := range validRows[i] {
			if _, exists := rowToID[i][row]; !exists {
				rowToID[i][row] = idx
			}
		}
	}

	// 动态规划
	prev := make([]int, len(validRows[0]))
	for i := 0; i < len(validRows[0]); i++ {
		prev[i] = 1
	}

	for i := 1; i < H; i++ {
		current := make([]int, len(validRows[i]))
		for j, currRow := range validRows[i] {
			for k, prevRow := range validRows[i-1] {
				conflict := false
				for c := 0; c < W; c++ {
					if currRow[c] == prevRow[c] {
						conflict = true
						break
					}
				}
				if !conflict {
					current[j] = (current[j] + prev[k]) % MOD
				}
			}
		}
		prev = current
	}

	result := 0
	for _, val := range prev {
		result = (result + val) % MOD
	}
	io.Println(result)
}
