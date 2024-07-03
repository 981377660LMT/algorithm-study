// 双数组Trie(DAT)
// https://github.com/euclidr/darts/blob/master/darts.go
// 1.build double array data structure from ordered strings
// 2.check whether a string exists and get its index
// 3.scan matched prefixes of a string

package main

import (
	"fmt"
	"sort"
)

func main() {
	builder := DoubleArrayBuilder{}
	keyset := []string{"印度", "印度尼西亚", "印加帝国", "瑞士", "瑞典", "巴基斯坦", "巴勒斯坦", "以色列", "巴比伦", "土耳其"}
	sort.Strings(keyset)

	// Build darts
	builder.Build(keyset)

	// ExactMatchSearch
	key := "印度"
	result, matched := builder.ExactMatchSearch(key)
	if !matched {
		fmt.Println("not found")
		return
	}
	fmt.Println(keyset[result])

	// CommonPrefixSearch
	values := builder.CommonPrefixSearch("印度尼西亚啊")
	fmt.Printf("%s, %s", keyset[values[0]], keyset[values[1]])

}

type DoubleArray struct {
	units UnitPool
}

func (da *DoubleArray) FromBytes(byts []byte) error {
	if len(byts)%4 != 0 {
		return fmt.Errorf("invalid length of bytes")
	}

	cnt := len(byts) / 4
	var units UnitPool
	units.resize(uint32(cnt))
	for i := 0; i < cnt; i++ {
		s := i * 4
		u := uint32(byts[s])<<24 | uint32(byts[s+1])<<16 | uint32(byts[s+2])<<8 | uint32(byts[s+3])
		units.set(uint32(i), Unit(u))
	}
	da.units = units
	return nil
}

func (da *DoubleArray) ToBytes() (byts []byte) {
	cnt := int(da.units.size())
	byts = make([]byte, cnt*4)
	for i := 0; i < cnt; i++ {
		s := i * 4
		u := *da.units.at(uint32(i))
		byts[s] = byte(u >> 24)
		byts[s+1] = byte(u >> 16 & 0xFF)
		byts[s+2] = byte(u >> 8 & 0xFF)
		byts[s+3] = byte(u & 0xFF)
	}
	return
}

func (da *DoubleArray) ExactMatchSearch(key string) (value int, matched bool) {
	if da.units.size() == 0 {
		return 0, false
	}

	bkey := []byte(key)

	unit := *da.units.at(0)
	var nodePos uint32
	for _, label := range bkey {
		nodePos ^= unit.offset() ^ uint32(label)
		unit = *da.units.at(nodePos)
		if unit.label() != uint32(label) {
			return 0, false
		}
	}

	if !unit.hasLeaf() {
		return 0, false
	}

	unit = *da.units.at(nodePos ^ unit.offset())
	return unit.value(), true
}

func (da *DoubleArray) CommonPrefixSearch(key string) (values []int) {
	values = make([]int, 0)
	if da.units.size() == 0 {
		return
	}

	bkey := []byte(key)
	unit := *da.units.at(0)
	nodePos := unit.offset()
	for _, label := range bkey {
		nodePos ^= uint32(label)
		unit = *da.units.at(nodePos)
		if unit.label() != uint32(label) {
			return
		}

		nodePos ^= unit.offset()
		if unit.hasLeaf() {
			values = append(values, da.units.at(nodePos).value())
		}
	}
	return
}

type Unit uint32

func (u *Unit) hasLeaf() bool {
	return ((*u >> 8) & 1) == 1
}

func (u *Unit) value() int {
	return int(uint32(*u) & ((1 << 31) - 1))
}

func (u *Unit) label() uint32 {
	return uint32(*u) & ((1 << 31) | 0xFF)
}

func (u *Unit) offset() uint32 {
	i := uint32(*u)
	return (i >> 10) << ((i & (1 << 9)) >> 6)
}

func (u *Unit) setHasLeaf(hasLeaf bool) {
	if hasLeaf {
		*u |= Unit(1) << 8
	} else {
		*u &= ^(Unit(1) << 8)
	}
}

// value must be positive?
func (u *Unit) setValue(value uint32) {
	*u = Unit(value | (1 << 31))
}

func (u *Unit) setLabel(label uint8) {
	*u = (*u & ^Unit(0xFF)) | Unit(label)
}

func (u *Unit) setOffset(offset uint32) {
	if offset >= 1<<29 {
		panic("failed to modify unit: too large offset")
	}

	*u &= (1 << 31) | (1 << 8) | 0xFF
	if offset < (1 << 21) {
		*u |= (Unit(offset) << 10)
	} else {
		*u |= ((Unit(offset) << 2) | (1 << 9))
	}
}

