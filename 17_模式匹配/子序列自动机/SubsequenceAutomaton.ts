/* eslint-disable no-inner-declarations */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// !子序列自动机.适用于多次子序列匹配的场景.
// API:
// - move(pos, newValue) -> nextPos:
//     查询当前位置的下一个特定字符的位置(下标严格大于pos).如果不存在，则为 n. 0<=pos<n
// - includes(t, sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> bool:
//     查询s[sStart:sEnd]是否含有某序列t[tStart:tEnd].时间复杂度O(len(t)logn).
// !- !match(t, sStart=0, sEnd=-1, tStart=0, tEnd=-1) -> (hit,end):
//     在 s[sStart:sEnd] 中寻找子序列 t[tStart:tEnd].时间复杂度 O(len(t)logn).
//     适合处理需要多次匹配的场景.

/**
 * `O(∑*n)`预处理,`∑`为字符集大小.
 * `O(len(t))`查询,`len(t)`为待匹配序列的长度.
 */
class SubsequnceAutomaton1 {
  /**
   * `nexts[i*charset+j]`表示从i索引开始,下一个字符j的索引(严格大于i).
   * 如果不存在,则为n.
   */
  private readonly _nexts: Uint32Array
  private readonly _s: string
  private readonly _charset: number
  private readonly _offset: number

  /**
   * @param s 主串.
   * @param charset 字符集大小.
   * @param offset 字符集偏移量(起始字符).
   */
  constructor(s: string, charset = 26, offset = 97) {
    this._s = s
    this._charset = charset
    this._offset = offset
    this._nexts = this._build()
  }

  /**
   * 查询当前位置的下一个特定字符的位置(下标严格大于pos).
   * 如果不存在，则为 n.
   * 0<=pos<n.
   */
  move(pos: number, newValue: string): number {
    const v = newValue.charCodeAt(0) - this._offset
    return this._nexts[pos * this._charset + v]
  }

  /**
   * 查询`s[sStart:sEnd]`是否含有某序列`t[tStart:tEnd]`.
   */
  includes(
    t: string,
    options?: { sStart?: number; sEnd?: number; tStart?: number; tEnd?: number }
  ): boolean {
    const [hit] = this.match(t, options)
    let tLen: number
    if (!options) {
      tLen = t.length
    } else {
      const { tStart = 0, tEnd = t.length } = options
      tLen = tEnd - tStart
    }
    return hit >= tLen
  }

  /**
   * 在`s[sStart:sEnd]`中寻找子序列`t[tStart:tEnd]`.
   * @returns `(hit,end)`: (`匹配到的的t的长度`, `匹配结束时s的索引`)
   * 此时,匹配结束时t的索引为`tStart+hit`.
   * 耗去的s的长度为`end-sStart`.
   */
  match(
    t: string,
    options?: { sStart?: number; sEnd?: number; tStart?: number; tEnd?: number }
  ): [hit: number, end: number] {
    if (!options) options = {}
    const { sStart = 0, sEnd = this._s.length, tStart = 0, tEnd = t.length } = options
    if (sStart >= sEnd) return [0, sStart]
    if (tStart >= tEnd) return [0, sStart]
    const n = this._s.length
    let si = sStart
    let ti = tStart
    if (this._s[si] === t[ti]) ti++ // !注意需要先判断第一个字符
    while (si < sEnd && ti < tEnd) {
      const nextPos = this.move(si, t[ti])
      if (nextPos === n) return [ti - tStart, si + 1]
      si = nextPos
      ti++
    }
    return [ti - tStart, si + 1]
  }

  private _build(): Uint32Array {
    const row = this._s.length
    const col = this._charset
    const nexts = new Uint32Array(row * col)
    nexts.fill(row, -col) // 最后一行全为row
    for (let i = row - 2; ~i; i--) {
      nexts.copyWithin(i * col, (i + 1) * col, (i + 2) * col)
      const v = this._s.charCodeAt(i + 1) - this._offset
      nexts[i * col + v] = i + 1
    }
    return nexts
  }
}

