// // XorTrie.
// type BinaryTrie struct {
// 	_multiset                        bool
// 	_maxLog, _xEnd, _maxV, _lazy     int
// 	_vList, _edges, _size, _endCount []int
// }

// // max: max of x
// // addLimit: max number of add and query operations
// // allowMultipleElements: allow multiple elements with the same value (SortedList or SortedSet)
// func NewBinaryTrie(max, addLimit int, allowMultipleElements bool) *BinaryTrie {
// 	maxLog := bits.Len(uint(max))
// 	n := maxLog*addLimit + 1
// 	edges := make([]int, 2*n)
// 	for i := range edges {
// 		edges[i] = -1
// 	}

// 	return &BinaryTrie{
// 		_multiset: allowMultipleElements,
// 		_maxLog:   maxLog,
// 		_xEnd:     1 << maxLog,
// 		_vList:    make([]int, maxLog+1),
// 		_edges:    edges,
// 		_size:     make([]int, n),
// 		_endCount: make([]int, n),
// 	}
// }

class XorTrie {
  private readonly _multiset: boolean
  private readonly _maxLog: number
  private readonly _xEnd: number
  private _maxV: number
  private readonly _vList: number[]
  private readonly _edges: Int32Array
  private readonly _size: Int32Array
  private readonly _endCount: Int32Array
  private _lazy: number

  constructor(
    options: {
      max?: number
      addLimit?: number
      allowMultipleElements?: boolean
    } = {}
  ) {
    const { max = (1 << 30) + 10, addLimit = 1e5 + 10, allowMultipleElements = true } = options
    const maxLog = 32 - Math.clz32(max)
    const n = maxLog * addLimit + 1
    this._multiset = allowMultipleElements
    this._maxLog = maxLog
    this._xEnd = 1 << maxLog
    this._maxV = 0
    this._vList = Array(maxLog + 1).fill(0)
    this._edges = new Int32Array(2 * n).fill(-1)
    this._size = new Int32Array(n)
    this._endCount = new Int32Array(n)
    this._lazy = 0
  }

  add(int32: number): void {
    if (int32 < 0 || int32 >= this._xEnd) return
    int32 ^= this._lazy
    let v = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const d = (int32 >>> i) & 1
      if (this._edges[2 * v + d] === -1) {
        this._edges[2 * v + d] = ++this._maxV
      }
      v = this._edges[2 * v + d]
      this._vList[i] = v
    }

    if (this._multiset || !this._endCount[v]) {
      this._endCount[v]++
      for (let i = 0; i < this._vList.length; i++) {
        this._size[this._vList[i]]++
      }
    }
  }

  discard(int32: number): void {
    if (int32 < 0 || int32 >= this._xEnd) return
    int32 ^= this._lazy
    let v = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const d = (int32 >>> i) & 1
      if (this._edges[2 * v + d] === -1) return
      v = this._edges[2 * v + d]
      this._vList[i] = v
    }
    if (this._endCount[v] > 0) {
      this._endCount[v]--
      for (let i = 0; i < this._vList.length; i++) {
        this._size[this._vList[i]]--
      }
    }
  }

  /**
   * 删除count个int32.
   * @param count 删除的个数, 默认为1.设为-1时删除所有.
   */
  erase(int32: number, count = 1): void {
    if (int32 < 0 || int32 >= this._xEnd) return
  }

  count(int32: number): number {
    if (int32 < 0 || int32 >= this._xEnd) return 0
  }

  bisectLeft(int32: number): number {}

  bisectRight(int32: number): number {
    return this.bisectLeft(int32 + 1)
  }

  indexOf(int32: number): number {}

  at(index: number): number {}

  xorAll(int32: number): number {
    this._lazy ^= int32
  }

  has(int32: number): boolean {
    return !!this.count(int32)
  }

  forEach(callbackfn: (value: number, index: number) => void): void {}

  toString(): string {}

  get size(): number {
    return this._size[0]
  }

  get min(): number {}

  get max(): number {}
}

export { XorTrie, XorTrie as BinaryTrie }
