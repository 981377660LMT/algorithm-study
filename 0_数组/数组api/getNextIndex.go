package main

import "fmt"

func main() {
	fmt.Println(GetPrevIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
	fmt.Println(GetNextIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
	fmt.Println(GetPrevAndNextIndex([]int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}))
}

// GetPrevIndex 获取数组中相同元素的前一个元素的位置.不存在则返回-1.
func GetPrevIndex[T comparable](arr []T) []int32 {
	pool := make(map[T]int32)
	newNums := make([]int32, len(arr))
	for i, v := range arr {
		if id, ok := pool[v]; !ok {
			newId := int32(len(pool))
			pool[v] = newId
			newNums[i] = newId
		} else {
			newNums[i] = id
		}
	}

	n := int32(len(arr))
	prevs := make([]int32, n)
	valuePrevs := make([]int32, len(pool))
	for i := range valuePrevs {
		valuePrevs[i] = -1
	}
	for i, v := range newNums {
		if valuePrevs[v] != -1 {
			prevs[i] = valuePrevs[v]
		} else {
			prevs[i] = -1
		}
		valuePrevs[v] = int32(i)
	}
	return prevs
}

// GetNextIndex 获取数组中相同元素的后一个元素的位置.不存在则返回-1.
func GetNextIndex[T comparable](arr []T) []int32 {
	pool := make(map[T]int32)
	newNums := make([]int32, len(arr))
	for i, v := range arr {
		if id, ok := pool[v]; !ok {
			newId := int32(len(pool))
			pool[v] = newId
			newNums[i] = newId
		} else {
			newNums[i] = id
		}
	}

	n := int32(len(arr))
	nexts := make([]int32, n)
	valueNexts := make([]int32, len(pool))
	for i := range valueNexts {
		valueNexts[i] = -1
	}
	for i := n - 1; i >= 0; i-- {
		v := newNums[i]
		if valueNexts[v] != -1 {
			nexts[i] = valueNexts[v]
		} else {
			nexts[i] = -1
		}
		valueNexts[v] = i
	}
	return nexts
}

func GetPrevAndNextIndex[T comparable](arr []T) ([]int32, []int32) {
	pool := make(map[T]int32)
	newNums := make([]int32, len(arr))
	for i, v := range arr {
		if id, ok := pool[v]; !ok {
			newId := int32(len(pool))
			pool[v] = newId
			newNums[i] = newId
		} else {
			newNums[i] = id
		}
	}

	n := int32(len(arr))
	nexts := make([]int32, n)
	prevs := make([]int32, n)
	valueNexts := make([]int32, len(pool))
	valuePrevs := make([]int32, len(pool))
	for i := range valueNexts {
		valueNexts[i] = -1
		valuePrevs[i] = -1
	}
	for i, v := range newNums {
		if valueNexts[v] != -1 {
			prevs[i] = valueNexts[v]
		} else {
			prevs[i] = -1
		}
		valueNexts[v] = int32(i)
	}
	for i := n - 1; i >= 0; i-- {
		v := newNums[i]
		if valuePrevs[v] != -1 {
			nexts[i] = valuePrevs[v]
		} else {
			nexts[i] = -1
		}
		valuePrevs[v] = i
	}
	return prevs, nexts
}
