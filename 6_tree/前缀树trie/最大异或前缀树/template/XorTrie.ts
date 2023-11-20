/* eslint-disable max-len */
/* eslint-disable no-inner-declarations */

class XorTrie {
  private readonly _multiset: boolean
  private readonly _maxLog: number
  private readonly _xEnd: number
  private readonly _vList: number[]
  private readonly _edges: Int32Array
  private readonly _size: Int32Array
  private readonly _endCount: Int32Array
  private _maxV: number
  private _lazy: number

  constructor(
    options: {
      /** 最大值, 默认为2^30. */
      maxInt32?: number
      /** 最多添加的元素个数, 默认为1e5. */
      addLimit?: number
      /** 是否允许重复元素, 默认为true. */
      allowMultipleElements?: boolean
    } = {}
  ) {
    const { maxInt32 = (1 << 30) + 10, addLimit = 1e5 + 10, allowMultipleElements = true } = options
    const maxLog = 32 - Math.clz32(maxInt32)
    const n = maxLog * addLimit + 1
    this._multiset = allowMultipleElements
    this._maxLog = maxLog
    this._xEnd = 2 ** maxLog // !不要用1 << maxLog, 会溢出
    this._vList = Array(maxLog + 1).fill(0)
    this._edges = new Int32Array(2 * n).fill(-1)
    this._size = new Int32Array(n)
    this._endCount = new Int32Array(n)
    this._maxV = 0
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
    int32 ^= this._lazy
    let v = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const d = (int32 >>> i) & 1
      if (this._edges[2 * v + d] === -1) return
      v = this._edges[2 * v + d]
      this._vList[i] = v
    }
    if (count === -1 || this._endCount[v] < count) {
      count = this._endCount[v]
    }
    if (this._endCount[v] > 0) {
      this._endCount[v] -= count
      for (let i = 0; i < this._vList.length; i++) {
        this._size[this._vList[i]] -= count
      }
    }
  }

  count(int32: number): number {
    if (int32 < 0 || int32 >= this._xEnd) return 0
    int32 ^= this._lazy
    let v = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const d = (int32 >>> i) & 1
      if (this._edges[2 * v + d] === -1) return 0
      v = this._edges[2 * v + d]
    }
    return this._endCount[v]
  }

  bisectLeft(int32: number): number {
    if (int32 < 0) return 0
    if (int32 >= this._xEnd) return this.size
    let v = 0
    let res = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const d = (int32 >>> i) & 1
      const left = (this._lazy >>> i) & 1
      let lc = this._edges[2 * v]
      let rc = this._edges[2 * v + 1]
      if (left) {
        const tmp = lc
        lc = rc
        rc = tmp
      }
      if (d) {
        if (lc !== -1) {
          res += this._size[lc]
        }
        if (rc === -1) break
        v = rc
      } else {
        if (lc === -1) break
        v = lc
      }
    }
    return res
  }

  bisectRight(int32: number): number {
    return this.bisectLeft(int32 + 1)
  }

  at(index: number): number | undefined {
    if (index < 0) index += this.size
    if (index < 0 || index >= this.size) return undefined
    let v = 0
    let res = 0
    for (let i = this._maxLog - 1; ~i; i--) {
      const left = (this._lazy >>> i) & 1
      let lc = this._edges[2 * v]
      let rc = this._edges[2 * v + 1]
      if (left) {
        const tmp = lc
        lc = rc
        rc = tmp
      }
      if (lc === -1) {
        v = rc
        res |= 1 << i
        continue
      }
      if (this._size[lc] <= index) {
        index -= this._size[lc]
        v = rc
        res |= 1 << i
      } else {
        v = lc
      }
    }
    return res
  }

  xorAll(int32: number): void {
    this._lazy ^= int32
  }

  has(int32: number): boolean {
    return !!this.count(int32)
  }

  forEach(callbackfn: (value: number, index: number) => void): void {
    let queue: [number, number][] = [[0, 0]]
    for (let i = this._maxLog - 1; ~i; i--) {
      const left = (this._lazy >>> i) & 1
      const nextQueue: [number, number][] = []
      for (let j = 0; j < queue.length; j++) {
        const { 0: v, 1: x } = queue[j]
        let lc = this._edges[2 * v]
        let rc = this._edges[2 * v + 1]
        if (left === 1) {
          const tmp = lc
          lc = rc
          rc = tmp
        }
        if (lc !== -1) {
          nextQueue.push([lc, 2 * x])
        }
        if (rc !== -1) {
          nextQueue.push([rc, 2 * x + 1])
        }
      }
      queue = nextQueue
    }

    let ptr = 0
    for (let j = 0; j < queue.length; j++) {
      const { 0: v, 1: x } = queue[j]
      for (let _ = 0; _ < this._endCount[v]; _++) {
        callbackfn(x, ptr++)
      }
    }
  }

  toString(): string {
    const res: number[] = []
    this.forEach(x => res.push(x))
    return `XorTrie { ${res.join(', ')} }`
  }

  get size(): number {
    return this._size[0]
  }

  get min(): number | undefined {
    return this.at(0)
  }

  get max(): number | undefined {
    return this.at(-1)
  }
}

export { XorTrie, XorTrie as BinaryTrie }

if (require.main === module) {
  // https://leetcode.cn/problems/maximum-xor-of-two-numbers-in-an-array/description/
  function findMaximumXOR(nums: number[]): number {
    const trie = new XorTrie({ maxInt32: Math.max(...nums), addLimit: nums.length, allowMultipleElements: true })
    let res = 0
    nums.forEach(num => {
      trie.add(num)
      trie.xorAll(num)
      res = Math.max(res, trie.max!)
      trie.xorAll(num)
    })
    return res
  }

  // 1803. 统计异或值在范围内的数对有多少
  // https://leetcode.cn/problems/count-pairs-with-xor-in-a-range/description/
  function countPairs(nums: number[], low: number, high: number): number {
    const n = nums.length
    const xorTrie = new XorTrie({ maxInt32: 1e5 + 10, addLimit: n, allowMultipleElements: true })
    nums.forEach(x => xorTrie.add(x))
    let res = 0
    for (let i = 0; i < n; i++) {
      xorTrie.xorAll(nums[i])
      res += xorTrie.bisectRight(high) - xorTrie.bisectLeft(low)
      xorTrie.xorAll(nums[i])
    }
    return res / 2
  }

  // 2935. 找出强数对的最大异或值 II
  // https://leetcode.cn/problems/maximum-strong-pair-xor-ii/description/

  function maximumStrongPairXor(nums: number[]): number {
    nums.sort((a, b) => a - b)
    const n = nums.length
    let res = 0
    let left = 0
    const trie = new XorTrie({ maxInt32: Math.max(...nums), addLimit: n, allowMultipleElements: true })
    for (let right = 0; right < n; right++) {
      trie.add(nums[right])
      while (left <= right && nums[right] > 2 * nums[left]) {
        trie.discard(nums[left])
        left++
      }
      trie.xorAll(nums[right])
      res = Math.max(res, trie.max!)
      trie.xorAll(nums[right])
    }
    return res
  }
}
