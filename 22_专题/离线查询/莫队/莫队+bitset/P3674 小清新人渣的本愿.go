// 莫队+bitset

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"math/bits"
	"os"
	"sort"
	"strconv"
)

// 给一个数组，每次询问给定op,start, end, x，(1->差,2->和,3->积)
// 问区间[start,end)内是否能找出两个数，它们的差/和/积为x。
// 询问次数、值域均不超过1e5
//
// !1.判断差k是否存在:右移k位，然后与bitset做与运算，如果结果不为0，则存在
// !2.判断和k是否存在:rSet右移(n-k)位，然后与bitset做与运算，如果结果不为0，则存在
// !3.枚举x的因子，并判断两者是否同时存在.
func 小清新人渣的本愿(nums []int, queries [][4]int) []bool {
	n, q := len(nums), len(queries)
	mo := NewMoAlgo(n, q)
	max_ := 0
	for _, v := range nums {
		max_ = max(max_, v)
	}
	for _, query := range queries {
		start, end := query[1], query[2]
		mo.AddQuery(start, end)
		max_ = max(max_, query[3])
	}

	set := _NewBitset(max_)
	rSet := _NewBitset(max_)
	counter := make([]int, max_+1)

	res := make([]bool, q)
	mo.Run(
		func(index, _ int) {
			num := nums[index]
			counter[num]++
			if counter[num] == 1 {
				set.Set(num)
				rSet.Set(max_ - num)
			}
		},
		func(index, _ int) {
			num := nums[index]
			counter[num]--
			if counter[num] == 0 {
				set.Reset(num)
				rSet.Reset(max_ - num)
			}
		},
		func(qid int) {
			op, k := queries[qid][0], queries[qid][3]
			if op == 1 {
				// !区间内是否存在差为k
				tmp := set.Copy()
				tmp.Rsh(k)
				tmp.IAnd(set)
				res[qid] = tmp.Index1() != -1
			} else if op == 2 {
				// !区间内是否存在和为k
				tmp := rSet.Copy()
				tmp.Rsh(max_ - k)
				tmp.IAnd(set)
				res[qid] = tmp.Index1() != -1
			} else if op == 3 {
				// !区间内是否存在积为k
				factors := GetFactors(k)
				for _, factor := range factors {
					if set.Has(factor) && set.Has(k/factor) {
						res[qid] = true
						break
					}
				}
			}
		},
	)

	return res
}

type MoAlgo struct {
	queryOrder int
	chunkSize  int
	buckets    [][]query
}

type query struct{ qi, left, right int }

func NewMoAlgo(n, q int) *MoAlgo {
	chunkSize := max(1, n/max(1, int(math.Sqrt(float64(q*2/3)))))
	buckets := make([][]query, n/chunkSize+1)
	return &MoAlgo{chunkSize: chunkSize, buckets: buckets}
}

// 添加一个查询，查询范围为`左闭右开区间` [left, right).
//
//	0 <= left <= right <= n
func (mo *MoAlgo) AddQuery(left, right int) {
	index := left / mo.chunkSize
	mo.buckets[index] = append(mo.buckets[index], query{mo.queryOrder, left, right})
	mo.queryOrder++
}

// 返回每个查询的结果.
//
//	add: 将数据添加到窗口. delta: 1 表示向右移动，-1 表示向左移动.
//	remove: 将数据从窗口移除. delta: 1 表示向右移动，-1 表示向左移动.
//	query: 查询窗口内的数据.
func (mo *MoAlgo) Run(
	add func(index, delta int),
	remove func(index, delta int),
	query func(qid int),
) {
	left, right := 0, 0

	for i, bucket := range mo.buckets {
		if i&1 == 1 {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right < bucket[j].right })
		} else {
			sort.Slice(bucket, func(i, j int) bool { return bucket[i].right > bucket[j].right })
		}

		for _, q := range bucket {
			// !窗口扩张
			for left > q.left {
				left--
				add(left, -1)
			}
			for right < q.right {
				add(right, 1)
				right++
			}

			// !窗口收缩
			for left < q.left {
				remove(left, 1)
				left++
			}
			for right > q.right {
				right--
				remove(right, -1)
			}

			query(q.qi)
		}
	}
}

