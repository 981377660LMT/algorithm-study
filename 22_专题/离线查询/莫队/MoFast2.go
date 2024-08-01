// MoFast-理论最优的莫队算法

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// 静态区间逆序对.
func main() {
	const eof = 0
	in := os.Stdin
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	_i, _n, buf := 0, 0, make([]byte, 1<<12)

	rc := func() byte {
		if _i == _n {
			_n, _ = in.Read(buf)
			if _n == 0 {
				return eof
			}
			_i = 0
		}
		b := buf[_i]
		_i++
		return b
	}

	// 读一个整数，支持负数
	NextInt := func() (x int) {
		neg := false
		b := rc()
		for ; '0' > b || b > '9'; b = rc() {
			if b == eof {
				return
			}
			if b == '-' {
				neg = true
			}
		}
		for ; '0' <= b && b <= '9'; b = rc() {
			x = x*10 + int(b&15)
		}
		if neg {
			return -x
		}
		return
	}
	_ = NextInt
	n, q := int32(NextInt()), int32(NextInt())

	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = int32(NextInt())
	}
	mo := NewMoFast(n, q)
	for i := int32(0); i < q; i++ {
		l, r := int32(NextInt()), int32(NextInt())
		mo.AddQuery(l, r)
	}

	newNums, origin := Discretize(nums)
	bit := newBitArrayFrom(int32(len(origin)), func(i int32) int32 { return 0 })

	count := 0
	res := make([]int, q)
	addLeft := func(i int32) {
		count += int(bit.QueryPrefix(newNums[i]))
		bit.Add(newNums[i], 1)
	}
	addRight := func(i int32) {
		count += int(bit.QueryRange(newNums[i]+1, bit.Size()))
		bit.Add(newNums[i], 1)
	}
	removeLeft := func(i int32) {
		count -= int(bit.QueryPrefix(newNums[i]))
		bit.Add(newNums[i], -1)
	}
	removeRight := func(i int32) {
		count -= int(bit.QueryRange(newNums[i]+1, bit.Size()))
		bit.Add(newNums[i], -1)
	}

	query := func(i int32) {
		res[i] = count
	}

	mo.Run(addLeft, addRight, removeLeft, removeRight, query)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}

}

type MoFast struct {
	n, q, width        int32
	left, right, order []int32
}

func NewMoFast(n, q int32) *MoFast {
	width := max32(1, max32(1, n/max32(1, int32(math.Sqrt(float64(q/2))))))
	order := make([]int32, q)
	for i := int32(0); i < q; i++ {
		order[i] = i
	}
	left, right := make([]int32, 0, q), make([]int32, 0, q)
	return &MoFast{n: n, q: q, width: width, order: order, left: left, right: right}
}

func (mo *MoFast) AddQuery(l, r int32) {
	mo.left = append(mo.left, l)
	mo.right = append(mo.right, r)
}

func (mo *MoFast) Run(
	addLeft, addRight,
	removeLeft, removeRight,
	query func(i int32),
) {
	mo.build()
	nl, nr := int32(0), int32(0)
	for _, idx := range mo.order {
		for nl > mo.left[idx] {
			nl--
			addLeft(nl)
		}
		for nr < mo.right[idx] {
			addRight(nr)
			nr++
		}
		for nl < mo.left[idx] {
			removeLeft(nl)
			nl++
		}
		for nr > mo.right[idx] {
			nr--
			removeRight(nr)
		}
		query(idx)
	}
}

func (mo *MoFast) build() {
	mo.sort()
	mo.climb(3, 5)
}

func (mo *MoFast) sort() {
	n, q, width := mo.n, mo.q, mo.width
	left, right, order := mo.left, mo.right, mo.order
	counter := make([]int32, n+1)
	buf := make([]int32, q)
	for i := int32(0); i < q; i++ {
		counter[right[i]]++
	}
	for i := int32(1); i <= n; i++ {
		counter[i] += counter[i-1]
	}
	for i := int32(0); i < q; i++ {
		counter[right[i]]--
		buf[counter[right[i]]] = i
	}
	b := make([]int32, q)
	for i := int32(0); i < q; i++ {
		b[i] = left[i] / width
	}
	counter = make([]int32, n/width+1)
	for i := int32(0); i < q; i++ {
		counter[b[i]]++
	}
	for i := int32(1); i < int32(len(counter)); i++ {
		counter[i] += counter[i-1]
	}
	for i := int32(0); i < q; i++ {
		counter[b[buf[i]]]--
		order[counter[b[buf[i]]]] = buf[i]
	}
	for i, j := int32(0), int32(0); i < q; i = j {
		bi := b[order[i]]
		j = i + 1
		for j != q && bi == b[order[j]] {
			j++
		}
		if bi&1 == 0 {
			for a, b := i, j-1; a < b; a, b = a+1, b-1 {
				order[a], order[b] = order[b], order[a]
			}
		}
	}
}

func (mo *MoFast) climb(iter, interval int32) {
	q := mo.q
	order := mo.order
	d := make([]int32, q-1)
	for i := int32(0); i < q-1; i++ {
		d[i] = mo.dist(order[i], order[i+1])
	}
	for iter > 0 {
		iter--
		for i := int32(1); i < q; i++ {
			pre1 := d[i-1]
			js, je := i+1, min32(i+interval, q-1)
			for j := je - 1; j >= js; j-- {
				pre2 := d[j]
				now1 := mo.dist(order[i-1], order[j])
				now2 := mo.dist(order[i], order[j+1])
				if now1+now2 < pre1+pre2 {
					for a, b := i, j; a < b; a, b = a+1, b-1 {
						order[a], order[b] = order[b], order[a]
					}
					for a, b := i, j-1; a < b; a, b = a+1, b-1 {
						d[a], d[b] = d[b], d[a]
					}
					pre1 = now1
					d[i-1] = now1
					d[j] = now2
				}
			}
		}
	}
}

func (mo *MoFast) dist(i, j int32) int32 {
	return abs32(mo.left[i]-mo.left[j]) + abs32(mo.right[i]-mo.right[j])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int32) (newNums []int32, origin []int32) {
	newNums = make([]int32, len(nums))
	origin = make([]int32, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}

// !Point Add Range Sum, 0-based.
type bITArray32 struct {
	n     int32
	total int32
	data  []int32
}

func newBitArray(n int32) *bITArray32 {
	res := &bITArray32{n: n, data: make([]int32, n)}
	return res
}

func newBitArrayFrom(n int32, f func(i int32) int32) *bITArray32 {
	total := int32(0)
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
		total += data[i]
	}
	for i := int32(1); i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] += data[i-1]
		}
	}
	return &bITArray32{n: n, total: total, data: data}
}

func (b *bITArray32) Add(index int32, v int32) {
	b.total += v
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *bITArray32) QueryPrefix(end int32) int32 {
	if end > b.n {
		end = b.n
	}
	res := int32(0)
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *bITArray32) QueryRange(start, end int32) int32 {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := int32(0), int32(0)
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}

func (b *bITArray32) QueryAll() int32 {
	return b.total
}

func (b *bITArray32) Size() int32 {
	return b.n
}
