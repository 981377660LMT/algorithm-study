/* eslint-disable no-inner-declarations */

/** 数组实现的AC自动机. */
class ACAutoMatonArray {
  /** wordPos[i] 表示加入的第i个模式串对应的节点编号. */
  readonly wordPos: number[] = []

  private readonly _sigma: number
  private readonly _offset: number
  private _children: Int32Array

  /** 又叫fail.指向当前节点最长真后缀对应结点，例如"bc"是"abc"的最长真后缀. */
  private _suffixLink!: Int32Array

  private _bfsOrder!: Int32Array
  private _needUpdateChildren!: boolean
  private _nodeCount = 0

  constructor(options: { lengthSum: number; sigma?: number; offset?: number }) {
    const { lengthSum, sigma = 26, offset = 97 } = options
    this._sigma = sigma
    this._offset = offset
    this._children = new Int32Array((1 + lengthSum) * sigma)
    this._newNode()
  }

  addString(str: string): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      const ord = str.charCodeAt(i) - this._offset
      const hash = pos * this._sigma + ord
      if (this._children[hash] === -1) {
        const id = this._newNode()
        this._children[hash] = id
      }
      pos = this._children[hash]
    }
    this.wordPos.push(pos)
    return pos
  }

  addChar(pos: number, ord: number): number {
    ord -= this._offset
    const hash = pos * this._sigma + ord
    if (this._children[hash] === -1) {
      const id = this._newNode()
      this._children[hash] = id
    }
    return this._children[hash]
  }

  move(pos: number, ord: number): number {
    ord -= this._offset
    if (this._needUpdateChildren) {
      return this._children[pos * this._sigma + ord]
    }
    while (true) {
      const hash = pos * this._sigma + ord
      if (this._children[hash] !== -1) return this._children[hash]
      if (pos === 0) return 0
      pos = this._suffixLink[pos]
    }
  }

  /**
   * @param needUpdateChildren move调用较少时，设置为false更快.
   */
  buildSuffixLink(needUpdateChildren = false) {
    this._needUpdateChildren = needUpdateChildren
    this._suffixLink = new Int32Array(this._nodeCount).fill(-1)
    this._bfsOrder = new Int32Array(this._nodeCount)
    let head = 0
    let tail = 1
    while (head < tail) {
      let v = this._bfsOrder[head]
      head++
      const offset = v * this._sigma
      for (let i = 0; i < this._sigma; i++) {
        const next = this._children[offset + i]
        if (next === -1) {
          continue
        }
        this._bfsOrder[tail++] = next
        let f = this._suffixLink[v]
        while (f !== -1 && this._children[f * this._sigma + i] === -1) {
          f = this._suffixLink[f]
        }
        this._suffixLink[next] = f
        if (f === -1) {
          this._suffixLink[next] = 0
        } else {
          this._suffixLink[next] = this._children[f * this._sigma + i]
        }
      }
    }
    if (!needUpdateChildren) return
    for (let i = 0; i < this._nodeCount; i++) {
      const v = this._bfsOrder[i]
      const offset = v * this._sigma
      for (let j = 0; j < this._sigma; j++) {
        const next = this._children[offset + j]
        if (next === -1) {
          const f = this._suffixLink[v]
          if (f === -1) {
            this._children[offset + j] = 0
          } else {
            this._children[offset + j] = this._children[f * this._sigma + j]
          }
        }
      }
    }
  }

  /** 获取每个状态匹配到的模式串的个数. */
  getCounter(): Uint32Array {
    const counter = new Uint32Array(this._nodeCount)
    for (let i = 0; i < this.wordPos.length; i++) {
      counter[this.wordPos[i]]++
    }
    for (let i = 0; i < this._bfsOrder.length; i++) {
      const v = this._bfsOrder[i]
      if (v !== 0) {
        counter[v] += counter[this._suffixLink[v]]
      }
    }
    return counter
  }

  /** 获取每个状态匹配到的模式串的索引. */
  getIndexes(): number[][] {
    const res: number[][] = Array(this._nodeCount)
    for (let i = 0; i < res.length; i++) res[i] = []
    for (let i = 0; i < this.wordPos.length; i++) res[this.wordPos[i]].push(i)
    for (let i = 0; i < this._bfsOrder.length; i++) {
      const v = this._bfsOrder[i]
      if (v !== 0) {
        const from = this._suffixLink[v]
        const arr1 = res[from]
        const arr2 = res[v]
        const arr3 = Array(arr1.length + arr2.length)
        let p1 = 0
        let p2 = 0
        let p3 = 0
        while (p1 < arr1.length && p2 < arr2.length) {
          if (arr1[p1] < arr2[p2]) {
            arr3[p3++] = arr1[p1++]
          } else if (arr1[p1] > arr2[p2]) {
            arr3[p3++] = arr2[p2++]
          } else {
            arr3[p3++] = arr1[p1]
            p1++
            p2++
          }
        }
        while (p1 < arr1.length) arr3[p3++] = arr1[p1++]
        while (p2 < arr2.length) arr3[p3++] = arr2[p2++]
        res[v] = arr3
      }
    }
    return res
  }

  dp(f: (from: number, to: number) => void): void {
    for (let i = 0; i < this._bfsOrder.length; i++) {
      const v = this._bfsOrder[i]
      if (v !== 0) {
        f(this._suffixLink[v], v)
      }
    }
  }

  buildFailTree(): number[][] {
    const res: number[][] = Array(this.size)
    for (let i = 0; i < res.length; i++) res[i] = []
    this.dp((pre, cur) => {
      res[pre].push(cur)
    })
    return res
  }

  empty(): boolean {
    return this._nodeCount === 1
  }

  /**
   * 返回str在trie树上的节点位置.如果不存在，返回0.
   */
  search(str: string): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      const ord = str.charCodeAt(i) - this._offset
      const hash = pos * this._sigma + ord
      if (hash < 0 || hash >= this._children.length || this._children[hash] === -1) {
        return 0
      }
      pos = this._children[hash]
    }
    return pos
  }

  get size(): number {
    return this._nodeCount
  }

  private _newNode(): number {
    const start = this._nodeCount * this._sigma
    const end = start + this._sigma
    if (end > this._children.length) {
      const tmp = new Int32Array(this._children.length * 2 + this._sigma)
      tmp.set(this._children)
      this._children = tmp
    }
    this._children.fill(-1, start, end)
    this._nodeCount++
    return this._nodeCount - 1
  }
}