type ExtraUnit struct {
	prev    uint32
	next    uint32
	isFixed bool
	isUsed  bool
}

// constants
const (
	BlockSize      uint32 = 256
	NumExtraBlocks uint32 = 16
	NumExtra       uint32 = BlockSize * NumExtraBlocks

	UpperMask uint32 = 0xFF << 21
	LowerMask uint32 = 0xFF
)

type UnitPool struct {
	pool []Unit
}

func (up *UnitPool) set(idx uint32, u Unit) {
	up.pool[idx] = u
}

func (up *UnitPool) at(idx uint32) *Unit {
	return &up.pool[idx]
}

func (up *UnitPool) size() uint32 {
	return uint32(len(up.pool))
}

func (up *UnitPool) resize(size uint32) {
	if uint32(len(up.pool)) > size {
		up.pool = up.pool[:size]
	}
	if size > uint32(cap(up.pool)) {
		up.resizeBuf(size)
	}
	for uint32(len(up.pool)) < size {
		up.pool = append(up.pool, 0)
	}
}

func (up *UnitPool) reserve(size uint32) {
	if size > uint32(cap(up.pool)) {
		up.resizeBuf(size)
	}
}

func (up *UnitPool) resizeBuf(size uint32) {
	var capacity uint32
	if size >= uint32(cap(up.pool))*2 {
		capacity = size
	} else {
		capacity = 1
		for capacity < size {
			capacity <<= 1
		}
	}

	pool2 := make([]Unit, len(up.pool), capacity)
	copy(pool2, up.pool)
	up.pool = pool2
}

type DoubleArrayBuilder struct {
	DoubleArray
	labels     []byte
	extrasHead uint32

	_extras []ExtraUnit
}

func (bd *DoubleArrayBuilder) resetExtra() {
	bd._extras = make([]ExtraUnit, NumExtra)
}

func (bd *DoubleArrayBuilder) resetLabel() {
	bd.labels = make([]byte, 0, 256)
}

func (bd *DoubleArrayBuilder) extras(id uint32) *ExtraUnit {
	return &bd._extras[id%NumExtra]
}

func (bd *DoubleArrayBuilder) Build(keyset []string) {
	ks := make([][]byte, len(keyset))
	for i := range keyset {
		ks[i] = []byte(keyset[i] + "\x00")
	}
	bd.buildFromKeyset(ks)
}

func (bd *DoubleArrayBuilder) buildFromKeyset(keyset [][]byte) {
	var numUnits uint32 = 1
	for numUnits < uint32(len(keyset)) {
		numUnits <<= 1
	}

	bd.units.reserve(numUnits)
	bd.resetExtra()
	bd.resetLabel()

	bd.reserveID(0)
	bd.extras(0).isUsed = true
	bd.units.at(0).setOffset(1)
	bd.units.at(0).setLabel(0)

	if len(keyset) > 0 {
		bd.buildFromKeysetRange(keyset, 0, uint32(len(keyset)), 0, 0)
	}

	bd.fixAllBlocks()

	bd._extras = nil
	bd.labels = nil
}

func (bd *DoubleArrayBuilder) buildFromKeysetRange(
	keyset [][]byte,
	begin uint32,
	end uint32,
	depth uint32,
	dicID uint32) {
	offset := bd.arrangeFromKeyset(keyset, begin, end, depth, dicID)

	for begin < end {
		if keyset[begin][depth] != 0 {
			break
		}
		begin++
	}
	if begin == end {
		return
	}

	lastBegin := begin
	lastLabel := keyset[begin][depth]
	begin++
	for begin < end {
		label := keyset[begin][depth]
		if label != lastLabel {
			bd.buildFromKeysetRange(
				keyset,
				lastBegin,
				begin,
				depth+1,
				offset^uint32(lastLabel))
			lastBegin = begin
			lastLabel = keyset[begin][depth]
		}
		begin++
	}
	bd.buildFromKeysetRange(keyset, lastBegin, end, depth+1, offset^uint32(lastLabel))
}

