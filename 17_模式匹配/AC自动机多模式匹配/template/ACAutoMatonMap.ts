/* eslint-disable no-inner-declarations */
/* eslint-disable no-loop-func */
/* eslint-disable max-len */

class ACAutoMatonMap<T = string> {
  /** wordPos[i] 表示加入的第i个模式串对应的节点编号. */
  readonly wordPos: number[] = []
  private _children: Map<T, number>[] = [new Map()]
  private _link!: Int32Array
  private _linkWord?: Int32Array
  private _bfsOrder!: Int32Array

  addString(str: ArrayLike<T>): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      const ord = str[i]
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

  addChar(pos: number, ord: T): number {
    let nexts = this._children[pos]
    if (nexts.has(ord)) {
      return nexts.get(ord)!
    }
    const nextState = this._children.length
    nexts.set(ord, nextState)
    this._children.push(new Map())
    return nextState
  }

  move(pos: number, ord: T): number {
    while (true) {
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        return nexts.get(ord)!
      }
      if (pos === 0) {
        return 0
      }
      pos = this._link[pos]
    }
  }

  buildSuffixLink() {
    this._link = new Int32Array(this._children.length).fill(-1)
    this._bfsOrder = new Int32Array(this._children.length)
    let head = 0
    let tail = 1
    while (head < tail) {
      const v = this._bfsOrder[head]
      head++
      this._children[v].forEach((next, char) => {
        this._bfsOrder[tail] = next
        tail++
        let f = this._link[v]
        while (f !== -1) {
          if (this._children[f].has(char)) {
            break
          }
          f = this._link[f]
        }
        if (f === -1) {
          this._link[next] = 0
        } else {
          this._link[next] = this._children[f].get(char)!
        }
      })
    }
  }

  linkWord(pos: number): number {
    if (this._linkWord) return this._linkWord[pos]
    const size = this.size
    this._linkWord = new Int32Array(size)
    const hasWord = new Uint8Array(size)
    for (let i = 0; i < this.wordPos.length; i++) hasWord[this.wordPos[i]] = 1
    const link = this._link
    const linkWord = this._linkWord
    for (let i = 0; i < this._bfsOrder.length; i++) {
      const v = this._bfsOrder[i]
      if (v !== 0) {
        const p = link[v]
        linkWord[v] = hasWord[p] ? p : linkWord[p]
      }
    }
    return this._linkWord[pos]
  }

  getCounter(): Uint32Array {
    const counter = new Uint32Array(this._children.length)
    this.wordPos.forEach(pos => {
      counter[pos]++
    })
    this._bfsOrder.forEach(v => {
      if (v !== 0) {
        counter[v] += counter[this._link[v]]
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
        const from = this._link[v]
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
        f(this._link[v], v)
      }
    })
  }

  buildFailTree(): number[][] {
    const res: number[][] = Array(this.size)
    for (let i = 0; i < res.length; i++) res[i] = []
    this.dp((pre, cur) => {
      res[pre].push(cur)
    })
    return res
  }

  buildTrieTree(): number[][] {
    const res: number[][] = Array(this.size)
    for (let i = 0; i < res.length; i++) res[i] = []
    const dfs = (cur: number): void => {
      this._children[cur].forEach((next, _) => {
        res[cur].push(next)
        dfs(next)
      })
    }
    dfs(0)
    return res
  }

  /**
   * 返回str在trie树上的节点位置.如果不存在，返回0.
   */
  search(str: ArrayLike<T>): number {
    if (str.length === 0) return 0
    let pos = 0
    for (let i = 0; i < str.length; i++) {
      if (pos < 0 || pos >= this._children.length) {
        return 0
      }
      const ord = str[i]
      const nexts = this._children[pos]
      if (nexts.has(ord)) {
        pos = nexts.get(ord)!
      } else {
        return 0
      }
    }
    return pos
  }

  empty(): boolean {
    return this._children.length === 1
  }

  clear(): void {
    this.wordPos.length = 0
    this._children = [new Map()]
    this._link = new Int32Array(0)
    this._linkWord = undefined
    this._bfsOrder = new Int32Array(0)
  }

  get size(): number {
    return this._children.length
  }
}

export { ACAutoMatonMap }

if (require.main === module) {
  const INF = 2e15

  // 100350. 最小代价构造字符串
  // https://leetcode.cn/problems/construct-string-with-minimum-cost/description/
  function minimumCost(target: string, words: string[], costs: number[]): number {
    const acm = new ACAutoMatonMap()
    for (let i = 0; i < words.length; i++) {
      acm.addString(words[i])
    }
    acm.buildSuffixLink()

    const nodeCosts = new Uint32Array(acm.size).fill(-1)
    const nodeDepth = new Uint32Array(acm.size)
    for (let i = 0; i < acm.wordPos.length; i++) {
      const pos = acm.wordPos[i]
      nodeCosts[pos] = Math.min(nodeCosts[pos], costs[i])
      nodeDepth[pos] = words[i].length
    }

    const dp = Array(target.length + 1).fill(INF)
    dp[0] = 0
    let pos = 0
    for (let i = 0; i < target.length; i++) {
      pos = acm.move(pos, target[i])
      for (let cur = pos; cur !== 0; cur = acm.linkWord(cur)) {
        dp[i + 1] = Math.min(dp[i + 1], dp[i + 1 - nodeDepth[cur]] + nodeCosts[cur])
      }
    }
    return dp[target.length] === INF ? -1 : dp[target.length]
  }

  // https://leetcode.cn/problems/length-of-the-longest-valid-substring/description/
  function longestValidSubstring(word: string, forbidden: string[]): number {
    const acm = new ACAutoMatonMap()
    forbidden.forEach(w => acm.addString(w))
    acm.buildSuffixLink()

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
      pos = acm.move(pos, word[right])
      left = Math.max(left, right - minBorder[pos] + 2)
      res = Math.max(res, right - left + 1)
    }
    return res
  }

  // 1032. 字符流
  // https://leetcode.cn/problems/stream-of-characters/description/
  class StreamChecker {
    private readonly _acm: ACAutoMatonMap<number>
    private readonly _counter: Uint32Array
    private _pos = 0

    constructor(words: string[]) {
      const acm = new ACAutoMatonMap<number>()
      words.forEach(w => acm.addString(w.split('').map(c => c.charCodeAt(0))))
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
