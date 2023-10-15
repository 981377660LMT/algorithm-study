// 大力分块，对每个块维护bitset
// 修改直接做，查询直接合并bitset
// O(n^2/w)
// !TLE

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"math/bits"
	"os"
	"strconv"
)

func PointSetRangeType(nums []int, operations [][3]int) []int {
	block := UseBlock(nums, int(30*math.Sqrt(float64(len(nums)))+1))
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	n := len(nums)
	D := NewDictionary()
	newNums := make([]int, n)
	for i := range nums {
		newNums[i] = D.Id(nums[i])
	}
	for i := range operations {
		if operations[i][0] == 1 {
			operations[i][2] = D.Id(operations[i][2])
		}
	}

	blockType := make([]_Bitset, blockCount)
	blockCounter := make([][]int, blockCount)
	for bid := 0; bid < blockCount; bid++ {
		blockType[bid] = _NewBitset(D.Size())
		blockCounter[bid] = make([]int, D.Size())
		for j := blockStart[bid]; j < blockEnd[bid]; j++ {
			blockType[bid].Set(newNums[j])
			blockCounter[bid][newNums[j]]++
		}
	}

	res := []int{}
	bs := _NewBitset(D.Size())
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			bs.Clear()
			start, end := op[1], op[2]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
				onesCount := 0
				for i := start; i < end; i++ {
					if !bs.Has(newNums[i]) {
						bs.Set(newNums[i])
						onesCount++
					}
				}
				res = append(res, onesCount)
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					bs.Set(newNums[i])
				}
				for bid := bid1 + 1; bid < bid2; bid++ {
					bs.IOr(blockType[bid])
				}
				for i := blockStart[bid2]; i < end; i++ {
					bs.Set(newNums[i])
				}
				res = append(res, bs.OnesCount())
			}
		} else {
			pos, target := op[1], op[2]
			if newNums[pos] == target {
				continue
			}
			bid := belong[pos]
			blockCounter[bid][newNums[pos]]--
			if blockCounter[bid][newNums[pos]] == 0 {
				blockType[bid].Reset(newNums[pos])
			}
			blockCounter[bid][target]++
			if blockCounter[bid][target] == 1 {
				blockType[bid].Set(target)
			}
			newNums[pos] = target
		}
	}
	return res
}

// from https://atcoder.jp/users/ccppjsrb
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

	n, q := io.NextInt(), io.NextInt()
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = io.NextInt()
	}

	operations := [][3]int{}
	for i := 0; i < q; i++ {
		var op string
		op = io.Text()
		if op == "Q" {
			l, r := io.NextInt(), io.NextInt()
			l--
			operations = append(operations, [3]int{0, l, r})
		} else {
			p, col := io.NextInt(), io.NextInt()
			p--
			operations = append(operations, [3]int{1, p, col})
		}
	}

	res := PointSetRangeType(nums, operations)
	for _, v := range res {
		io.Println(v)
	}

}

type _Bitset []uint

func _NewBitset(n int) _Bitset { return make(_Bitset, n>>6+1) } // (n+64-1)>>6

func (b _Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b _Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b _Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b _Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0
func (b _Bitset) Clear() {
	for i := range b {
		b[i] = 0
	}
}
func (b _Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}

// 将 c 的元素合并进 b
func (b _Bitset) IOr(c _Bitset) {
	for i, v := range c {
		b[i] |= v
	}
}

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