func (bd *DoubleArrayBuilder) arrangeFromKeyset(
	keyset [][]byte,
	begin uint32,
	end uint32,
	depth uint32,
	dicID uint32) (offset uint32) {
	bd.labels = bd.labels[:0]

	var value uint32 = 0xFFFFFFFF
	for i := begin; i < end; i++ {
		label := keyset[i][depth]
		if label == 0 {
			if value == 0xFFFFFFFF {
				value = i
			}
		}

		if len(bd.labels) == 0 {
			bd.labels = append(bd.labels, label)
		} else if label != bd.labels[len(bd.labels)-1] {
			if label < bd.labels[len(bd.labels)-1] {
				panic("failed to build double-array: wrong key order")
			}
			bd.labels = append(bd.labels, label)
		}
	}

	offset = bd.findValidOffset(dicID)
	bd.units.at(dicID).setOffset(dicID ^ offset)

	for i := 0; i < len(bd.labels); i++ {
		dicChildID := offset ^ uint32(bd.labels[i])
		bd.reserveID(dicChildID)
		if bd.labels[i] == 0 {
			bd.units.at(dicID).setHasLeaf(true)
			bd.units.at(dicChildID).setValue(value)
		} else {
			bd.units.at(dicChildID).setLabel(bd.labels[i])
		}
	}

	bd.extras(offset).isUsed = true

	return offset
}

func (bd *DoubleArrayBuilder) findValidOffset(id uint32) (offset uint32) {
	if bd.extrasHead >= bd.units.size() {
		return bd.units.size() | (id & LowerMask)
	}

	unfixedID := bd.extrasHead
	for {
		offset := unfixedID ^ uint32(bd.labels[0])
		if bd.isValidOffset(id, offset) {
			return offset
		}
		unfixedID = bd.extras(unfixedID).next
		if unfixedID == bd.extrasHead {
			break
		}
	}
	return bd.units.size() | (id & LowerMask)
}

func (bd *DoubleArrayBuilder) isValidOffset(id uint32, offset uint32) bool {
	if bd.extras(offset).isUsed {
		return false
	}

	relOffset := id ^ offset
	if (relOffset&LowerMask) != 0 && (relOffset&UpperMask) != 0 {
		return false
	}

	for i := 1; i < len(bd.labels); i++ {
		if bd.extras(offset ^ uint32(bd.labels[i])).isFixed {
			return false
		}
	}

	return true
}

func (bd *DoubleArrayBuilder) reserveID(id uint32) {
	if id >= bd.units.size() {
		bd.expandUnits()
	}

	if id == bd.extrasHead {
		bd.extrasHead = bd.extras(id).next
		if bd.extrasHead == id {
			bd.extrasHead = bd.units.size()
		}
	}

	bd.extras(bd.extras(id).prev).next = bd.extras(id).next
	bd.extras(bd.extras(id).next).prev = bd.extras(id).prev
	bd.extras(id).isFixed = true
}

func (bd *DoubleArrayBuilder) numBlocks() uint32 {
	return bd.units.size() / BlockSize
}

func (bd *DoubleArrayBuilder) expandUnits() {
	var srcNumUnits = bd.units.size()
	var srcNumBlocks = bd.numBlocks()

	var destNumUnits = srcNumUnits + BlockSize
	var destNumBlocks = srcNumBlocks + 1

	if destNumBlocks > NumExtraBlocks {
		bd.fixBlock(srcNumBlocks - NumExtraBlocks)
	}

	bd.units.resize(destNumUnits)

	if destNumBlocks > NumExtraBlocks {
		for id := srcNumUnits; id < destNumUnits; id++ {
			bd.extras(id).isUsed = false
			bd.extras(id).isFixed = false
		}
	}

	for i := srcNumUnits + 1; i < destNumUnits; i++ {
		bd.extras(i - 1).next = i
		bd.extras(i).prev = i - 1
	}

	// bd.extras(srcNumUnits).prev = destNumUnits - 1
	// bd.extras(destNumUnits - 1).next = srcNumUnits

	bd.extras(srcNumUnits).prev = bd.extras(bd.extrasHead).prev
	bd.extras(destNumUnits - 1).next = bd.extrasHead

	bd.extras(bd.extras(bd.extrasHead).prev).next = srcNumUnits
	bd.extras(bd.extrasHead).prev = destNumUnits - 1
}

func (bd *DoubleArrayBuilder) fixAllBlocks() {
	var begin uint32
	if bd.numBlocks() > NumExtraBlocks {
		begin = bd.numBlocks() - NumExtraBlocks
	}
	end := bd.numBlocks()

	for blockID := begin; blockID != end; blockID++ {
		bd.fixBlock(blockID)
	}
}

func (bd *DoubleArrayBuilder) fixBlock(blockID uint32) {
	var begin = blockID * BlockSize
	var end = begin + BlockSize

	var unusedOffset uint32
	for offset := begin; offset != end; offset++ {
		if !bd.extras(offset).isUsed {
			unusedOffset = offset
			break
		}
	}

	for id := begin; id != end; id++ {
		if !bd.extras(id).isFixed {
			bd.reserveID(id)
			bd.units.at(id).setLabel(uint8(id ^ unusedOffset))
		}
	}
}