/**
 * `O(n)`预处理.
 * `O(len(t)logn)`查询,`len(t)`为待匹配序列的长度.
 * !复杂度与字符种类数无关.占用内存更小.
 */
class SubsequnceAutomaton2<V> {
  private readonly _arr: ArrayLike<V>
  private readonly _indexes: Map<V, number[]>

  constructor(arr: ArrayLike<V>) {
    this._arr = arr
    this._indexes = this._build()
  }

  /**
   * 查询当前位置的下一个特定字符的位置(下标严格大于pos).
   * 如果不存在，则为 n.
   * 0<=pos<n.
   */
  move(pos: number, newValue: V): number {
    const indexes = this._indexes.get(newValue)
    if (!indexes) return this._arr.length
    const nextPos = SubsequnceAutomaton2._bisectRight(indexes, pos)
    return nextPos < indexes.length ? indexes[nextPos] : this._arr.length
  }

  /**
   * 查询`s[sStart:sEnd]`是否含有某序列`t[tStart:tEnd]`.
   */
  includes(
    t: ArrayLike<V>,
    options?: { sStart?: number; sEnd?: number; tStart?: number; tEnd?: number }
  ): boolean {
    const [hit] = this.match(t, options)
    let tLen: number
    if (!options) {
      tLen = t.length
    } else {
      const { tStart = 0, tEnd = t.length } = options
      tLen = tEnd - tStart
    }
    return hit >= tLen
  }

  /**
   * 在`s[sStart:sEnd]`中寻找子序列`t[tStart:tEnd]`.
   * @returns `(hit,end)`: (`匹配到的的t的长度`, `匹配结束时s的索引`)
   * 此时,匹配结束时t的索引为`tStart+hit`.
   * 耗去的s的长度为`end-sStart`.
   */
  match(
    t: ArrayLike<V>,
    options?: { sStart?: number; sEnd?: number; tStart?: number; tEnd?: number }
  ): [hit: number, end: number] {
    if (!options) options = {}
    const { sStart = 0, sEnd = this._arr.length, tStart = 0, tEnd = t.length } = options
    if (sStart >= sEnd) return [0, sStart]
    if (tStart >= tEnd) return [0, sStart]
    const n = this._arr.length
    let si = sStart
    let ti = tStart
    if (this._arr[si] === t[ti]) ti++ // !注意需要先判断第一个字符
    while (si < sEnd && ti < tEnd) {
      const nextPos = this.move(si, t[ti])
      if (nextPos === n) return [ti - tStart, si + 1]
      si = nextPos
      ti++
    }
    return [ti - tStart, si + 1]
  }

  private _build(): Map<V, number[]> {
    const indexes = new Map<V, number[]>()
    for (let i = 0; i < this._arr.length; i++) {
      const v = this._arr[i]
      if (!indexes.has(v)) indexes.set(v, [])
      indexes.get(v)!.push(i)
    }
    return indexes
  }

  private static _bisectRight<V>(arr: ArrayLike<V>, value: V): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >>> 1
      if (arr[mid] <= value) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }
}

export { SubsequnceAutomaton1, SubsequnceAutomaton2 }

if (require.main === module) {
  // 792. 匹配子序列的单词数
  // https://leetcode.cn/problems/number-of-matching-subsequences/
  function numMatchingSubseq(s: string, words: string[]): number {
    const S = new SubsequnceAutomaton1(s)
    return words.filter(w => S.includes(w)).length
  }

  // 727. 最小窗口子序列
  // https://leetcode.cn/problems/minimum-window-subsequence/
  function minWindow(s1: string, s2: string): string {
    const S = new SubsequnceAutomaton1(s1)
    const starts: number[] = []
    for (let i = 0; i < s1.length; i++) if (s1[i] === s2[0]) starts.push(i)

    let res: [number, number] | null = null
    starts.forEach(sStart => {
      const [hit, sEnd] = S.match(s2, { sStart })
      if (hit !== s2.length) return
      const sLen = sEnd - sStart
      if (!res || sLen < res[1] - res[0]) res = [sStart, sEnd]
    })

    return res ? s1.slice(res[0], res[1]) : ''
  }
}
