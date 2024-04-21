package main

func main() {

}

type AVLTreeBitVector struct {
}

type Value = int32

func NewAVLTreeBitVector(n int32, f func(i int32) Value) *AVLTreeBitVector {
	return &AVLTreeBitVector{}
}

func (t *AVLTreeBitVector) Reserve(n int32)                       {}
func (t *AVLTreeBitVector) Insert(index int32, key Value)         {}
func (t *AVLTreeBitVector) Pop(index int32) int32                 { return 0 }
func (t *AVLTreeBitVector) Set(index int32, key Value)            {}
func (t *AVLTreeBitVector) Get(index int32) Value                 { return 0 }
func (t *AVLTreeBitVector) Count0(end int32) int32                { return 0 }
func (t *AVLTreeBitVector) Count1(end int32) int32                { return 0 }
func (t *AVLTreeBitVector) Count(end int32, key Value) int32      { return 0 }
func (t *AVLTreeBitVector) Kth0(k int32) int32                    { return 0 }
func (t *AVLTreeBitVector) Kth1(k int32) int32                    { return 0 }
func (t *AVLTreeBitVector) Kth(k int32, key Value) int32          { return 0 }
func (t *AVLTreeBitVector) Len() int32                            { return 0 }
func (t *AVLTreeBitVector) ToList() []Value                       { return nil }
func (t *AVLTreeBitVector) Debug()                                {}
func (t *AVLTreeBitVector) _build(n int32, f func(i int32) Value) {}
func (t *AVLTreeBitVector) _rotateLeft(node int32) int32          { return 0 }
func (t *AVLTreeBitVector) _rotateRight(node int32) int32         { return 0 }
func (t *AVLTreeBitVector) _rotateLeftRight(node int32) int32     { return 0 }
func (t *AVLTreeBitVector) _rotateRightLeft(node int32) int32     { return 0 }
func (t *AVLTreeBitVector) _updateBalance(node int32)             {}
func (t *AVLTreeBitVector) _pref(r int32) int32                   { return 0 }
func (t *AVLTreeBitVector) _makeNode(key int32, bitLen int32) int32 {
	return 0
}
func (t *AVLTreeBitVector) _popUnder(path []int32, d int32, node int32, res int32) {}
func (t *AVLTreeBitVector) _insertAndCount1(k int32, key int32) int32              { return 0 }
func (t *AVLTreeBitVector) _accessPopAndCount1(k int32) int32                      { return 0 }
