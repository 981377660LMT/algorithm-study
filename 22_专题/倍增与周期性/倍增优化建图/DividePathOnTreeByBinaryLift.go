// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
// TODO: 与`CompressedLCA`功能保持一致，并增加拆路径的功能.
// 一共拆分成[0,log]层，每层有n个元素.
// !jumpId = level*n + index 表示第level层的第index个元素(0<=level<log+1,0<=index<n).

package main

import "math/bits"

func main() {

}

type DividePathOnTreeByBinaryLift struct {
	n, log int32
	size   int32
}

func NewDividePathOnTreeByBinaryLift(n int32) *DividePathOnTreeByBinaryLift {
	log := int32(bits.Len(uint(n))) - 1
	size := n * (log + 1)
	return &DividePathOnTreeByBinaryLift{n: n, log: log, size: size}
}

func (d *DividePathOnTreeByBinaryLift) EnumeratePath(u, v int32, f func(level, index int32)) {}

// func (d *DividePathOnTreeByBinaryLift) _findLca()

func (d *DividePathOnTreeByBinaryLift) Size() int32 {
	return d.size
}
