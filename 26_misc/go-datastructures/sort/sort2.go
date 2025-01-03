// 并行排序
// 核心思想是利用多个线程对数据进行并行的分块排序（bucket sort），
// 然后递归地进行对称归并（Symmetric Merge）.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"slices"
	"sync"
	"time"
)

func main() {
	// demo()
	test()
}

func demo() {
	arr := []int{1, 3, 2, 4, 5, 7, 6, 8, 9, 10}
	MultithreadedSortSlice(arr, func(a, b int) int { return a - b })
	fmt.Println(arr)
}

func test() {
	n := int(1e7)
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = n - i
	}
	rand.Shuffle(n, func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })

	timeit := func(f func()) int {
		start := time.Now()
		f()
		return int(time.Since(start).Milliseconds())
	}

	run1 := func() {
		toBeSorted := append(slice[:0:0], slice...)
		cmp := func(a, b int) int { return a - b }
		MultithreadedSortSlice(toBeSorted, cmp)
	}

	run2 := func() {
		toBeSorted := append(slice[:0:0], slice...)
		cmp := func(a, b int) int { return a - b }
		slices.SortFunc(toBeSorted, cmp)
	}

	fmt.Println("MultithreadedSortSlice:", timeit(run1), "ms")
	fmt.Println("sort.Slice:", timeit(run2), "ms")
}

// MultithreadedSortSlice 对切片 slice 进行多线程排序，使用给定的 cmp 函数。
func MultithreadedSortSlice[S ~[]T, T any](slice S, cmp func(T, T) int) {
	var wg sync.WaitGroup

	numCPU := int64(runtime.NumCPU())
	if numCPU == 1 {
		numCPU = 2
	} else {
		numCPU = int64(prevPowerOfTwo(uint64(numCPU)))
	}

	chunks := chunk(slice, numCPU)
	wg.Add(len(chunks))
	for i := 0; i < len(chunks); i++ {
		go func(i int) {
			sortBucket(chunks[i], cmp)
			wg.Done()
		}(i)
	}
	wg.Wait()

	// 对每对相邻的块进行对称归并（SymMerge），结果存储在 todo 切片中
	todo := make([][]T, len(chunks)/2)
	for {
		todo = todo[:len(chunks)/2]
		wg.Add(len(chunks) / 2)
		for i := 0; i < len(chunks); i += 2 {
			go func(i int) {
				todo[i/2] = SymMerge(chunks[i], chunks[i+1], cmp)
				wg.Done()
			}(i)
		}
		wg.Wait()

		chunks = copyChunk(todo)
		if len(chunks) == 1 {
			break
		}
	}
}

func sortBucket[T any](slice []T, cmp func(T, T) int) {
	slices.SortFunc(slice, cmp)
}

func chunk[T any](slice []T, numParts int64) [][]T {
	parts := make([][]T, numParts)
	n := int64(len(slice))
	for i := int64(0); i < numParts; i++ {
		start := i * n / numParts
		end := (i + 1) * n / numParts
		parts[i] = slice[start:end]
	}
	return parts
}

func copyChunk[T any](chunk [][]T) [][]T {
	cp := make([][]T, len(chunk))
	copy(cp, chunk)
	return cp
}

// prevPowerOfTwo 返回 <= x 的最大的 2 的幂
func prevPowerOfTwo(x uint64) uint64 {
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >> 16)
	x = x | (x >> 32)
	return x - (x >> 1)
}

// SymMerge 假设 u, w 预先各自有序，合并得到新的有序切片。
// 如果长度差别较大，会先做拆分 + 递归 merge
func SymMerge[T any](u, w []T, cmp func(T, T) int) []T {
	lenU, lenW := len(u), len(w)
	if lenU == 0 {
		return w
	}
	if lenW == 0 {
		return u
	}

	diff := lenU - lenW
	if math.Abs(float64(diff)) > 1 {
		u1, w1, u2, w2 := prepareForSymMerge(u, w, cmp)

		lenU1 := len(u1)
		lenU2 := len(u2)
		u = append(u1, w1...)
		w = append(u2, w2...)

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			symMerge(u, 0, lenU1, len(u), cmp)
			wg.Done()
		}()
		go func() {
			symMerge(w, 0, lenU2, len(w), cmp)
			wg.Done()
		}()
		wg.Wait()

		u = append(u, w...)
		return u
	}

	u = append(u, w...)
	symMerge(u, 0, lenU, len(u), cmp)
	return u
}

