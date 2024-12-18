package main

// 合并K个有序数据结构.
// 时间复杂度 O(nlogk), 空间复杂度 O(logk).k为有序数据结构的个数，n为所有数据的总个数.
func MergeKSorted[E any](sorted []E, merge func(E, E) E) (res E) {
	n := len(sorted)
	if n == 0 {
		return
	}
	if n == 1 {
		return sorted[0]
	}
	if n == 2 {
		return merge(sorted[0], sorted[1])
	}

	var f func(start, end int) E
	f = func(start, end int) E {
		if end-start == 1 {
			return sorted[start]
		}
		mid := (start + end) >> 1
		return merge(f(start, mid), f(mid, end))
	}
	return f(0, n)
}