type _Bitset []uint

func _NewBitset(n int) _Bitset { return make(_Bitset, n>>6+1) } // (n+64-1)>>6

func (b _Bitset) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b _Bitset) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b _Bitset) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b _Bitset) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b _Bitset) Copy() _Bitset {
	res := make(_Bitset, len(b))
	copy(res, b)
	return res
}

func (bs _Bitset) Clear() {
	for i := range bs {
		bs[i] = 0
	}
}

// 遍历所有 1 的位置
// 如果对范围有要求，可在 f 中 return p < n
func (b _Bitset) Foreach(f func(p int) (Break bool)) {
	for i, v := range b {
		for ; v > 0; v &= v - 1 {
			j := i<<6 | bits.TrailingZeros(v)
			if f(j) {
				return
			}
		}
	}
}

// 返回第一个 0 的下标，若不存在则返回-1。
func (b _Bitset) Index0() int {
	for i, v := range b {
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
	}
	return -1
}

// 返回第一个 1 的下标，若不存在则返回-1。
func (b _Bitset) Index1() int {
	for i, v := range b {
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 1 的下标，若不存在则返回-1。
func (b _Bitset) Next1(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i] & (^uint(0) << (p & 63)) // mask off bits below bound
		if v != 0 {
			return i<<6 | bits.TrailingZeros(v)
		}
		for i++; i < len(b); i++ {
			if b[i] != 0 {
				return i<<6 | bits.TrailingZeros(b[i])
			}
		}
	}
	return -1
}

// 返回下标 >= p 的第一个 0 的下标，若不存在则返回-1。
func (b _Bitset) Next0(p int) int {
	if i := p >> 6; i < len(b) {
		v := b[i]
		if p&63 > 0 {
			v |= ^(^uint(0) << (p & 63))
		}
		if ^v != 0 {
			return i<<6 | bits.TrailingZeros(^v)
		}
		for i++; i < len(b); i++ {
			if ^b[i] != 0 {
				return i<<6 | bits.TrailingZeros(^b[i])
			}
		}
	}
	return -1
}

// 返回最后第一个 1 的下标，若不存在则返回 -1
func (b _Bitset) LastIndex1() int {
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0 {
			return i<<6 | (bits.Len(b[i]) - 1) // 如果再 +1，需要改成 i<<6 + bits.Len(b[i])
		}
	}
	return -1
}

// += 1 << i，模拟进位
func (b _Bitset) Add(i int) { b.FlipRange(i, b.Next0(i)) }

// -= 1 << i，模拟借位
func (b _Bitset) Sub(i int) { b.FlipRange(i, b.Next1(i)) }

// 判断 [l,r] 范围内的数是否全为 0
// https://codeforces.com/contest/1107/problem/D（标准做法是二维前缀和）
func (b _Bitset) All0(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == 0
	}
	if b[i]>>(l&63) != 0 {
		return false
	}
	for i++; i < r>>6; i++ {
		if b[i] != 0 {
			return false
		}
	}
	mask := ^uint(0) << (r & 63)
	return b[r>>6]&^mask == 0
}

// 判断 [l,r] 范围内的数是否全为 1
func (b _Bitset) All1(l, r int) bool {
	i := l >> 6
	if i == r>>6 {
		mask := ^uint(0)<<(l&63) ^ ^uint(0)<<(r&63)
		return b[i]&mask == mask
	}
	mask := ^uint(0) << (l & 63)
	if b[i]&mask != mask {
		return false
	}
	for i++; i < r>>6; i++ {
		if ^b[i] != 0 {
			return false
		}
	}
	mask = ^uint(0) << (r & 63)
	return ^(b[r>>6] | mask) == 0
}

