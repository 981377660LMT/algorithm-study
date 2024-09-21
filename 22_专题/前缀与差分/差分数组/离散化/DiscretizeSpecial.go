package main

import "fmt"

func main() {
	vToId, idToV := DiscretizeSpecial(10, func(i int32) bool { return i%2 == 0 })
	fmt.Println(vToId, idToV)
}

// 给定元素0~n-1,对数组中的某些特殊元素进行离散化.
// 返回离散化后的数组id和id对应的值.
// 特殊元素的id为0~len(idToV)-1, 非特殊元素的id为-1.
func DiscretizeSpecial(n int32, isSpecial func(i int32) bool) (vToId []int32, idToV []int32) {
	vToId = make([]int32, n)
	idToV = []int32{}
	ptr := int32(0)
	for i := int32(0); i < n; i++ {
		if isSpecial(i) {
			vToId[i] = ptr
			ptr++
			idToV = append(idToV, i)
		} else {
			vToId[i] = -1
		}
	}
	idToV = idToV[:len(idToV):len(idToV)]
	return
}
