// P3792 由乃与大母神原型和偶像崇拜-区间异或哈希

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"time"

	stdio "io"
	"sort"
	"strconv"
)

// 1 start end : 判断区间[start, end)是否可以重排为值域上连续的一段(区间值域连续)
// 2 pos val : 将 nums[pos] 修改为 val
// n,q<=5e5, nums[i]<=1e9
//
// 使用异或哈希来判断区间的集合相等,即把每个数映射到一个随机数上，然后维护异或和来判断是否是连续段。
//  1. 离散化,为了使原来不相邻的数离散化后也不相邻，需要将每个数+1一起离散化.
//  2. 将离散化后的长度记为n,然后将1-n每个数映射到一个随机数上，并预处理出前缀异或和preXor[i].
//  3. 用树状数组维护原数组的和、映射到的随机数的异或和；
//     查询时根据长度与区间和可以知道一个连续段包含哪些数，然后与实际的异或和比较即可。
func 由乃与大母神原型和偶像崇拜(nums []int, operations [][3]int) []bool {
	n := len(nums)
	res := []bool{}

	allNums := make([]int, 2*n)
	copy(allNums, nums)
	for i, v := range nums {
		allNums[i+n] = v + 1
	}
	for _, op := range operations {
		if op[0] == 1 {
			allNums = append(allNums, op[2], op[2]+1)
		}
	}
	rank, count := Compress(allNums, 1) // 离散化到 1-count

	xor := make([]uint64, count+1)
	preXor := make([]uint64, count+1)
	random := NewRandom()
	for i := 1; i <= count; i++ {
		rand := random.Rng61()
		xor[i] = rand
		preXor[i] = preXor[i-1] ^ rand
	}

	newNums := make([]int, n)
	newNumsXor := make([]uint64, n)
	for i, v := range nums {
		r := rank(v)
		newNums[i] = r
		newNumsXor[i] = xor[r]
	}
	bitSum := NewBitArrayFrom(newNums)
	bitXor := NewBitArrayXorFrom(newNumsXor)

	// 值域[min,max]中的数的异或和.
	// 1<=min<=max<=count
	queryXorOfInterval := func(min, max int) uint64 {
		return preXor[max] ^ preXor[min-1]
	}

	// 原数组区间[start,end)离散化后的数的异或和.
	queryXorOfRange := func(start, end int) uint64 {
		return bitXor.QueryRange(start, end)
	}

	// 从区间和与区间长度还原区间.
	recoverIntervalFrom := func(sum int, len int) (first, last int, ok bool) {
		// 首项: x 末项 x+(len-1) 项数 len
		// (2*x+len-1)*len=2*sum
		if 2*sum%len != 0 {
			return 0, 0, false
		}
		div := 2 * sum / len
		if (div+1-len)&1 == 1 {
			return 0, 0, false
		}
		x := (div + 1 - len) / 2
		if x <= 0 || x+len-1 > count {
			return 0, 0, false
		}
		return x, x + len - 1, true
	}

	for _, op := range operations {
		kind := op[0]
		if kind == 1 {
			pos, val := op[1], op[2]
			val = rank(val)
			if newNums[pos] == val {
				continue
			}
			bitSum.Add(pos, val-newNums[pos])
			bitXor.Add(pos, xor[val]^xor[newNums[pos]])
			newNums[pos] = val
		} else {
			start, end := op[1], op[2]
			sum := bitSum.QueryRange(start, end)
			first, last, ok := recoverIntervalFrom(sum, end-start)
			if !ok {
				res = append(res, false)
				continue
			}

			// 判断实际的区间异或和与值域内异或和是否相等
			xor1 := queryXorOfInterval(first, last)
			xor2 := queryXorOfRange(start, end)
			res = append(res, xor1 == xor2)
		}
	}

	return res
}

type Random struct {
	seed     uint64
	hashBase uint64
}

func NewRandom() *Random                 { return &Random{seed: uint64(time.Now().UnixNano()/2 + 1)} }
func NewRandomWithSeed(seed int) *Random { return &Random{seed: uint64(seed)} }

func (r *Random) Rng() uint64 {
	r.seed ^= r.seed << 7
	r.seed ^= r.seed >> 9
	return r.seed
}
func (r *Random) Rng61() uint64 { return r.Rng() & ((1 << 61) - 1) }

type BitArrayXor struct {
	n    int
	log  int
	data []uint64
}

func NewBitArrayXor(n int) *BitArrayXor {
	return &BitArrayXor{n: n, log: bits.Len(uint(n)), data: make([]uint64, n+1)}
}

func NewBitArrayXorFrom(arr []uint64) *BitArrayXor {
	res := NewBitArrayXor(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArrayXor) Build(arr []uint64) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] ^= b.data[i]
		}
	}
}

func (b *BitArrayXor) Add(i int, xor uint64) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] ^= xor
	}
}

func (b *BitArrayXor) Query(end int) uint64 {
	res := uint64(0)
	for ; end > 0; end &= end - 1 {
		res ^= b.data[end]
	}
	return res
}

func (b *BitArrayXor) QueryRange(start, end int) uint64 {
	return b.Query(end) ^ b.Query(start)
}

type BitArray struct {
	n    int
	log  int
	data []int
}

func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

func (b *BitArray) Add(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

func (b *BitArray) Query(end int) int {
	res := 0
	for ; end > 0; end &= end - 1 {
		res += b.data[end]
	}
	return res
}

func (b *BitArray) QueryRange(start, end int) int {
	return b.Query(end) - b.Query(start)
}

func Compress(nums []int, offset int) (rank func(int) int, count int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	count = len(set)
	allNums := make([]int, 0, count)
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	rank = func(x int) int { return sort.SearchInts(allNums, x) + offset }
	return
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
	operations := make([][3]int, q)
	for i := 0; i < q; i++ {
		op := io.NextInt()
		if op == 1 {
			pos, val := io.NextInt(), io.NextInt()
			pos--
			operations[i] = [3]int{op, pos, val}
		} else {
			start, end := io.NextInt(), io.NextInt()
			start--
			operations[i] = [3]int{op, start, end}
		}
	}

	res := 由乃与大母神原型和偶像崇拜(nums, operations)
	for _, ok := range res {
		if ok {
			io.Println("damushen")
		} else {
			io.Println("yuanxing")
		}
	}
}
