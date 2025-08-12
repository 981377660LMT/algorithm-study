package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MAXQ = 500010
	LOGK = 20
	LMAX = 1e18
)

type HeavyLightDecomposition struct {
	size  []int64
	light []Segment
	heavy [][]Segment

	data []interface{}
}

type Segment struct {
	idx int64
	ofs int64
}

func NewHLD(capacity int) *HeavyLightDecomposition {
	hld := &HeavyLightDecomposition{
		size:  make([]int64, capacity),
		light: make([]Segment, capacity),
		heavy: make([][]Segment, capacity),
		data:  make([]interface{}, capacity),
	}

	for i := 0; i < capacity; i++ {
		hld.size[i] = -1
		hld.light[i] = Segment{-1, -1}
		hld.heavy[i] = make([]Segment, LOGK)
		for j := 0; j < LOGK; j++ {
			hld.heavy[i][j] = Segment{-1, -1}
		}
	}

	return hld
}

func (hld *HeavyLightDecomposition) SetLeaf(idx int, size int64, data interface{}) {
	hld.size[idx] = size
	hld.data[idx] = data
}

func (hld *HeavyLightDecomposition) BuildNode(nodeIdx int, leftChild, rightChild int,
	leftSize, rightSize, leftOffset, rightOffset int64) {

	totalSize := leftSize + rightSize
	if totalSize > LMAX {
		totalSize = LMAX
		rightSize = totalSize - leftSize
		if rightSize < 0 {
			rightSize = 0
		}
	}

	hld.size[nodeIdx] = totalSize

	var heavyChild, lightChild int
	var heavyOffset, lightOffset int64

	if leftSize >= rightSize {
		heavyChild, lightChild = leftChild, rightChild
		heavyOffset, lightOffset = leftOffset, rightOffset
	} else {
		heavyChild, lightChild = rightChild, leftChild
		heavyOffset, lightOffset = rightOffset, leftOffset
	}

	hld.light[nodeIdx] = Segment{int64(lightChild), lightOffset}

	hld.heavy[nodeIdx][0] = Segment{int64(heavyChild), heavyOffset}

	cur := heavyChild
	curOffset := heavyOffset

	for k := 1; k < LOGK; k++ {
		if cur == -1 || hld.heavy[cur][k-1].idx == -1 {
			break
		}

		nextOffset := curOffset + hld.heavy[cur][k-1].ofs
		if nextOffset > LMAX {
			break
		}

		nextCur := int(hld.heavy[cur][k-1].idx)
		hld.heavy[nodeIdx][k] = Segment{int64(nextCur), nextOffset}

		cur = nextCur
		curOffset = nextOffset
	}
}

func (hld *HeavyLightDecomposition) Query(nodeIdx int, position int64) interface{} {
	return hld.queryRecursive(int64(nodeIdx), position)
}

func (hld *HeavyLightDecomposition) queryRecursive(nodeIdx, position int64) interface{} {
	if nodeIdx <= 1 {
		return hld.data[nodeIdx]
	}

	cur := nodeIdx
	pos := position

	for k := LOGK - 1; k >= 0; k-- {
		seg := hld.heavy[cur][k]
		if seg.idx == -1 {
			continue
		}

		targetSize := hld.size[seg.idx]
		if seg.ofs <= pos && pos < seg.ofs+targetSize {
			pos -= seg.ofs
			cur = seg.idx
		}
	}

	for hld.light[cur].idx != -1 {
		seg := hld.light[cur]
		targetSize := hld.size[seg.idx]

		if seg.ofs <= pos && pos < seg.ofs+targetSize {
			pos -= seg.ofs
			cur = seg.idx
			break
		} else {
			break
		}
	}

	return hld.queryRecursive(cur, pos)
}

func (hld *HeavyLightDecomposition) GetSize(idx int) int64 {
	if idx < 0 || idx >= len(hld.size) {
		return 0
	}
	return hld.size[idx]
}

type BinaryCatSolver struct {
	hld *HeavyLightDecomposition
}

func NewBinaryCatSolver() *BinaryCatSolver {
	solver := &BinaryCatSolver{
		hld: NewHLD(MAXQ),
	}

	solver.hld.SetLeaf(0, 1, byte('0'))
	solver.hld.SetLeaf(1, 1, byte('1'))

	return solver
}

func (solver *BinaryCatSolver) Concat(nodeIdx, leftIdx, rightIdx int) {
	leftSize := solver.hld.GetSize(leftIdx)
	rightSize := solver.hld.GetSize(rightIdx)

	solver.hld.BuildNode(nodeIdx, leftIdx, rightIdx,
		leftSize, rightSize, 0, leftSize)
}

func (solver *BinaryCatSolver) Query(nodeIdx int, position int64) byte {
	result := solver.hld.Query(nodeIdx, position-1)
	return result.(byte)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var Q int
	fmt.Fscanf(reader, "%d\n", &Q)

	solver := NewBinaryCatSolver()

	for i := 0; i < Q; i++ {
		var L, R, X int64
		fmt.Fscanf(reader, "%d %d %d\n", &L, &R, &X)

		solver.Concat(i+2, int(L), int(R))

		result := solver.Query(i+2, X)
		fmt.Fprintf(writer, "%c\n", result)
	}
}