// prepareForSymMerge 拆分与旋转，让长的数组拆出一段给短数组，使得后续合并更平衡
func prepareForSymMerge[T any](u, w []T, cmp func(T, T) int) ([]T, []T, []T, []T) {
	if len(u) > len(w) {
		u, w = w, u
	}
	v1, wActive, v2 := decomposeForSymMerge(len(u), w)
	i := symSearch(u, wActive, cmp)

	u1 := make([]T, i)
	copy(u1, u[:i])
	w1 := append(v1, wActive[:len(wActive)-i]...)

	u2 := make([]T, len(u)-i)
	copy(u2, u[i:])

	w2 := append(wActive[len(wActive)-i:], v2...)
	return u1, w1, u2, w2
}

// decomposeForSymMerge 从 comparators 中提取一段 w 作为 active site，前面 v1，后面 v2
func decomposeForSymMerge[T any](length int, slice []T) (v1, w, v2 []T) {
	if length >= len(slice) {
		panic(`INCORRECT PARAMS FOR SYM MERGE.`)
	}
	overhang := (len(slice) - length) / 2
	v1 = slice[:overhang]
	w = slice[overhang : overhang+length]
	v2 = slice[overhang+length:]
	return
}

// symSearch 在两个有序列表 u, w 上寻找合适的起始位置
// 这里假设 u,w 已排好序
func symSearch[T any](u, w []T, cmp func(T, T) int) int {
	start, stop := 0, len(u)
	p := len(w) - 1

	for start < stop {
		mid := (start + stop) / 2
		if cmp(w[p-mid], u[mid]) >= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}
	return start
}

// symMerge (递归版本)
func symMerge[T any](u []T, start1, start2, last int, cmp func(T, T) int) {
	if start1 < start2 && start2 < last {
		mid := (start1 + last) / 2
		n := mid + start2
		var start int
		if start2 > mid {
			start = symBinarySearch(u, n-last, mid, n-1, cmp)
		} else {
			start = symBinarySearch(u, start1, start2, n-1, cmp)
		}
		end := n - start

		symRotate(u, start, start2, end)
		symMerge(u, start1, start, mid, cmp)
		symMerge(u, mid, end, last, cmp)
	}
}

// symBinarySearch 在 [start..stop) 范围搜索，使得 u[mid] <= u[total-mid]
func symBinarySearch[T any](u []T, start, stop, total int, cmp func(T, T) int) int {
	for start < stop {
		mid := (start + stop) / 2
		if cmp(u[mid], u[total-mid]) <= 0 {
			start = mid + 1
		} else {
			stop = mid
		}
	}
	return start
}

// symRotate 做区间旋转
func symRotate[T any](u []T, start1, start2, end int) {
	i := start2 - start1
	if i == 0 {
		return
	}
	j := end - start2
	if j == 0 {
		return
	}
	if i == j {
		symSwap(u, start1, start2, i)
		return
	}
	p := start1 + i
	for i != j {
		if i > j {
			symSwap(u, p-i, p, j)
			i -= j
		} else {
			symSwap(u, p-i, p+j-i, i)
			j -= i
		}
	}
	symSwap(u, p-i, p, i)
}

// symSwap 交换 [start1..start1+length) 与 [start2..start2+length)
func symSwap[T any](u []T, start1, start2, length int) {
	for i := 0; i < length; i++ {
		u[start1+i], u[start2+i] = u[start2+i], u[start1+i]
	}
}
