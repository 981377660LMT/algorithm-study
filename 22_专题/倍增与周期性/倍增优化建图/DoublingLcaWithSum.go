// DividePathOnTreeBinaryLift/DoublingLca
// 倍增拆树上路径`path(from,to)`：倍增拆点将树上的一段路径拆成logn个点
// TODO: 与`CompressedLCA`功能保持一致，并增加拆路径的功能.

package main

import "math/bits"

func main() {

}

type DivideIntervalBinaryLift struct {
	n, log int32
	size   int32
}

func NewDivideIntervalBinaryLift(n int32) *DivideIntervalBinaryLift {
	log := int32(bits.Len(uint(n))) - 1
	size := n * (log + 1)
	return &DivideIntervalBinaryLift{n: n, log: log, size: size}
}

func (d *DivideIntervalBinaryLift) EnumerateRange(start int32, end int, f func(jumpId int32)) {}

func (d *DivideIntervalBinaryLift) EnumerateRange2(start1, end1 int, start2, end2 int32, f func(jumpId int32)) {
}

func (d *DivideIntervalBinaryLift) Size() int32 {
	return d.size
}
