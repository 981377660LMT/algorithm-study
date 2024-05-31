package main

import "fmt"

func main() {
	data := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9}
	maxSize := int32(4)
	chunked := Chunk(int32(len(data)), func(i int32) int32 { return data[i] }, maxSize)
	fmt.Println(chunked)
}

func Chunk[T any](n int32, f func(i int32) T, maxSize int32) [][]T {
	if maxSize <= 1 {
		res := make([][]T, n)
		for i := int32(0); i < n; i++ {
			res[i] = []T{f(i)}
		}
		return res
	}

	res := make([][]T, 0, (n+maxSize-1)/maxSize)
	ptr := int32(0)
	for ptr < n {
		curGroup := make([]T, 0, maxSize)
		for i := int32(0); i < maxSize && ptr < n; i++ {
			curGroup = append(curGroup, f(ptr))
			ptr++
		}
		curGroup = curGroup[:len(curGroup):len(curGroup)] // clip
		res = append(res, curGroup)
	}
	return res
}