export { ACAutoMatonArray }

if (require.main === module) {
  // https://leetcode.cn/problems/length-of-the-longest-valid-substring/description/
  function longestValidSubstring(word: string, forbidden: string[]): number {
    const lengthSum = forbidden.reduce((sum, w) => sum + w.length, 0)
    const acm = new ACAutoMatonArray({ lengthSum })
    forbidden.forEach(w => acm.addString(w))
    acm.buildSuffixLink(false)

    const minBorder = new Uint32Array(acm.size).fill(-1)
    for (let i = 0; i < forbidden.length; i++) {
      minBorder[acm.wordPos[i]] = Math.min(minBorder[acm.wordPos[i]], forbidden[i].length)
    }
    acm.dp((from, to) => {
      minBorder[to] = Math.min(minBorder[to], minBorder[from])
    })

    let res = 0
    let left = 0
    let pos = 0
    for (let right = 0; right < word.length; right++) {
      pos = acm.move(pos, word.charCodeAt(right))
      left = Math.max(left, right - minBorder[pos] + 2)
      res = Math.max(res, right - left + 1)
    }
    return res
  }

  // 1032. 字符流
  // https://leetcode.cn/problems/stream-of-characters/description/
  class StreamChecker {
    private readonly _acm: ACAutoMatonArray
    private readonly _counter: Uint32Array
    private _pos = 0

    constructor(words: string[]) {
      const acm = new ACAutoMatonArray({ lengthSum: words.reduce((sum, w) => sum + w.length, 0) })
      words.forEach(w => acm.addString(w))
      acm.buildSuffixLink(false)
      this._acm = acm
      this._counter = acm.getCounter()
    }

    query(letter: string): boolean {
      this._pos = this._acm.move(this._pos, letter.charCodeAt(0))
      return this._counter[this._pos] > 0
    }
  }

  /**
   * Your StreamChecker object will be instantiated and called as such:
   * var obj = new StreamChecker(words)
   * var param_1 = obj.query(letter)
   */
}
