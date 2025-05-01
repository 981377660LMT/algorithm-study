// 对角线排序(inplace).
// https://leetcode.cn/problems/sort-the-matrix-diagonally/solutions/2760094/dui-jiao-xian-pai-xu-fu-yuan-di-pai-xu-p-uts8/

package main

import "sort"

type diagonalSorter struct {
	mat [][]int
	k   int // 列-行
}

func (s diagonalSorter) Len() int {
	m, n := len(s.mat), len(s.mat[0])
	return min(s.k+n, m) - max(s.k, 0)
}

func (s diagonalSorter) Less(i, j int) bool {
	minI := max(s.k, 0)
	x := s.mat[minI+i][minI+i-s.k]
	y := s.mat[minI+j][minI+j-s.k]
	return x < y
}

func (s diagonalSorter) Swap(i, j int) {
	minI := max(s.k, 0)
	p := &s.mat[minI+i][minI+i-s.k]
	q := &s.mat[minI+j][minI+j-s.k]
	*p, *q = *q, *p
}

func diagonalSort(mat [][]int) [][]int {
	m, n := len(mat), len(mat[0])
	ds := diagonalSorter{mat: mat}
	for ds.k = 1 - n; ds.k < m; ds.k++ {
		sort.Sort(ds)
	}
	return mat
}