// 反转 [l,r) 范围内的比特
// https://codeforces.com/contest/1705/problem/E
func (b _Bitset) FlipRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] ^= maskL ^ maskR
		return
	}
	b[i] ^= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^b[i]
	}
	b[i] ^= ^maskR
}

// 将 [l,r) 范围内的比特全部置 1
func (b _Bitset) SetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] |= maskL ^ maskR
		return
	}
	b[i] |= maskL
	for i++; i < r>>6; i++ {
		b[i] = ^uint(0)
	}
	b[i] |= ^maskR
}

// 将 [l,r) 范围内的比特全部置 0
func (b _Bitset) ResetRange(l, r int) {
	maskL, maskR := ^uint(0)<<(l&63), ^uint(0)<<(r&63)
	i := l >> 6
	if i == r>>6 {
		b[i] &= ^maskL | maskR
		return
	}
	b[i] &= ^maskL
	for i++; i < r>>6; i++ {
		b[i] = 0
	}
	b[i] &= maskR
}

// 左移 k 位
func (b _Bitset) Lsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return
	}
	if offset == 0 {
		// Fast path
		copy(b[shift:], b)
	} else {
		for i := len(b) - 1; i > shift; i-- {
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(64-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
}

// 右移 k 位
func (b _Bitset) Rsh(k int) {
	if k == 0 {
		return
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return
	}
	lim := len(b) - 1 - shift
	if offset == 0 {
		// Fast path
		copy(b, b[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(64-offset)
		}
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
}

// 借用 bits 库中的一些方法的名字
func (b _Bitset) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}
func (b _Bitset) TrailingZeros() int { return b.Index1() }
func (b _Bitset) Len() int           { return b.LastIndex1() + 1 }

// 下面几个方法均需保证长度相同
func (b _Bitset) Equals(c _Bitset) bool {
	for i, v := range b {
		if v != c[i] {
			return false
		}
	}
	return true
}

func (b _Bitset) HasSubset(c _Bitset) bool {
	for i, v := range b {
		if v|c[i] != v {
			return false
		}
	}
	return true
}

// 将 c 的元素合并进 b
func (b _Bitset) IOr(c _Bitset) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b _Bitset) Or(c _Bitset) _Bitset {
	res := make(_Bitset, len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func (b _Bitset) IAnd(c _Bitset) {
	for i, v := range c {
		b[i] &= v
	}
}

func (b _Bitset) And(c _Bitset) _Bitset {
	res := make(_Bitset, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}

func (b _Bitset) IXor(c _Bitset) _Bitset {
	for i, v := range c {
		b[i] ^= v
	}
	return b
}

func (b _Bitset) Xor(c _Bitset) _Bitset {
	res := make(_Bitset, len(b))
	for i, v := range b {
		res[i] = v ^ c[i]
	}
	return res
}

// 返回 n 的所有因子. O(n^0.5).
func GetFactors(n int) []int {
	if n <= 0 {
		return nil
	}
	small := []int{}
	big := []int{}
	upper := int(math.Sqrt(float64(n)))
	for f := 1; f <= upper; f++ {
		if n%f == 0 {
			small = append(small, f)
			big = append(big, n/f)
		}
	}
	if small[len(small)-1] == big[len(big)-1] {
		big = big[:len(big)-1]
	}
	for i, j := 0, len(big)-1; i < j; i, j = i+1, j-1 {
		big[i], big[j] = big[j], big[i]
	}
	res := append(small, big...)
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		op, l, r, x := io.NextInt(), io.NextInt(), io.NextInt(), io.NextInt()
		l--
		queries[i] = [4]int{op, l, r, x}
	}

	res := 小清新人渣的本愿(nums, queries)
	for _, b := range res {
		if b {
			io.Println("hana")
		} else {
			io.Println("bi")
		}
	}
}
