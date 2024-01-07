package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
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
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	sqrt := int(math.Sqrt(float64(len(nums))) + 1)
	B := UseBlock(len(nums), sqrt)
	belong, blockStart, blockEnd, blockCount := B.belong, B.blockStart, B.blockEnd, B.blockCount
	blockSum := make([]int, blockCount)
	for i := range blockSum {
		for j := blockStart[i]; j < blockEnd[i]; j++ {
			blockSum[i] = (blockSum[i] + nums[j]) % MOD
		}
	}

	pre, suf := make([][]int, sqrt+1), make([][]int, sqrt+1)
	for i := range pre {
		pre[i] = make([]int, sqrt+1)
		suf[i] = make([]int, sqrt+1)
	}

	_sum := func(start, end int) int {
		bid1, bid2 := belong[start], belong[end-1]
		res := 0
		if bid1 == bid2 {
			for i := start; i < end; i++ {
				res = (res + nums[i]) % MOD
			}
		} else {
			for i := start; i < blockEnd[bid1]; i++ {
				res = (res + nums[i]) % MOD
			}
			for bid := bid1 + 1; bid < bid2; bid++ {
				res = (res + blockSum[bid]) % MOD
			}
			for i := blockStart[bid2]; i < end; i++ {
				res = (res + nums[i]) % MOD
			}
		}
		return res
	}

	update := func(start, step, delta int) {
		if step >= sqrt {
			for i := start; i < len(nums); i += step {
				bid := belong[i]
				nums[i] = (nums[i] + delta) % MOD
				blockSum[bid] = (blockSum[bid] + delta) % MOD
			}
		} else {
			curPre, curSuf := pre[step], suf[step]
			for i := start; i+1 < len(curPre); i++ {
				curPre[i+1] = (curPre[i+1] + delta) % MOD
			}
			for i := 0; i <= start && i < len(curSuf); i++ {
				curSuf[i] = (curSuf[i] + delta) % MOD
			}
		}
	}

	query := func(start, end int) int {
		res := _sum(start, end)
		// 加上每个step对应的和.
		for step := 1; step < sqrt; step++ {
			id1, id2 := start/step, (end-1)/step
			pos1, pos2 := start%step, (end-1)%step
			curPre, curSuf := pre[step], suf[step]
			if id1 == id2 {
				res = (res + curPre[pos2+1] - curPre[pos1] + MOD) % MOD
			} else {
				res = (res + curSuf[pos1]) % MOD
				res = (res + (id2-id1-1)*curPre[step]) % MOD
				res = (res + curPre[pos2+1]) % MOD
			}
		}

		return res % MOD
	}

	// res := []int{}
	// for _, op := range operations {
	// 	kind := op[0]
	// 	if kind == 0 {
	// 		start, step, delta := op[1], op[2], op[3]
	// 		update(start, step, delta)
	// 	} else {
	// 		start, end := op[1], op[2]
	// 		res = append(res, query(start, end))
	// 	}
	// }

	// return res
	update(0, n, 1)
	for i := 0; i < n; i++ {
		cur := query(i, i+1)
		update(i+1, nums[i], cur)
	}
	io.Println(query(n-1, n) % MOD)

}

const MOD int = 998244353

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func bisectLeft(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
func bisectRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
