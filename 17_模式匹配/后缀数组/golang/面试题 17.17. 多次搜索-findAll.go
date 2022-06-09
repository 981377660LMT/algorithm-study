package main

import (
	"index/suffixarray"
	"sort"
)

func multiSearch(big string, smalls []string) [][]int {
	res := make([][]int, len(smalls))
	for i, small := range smalls {
		res[i] = indexOfAll(big, small)
	}
	return res
}

// !查询所有索引位置
func indexOfAll(rawString, searchString string) []int {
	sa := suffixarray.New([]byte(rawString))
	indexes := sa.Lookup([]byte(searchString), -1)
	sort.Ints(indexes)
	return indexes
}
