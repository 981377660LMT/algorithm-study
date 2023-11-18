/* eslint-disable no-inner-declarations */
/* eslint-disable no-loop-func */
/* eslint-disable max-len */

class ACAutoMatonMap {
  /** wordPos[i] 表示加入的第i个模式串对应的节点编号. */
  readonly wordPos: number[] = [] //
  private readonly _children: Map<number, number>[] = [new Map()]
  private _suffixLink!: Int32Array
  private _bfsOrder!: Int32Array

  addString(str: string): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      const ord = str[i].charCodeAt(0)
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        pos = nexts.get(ord)!
      } else {
        const nextState = this._children.length
        nexts.set(ord, nextState)
        pos = nextState
        this._children.push(new Map())
      }
    }
    this.wordPos.push(pos)
    return pos
  }

  addChar(pos: number, ord: number): number {
    let nexts = this._children[pos]
    if (nexts.has(ord)) {
      return nexts.get(ord)!
    }
    const nextState = this._children.length
    nexts.set(ord, nextState)
    this._children.push(new Map())
    return nextState
  }

  move(pos: number, ord: number): number {
    while (true) {
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        return nexts.get(ord)!
      }
      if (pos === 0) {
        return 0
      }
      pos = this._suffixLink[pos]
    }
  }

  buildSuffixLink() {
    this._suffixLink = new Int32Array(this._children.length).fill(-1)
    this._bfsOrder = new Int32Array(this._children.length)
    let head = 0
    let tail = 1
    while (head < tail) {
      const v = this._bfsOrder[head]
      head++
      this._children[v].forEach((next, char) => {
        this._bfsOrder[tail] = next
        tail++
        let f = this._suffixLink[v]
        while (f !== -1) {
          if (this._children[f].has(char)) {
            break
          }
          f = this._suffixLink[f]
        }
        if (f === -1) {
          this._suffixLink[next] = 0
        } else {
          this._suffixLink[next] = this._children[f].get(char)!
        }
      })
    }
  }

  getCounter(): Uint32Array {
    const counter = new Uint32Array(this._children.length)
    this.wordPos.forEach(pos => {
      counter[pos]++
    })
    this._bfsOrder.forEach(v => {
      if (v !== 0) {
        counter[v] += counter[this._suffixLink[v]]
      }
    })
    return counter
  }

  getIndexes(): number[][] {
    let res: number[][] = Array(this._children.length)
    for (let i = 0; i < res.length; i++) res[i] = []
    for (let i = 0; i < this.wordPos.length; i++) {
      let pos = this.wordPos[i]
      res[pos].push(i)
    }
    this._bfsOrder.forEach(v => {
      if (v !== 0) {
        const from = this._suffixLink[v]
        const arr1 = res[from]
        const arr2 = res[v]
        const arr3 = Array(arr1.length + arr2.length)
        let i = 0
        let j = 0
        let k = 0
        while (i < arr1.length && j < arr2.length) {
          if (arr1[i] < arr2[j]) {
            arr3[k++] = arr1[i++]
          } else if (arr1[i] > arr2[j]) {
            arr3[k++] = arr2[j++]
          } else {
            arr3[k++] = arr1[i]
            i++
            j++
          }
        }
        while (i < arr1.length) arr3[k++] = arr1[i++]
        while (j < arr2.length) arr3[k++] = arr2[j++]
        res[v] = arr3
      }
    })
    return res
  }

  dp(f: (from: number, to: number) => void): void {
    this._bfsOrder.forEach(v => {
      if (v !== 0) {
        f(this._suffixLink[v], v)
      }
    })
  }

  get size(): number {
    return this._children.length
  }
}

export { ACAutoMatonMap }

if (require.main === module) {
  // https://leetcode.cn/problems/length-of-the-longest-valid-substring/description/
  function longestValidSubstring(word: string, forbidden: string[]): number {
    const acm = new ACAutoMatonMap()
    forbidden.forEach(w => acm.addString(w))
    acm.buildSuffixLink()

    const minLen = new Uint32Array(acm.size).fill(-1)
    for (let i = 0; i < forbidden.length; i++) {
      minLen[acm.wordPos[i]] = Math.min(minLen[acm.wordPos[i]], forbidden[i].length)
    }
    acm.dp((from, to) => {
      minLen[to] = Math.min(minLen[to], minLen[from])
    })

    let res = 0
    let left = 0
    let pos = 0
    for (let right = 0; right < word.length; right++) {
      pos = acm.move(pos, word.charCodeAt(right))
      left = Math.max(left, right - minLen[pos] + 2)
      res = Math.max(res, right - left + 1)
    }
    return res
  }

  // 1032. 字符流
  // https://leetcode.cn/problems/stream-of-characters/description/
  class StreamChecker {
    private readonly _acm: ACAutoMatonMap
    private readonly _counter: Uint32Array
    private _pos = 0

    constructor(words: string[]) {
      const acm = new ACAutoMatonMap()
      words.forEach(w => acm.addString(w))
      acm.buildSuffixLink()
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
