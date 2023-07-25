// https://maspypy.github.io/library/ds/my_multiset.hpp
//
// SortedListWithSum/MultiSetWithSum
//
// Api:
// 	Build(values []E)
// 	Add(x E)
// 	Discard(x E)
// 	Count(x E) int
// 	GetKth(k int, suffix bool) (res, sum E)
// 	GetRange(floor E, higher E) (count int, sum E)
// 	GetAll() []E
// 	Size() int

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// abc241_d()
	abc281_e()
}

// https://atcoder.jp/contests/abc241/tasks/abc241_d
func abc241_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	sl := NewMultiSetWithSum()
	for i := 0; i < q; i++ {
		var kind, x int
		fmt.Fscan(in, &kind, &x)
		if kind == 1 {
			sl.Add(x)
			continue
		}

		var k int
		fmt.Fscan(in, &k)
		if kind == 2 {
			count, _ := sl.GetRange(-INF, x+1)
			count -= k
			if 0 <= count && count < sl.Size() {
				res, _ := sl.GetKth(count, false)
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, -1)
			}
		}

		if kind == 3 {
			count, _ := sl.GetRange(-INF, x)
			count += k - 1
			if 0 <= count && count < sl.Size() {
				res, _ := sl.GetKth(count, false)
				fmt.Fprintln(out, res)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}

}

// https://atcoder.jp/contests/abc281/tasks/abc281_e
func abc281_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sl := NewMultiSetWithSum()
	for i := 0; i < m; i++ {
		sl.Add(nums[i])
	}

	res := make([]int, 0, n-m+1)
	for i := m; i <= n; i++ {
		_, sum := sl.GetKth(k, false)
		res = append(res, sum)
		if i == n {
			break
		}
		sl.Add(nums[i])
		sl.Discard(nums[i-m])
	}

	for _, x := range res {
		fmt.Fprint(out, x, " ")
	}
}

const INF int = 1e18

const _BUCKET_RATIO int = 50
const _REBUILD_RATIO int = 170

type E = int

func e() E             { return 0 }
func op(a, b E) E      { return a + b }
func inv(a E) E        { return -a }
func less(a, b E) bool { return a < b }

type MultiSetWithSum struct {
	size  int
	data  [][]E
	sum   []E
	allSm E
}

func NewMultiSetWithSum() *MultiSetWithSum {
	return &MultiSetWithSum{data: [][]E{{}}, sum: []E{e()}, allSm: e()}
}

// 'Mutates' values and returns a new MultiSetWithSum。
func (ms *MultiSetWithSum) Build(values []E) {
	sort.Slice(values, func(i, j int) bool { return less(values[i], values[j]) })
	ms.size = len(values)
	bucketCount := int(math.Sqrt(float64(ms.size/_BUCKET_RATIO))) + 1

	ms.data = make([][]E, bucketCount)
	for i := 0; i < bucketCount; i++ {
		start := ms.size * i / bucketCount
		end := min(ms.size*(i+1)/bucketCount, ms.size)
		ms.data[i] = values[start:end]
	}

	ms.sum = make([]E, bucketCount)
	for i := 0; i < bucketCount; i++ {
		cur := e()
		for _, x := range ms.data[i] {
			cur = op(cur, x)
		}
		ms.sum[i] = cur
	}

	ms.allSm = e()
	for _, x := range ms.sum {
		ms.allSm = op(ms.allSm, x)
	}
}

func (ms *MultiSetWithSum) Add(x E) {
	if ms.size == 0 {
		ms.data[0] = append(ms.data[0], x)
		ms.size++
		ms.sum[0] = op(ms.sum[0], x)
		ms.allSm = op(ms.allSm, x)
		return
	}

	for i, v := range ms.data {
		if less(v[len(v)-1], x) && i < len(ms.data)-1 {
			continue
		}
		pos := ms._lowerBound(v, x)
		ms.data[i] = append(v[:pos], append([]E{x}, v[pos:]...)...)
		ms.size++
		ms.sum[i] = op(ms.sum[i], x)
		ms.allSm = op(ms.allSm, x)
		if len(ms.data[i]) > len(ms.data)*_REBUILD_RATIO {
			ms._rebuild()
		}
		break
	}
}

func (ms *MultiSetWithSum) Discard(x E) {
	for i, v := range ms.data {
		if less(v[len(v)-1], x) && i < len(ms.data)-1 {
			continue
		}
		pos := ms._lowerBound(v, x)
		ms.data[i] = append(v[:pos], v[pos+1:]...)
		ms.size--
		ms.sum[i] = op(ms.sum[i], inv(x))
		ms.allSm = op(ms.allSm, inv(x))
		if len(ms.data[i]) == 0 && len(ms.data) > 0 {
			ms.data = append(ms.data[:i], ms.data[i+1:]...)
			ms.sum = append(ms.sum[:i], ms.sum[i+1:]...)
		}
		break
	}
}

func (ms *MultiSetWithSum) Count(x E) int {
	res := 0
	for _, v := range ms.data {
		if less(v[len(v)-1], x) {
			continue
		}
		if less(x, v[0]) {
			break
		}
		if v[0] == v[len(v)-1] {
			res += len(v)
		} else {
			res += ms._upperBound(v, x) - ms._lowerBound(v, x)
		}
	}
	return res
}

// 返回 {第k小的元素，前k个元素的和}.当suffix为true时，返回{第k大的元素，后(n-k-1)个元素的和}
//
//	k：k个元素，0<=k<=ms.size.当k==ms.size时，返回{e()，所有数的和}
//	suffix: 是否取后缀中的元素
func (ms *MultiSetWithSum) GetKth(k int, suffix bool) (res, sum E) {
	if !(0 <= k && k <= ms.size) {
		panic("k out of range")
	}

	if suffix {
		if k == ms.size {
			return e(), ms.allSm
		}
		res, sum = ms.GetKth(ms.size-k-1, false)
		return res, op(ms.allSm, inv(op(sum, res)))
	}

	sum = e()
	for i, v := range ms.data {
		if k >= len(v) {
			k -= len(v)
			sum = op(sum, ms.sum[i])
			continue
		}
		for j := 0; j < k; j++ {
			sum = op(sum, v[j])
		}
		return v[k], sum
	}

	panic("unreachable")
}

// [floor, higher) 范围内的 {count, sum}
func (ms *MultiSetWithSum) GetRange(floor E, higher E) (count int, sum E) {
	if ms.size == 0 {
		return 0, e()
	}
	sum = e()
	for i, v := range ms.data {
		if less(v[len(v)-1], floor) {
			continue
		}
		if !less(v[0], higher) {
			break
		}
		if !less(v[0], floor) && less(v[len(v)-1], higher) {
			count += len(v)
			sum = op(sum, ms.sum[i])
			continue
		}
		for _, x := range v {
			if !less(x, floor) && less(x, higher) {
				count++
				sum = op(sum, x)
			}
		}
	}
	return
}

func (ms *MultiSetWithSum) GetAll() []E {
	res := make([]E, 0, ms.Size())
	for _, x := range ms.data {
		res = append(res, x...)
	}
	return res
}

func (ms *MultiSetWithSum) Size() int { return ms.size }

func (ms *MultiSetWithSum) _rebuild() {
	ms.Build(ms.GetAll())
}

func (ms *MultiSetWithSum) _lowerBound(nums []E, target E) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if less(nums[mid], target) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func (ms *MultiSetWithSum) _upperBound(nums []E, target E) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if less(target, nums[mid]) {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
