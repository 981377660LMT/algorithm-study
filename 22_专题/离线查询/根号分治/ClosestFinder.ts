/* eslint-disable @typescript-eslint/no-non-null-assertion */

// https://leetcode.cn/problems/find-closest-lcci/
// https://www.luogu.com.cn/blog/236866/sol-p5397 的无修改版

/**
 * 给定一个数组，每次查询给定`x`和`y`，查询数组中`x`和`y`的最近距离.
 * `O(nsqrt(n))`预处理，`O(sqrt(n))`查询.
 */
class ClosestFinder<T> {
  private readonly _threshold: number
  private readonly _ids: number[]
  private readonly _mp: number[][] = []
  private readonly _largeResRecord: Map<number, Uint32Array> = new Map()
  private readonly _valueToId = new Map<T, number>()

  /**
   * @threshold 当元素出现的次数超过这个阈值时，会对此元素进行时空复杂度`O(n)`的预处理.
   */
  constructor(arr: ArrayLike<T>, threshold = 4 * (1 + (Math.sqrt(arr.length) | 0))) {
    this._threshold = threshold
    const ids = Array(arr.length)
    for (let i = 0; i < arr.length; i++) {
      const id = this._getId(arr[i])
      ids[i] = id
      if (this._mp.length <= id) this._mp.push([])
      this._mp[id].push(i)
    }
    this._ids = ids
    this._mp.forEach((_, id) => {
      if (this._isLarge(id)) {
        this._largeResRecord.set(id, this._buildLargeRes(id))
      }
    })
  }

  /**
   * 查询`x`和`y`的最近距离.
   * 如果`x`和`y`相同，返回`0`.
   * 如果`x`或`y`不存在，返回`undefined`.
   */
  findClosest(x: T, y: T): number | undefined {
    if (x === y) return 0
    if (!this._valueToId.has(x) || !this._valueToId.has(y)) return void 0
    const id1 = this._getId(x)
    const id2 = this._getId(y)
    if (this._isLarge(id1)) return this._largeResRecord.get(id1)![id2]
    if (this._isLarge(id2)) return this._largeResRecord.get(id2)![id1]
    const pos1 = this._mp[id1]
    const pos2 = this._mp[id2]
    let i = 0
    let j = 0
    let res = this._ids.length
    while (i < pos1.length && j < pos2.length) {
      res = Math.min(res, Math.abs(pos1[i] - pos2[j]))
      if (pos1[i] < pos2[j]) i++
      else j++
    }
    return res
  }

  private _getId(value: T): number {
    const res = this._valueToId.get(value)
    if (res !== void 0) return res
    const id = this._valueToId.size
    this._valueToId.set(value, id)
    return id
  }

  private _isLarge(id: number): boolean {
    return this._mp[id].length > this._threshold
  }

  private _buildLargeRes(id: number): Uint32Array {
    const n = this._ids.length
    const res = new Uint32Array(this._valueToId.size).fill(n)
    let dist = n
    for (let i = 0; i < n; i++) {
      const cur = this._ids[i]
      if (cur === id) {
        dist = 0
      } else {
        dist++
        res[cur] = Math.min(res[cur], dist)
      }
    }
    dist = n
    for (let i = n - 1; i >= 0; i--) {
      const cur = this._ids[i]
      if (cur === id) {
        dist = 0
      } else {
        dist++
        res[cur] = Math.min(res[cur], dist)
      }
    }
    return res
  }
}

export { ClosestFinder }

if (require.main === module) {
  // 面试题 17.11. 单词距离
  // https://leetcode.cn/problems/find-closest-lcci/
  // eslint-disable-next-line no-inner-declarations
  function findClosest(words: string[], word1: string, word2: string): number {
    const finder = new ClosestFinder(words)
    return finder.findClosest(word1, word2)!
  }

  // 244. 最短单词距离 II
  // https://leetcode.cn/problems/shortest-word-distance-ii/description/
  class WordDistance {
    private readonly _finder: ClosestFinder<string>
    constructor(wordsDict: string[]) {
      this._finder = new ClosestFinder(wordsDict)
    }

    shortest(word1: string, word2: string): number {
      return this._finder.findClosest(word1, word2)!
    }
  }

  /**
   * Your WordDistance object will be instantiated and called as such:
   * var obj = new WordDistance(wordsDict)
   * var param_1 = obj.shortest(word1,word2)
   */

  // test time
  const n = 1e6
  const bigArr = Array(n)
  for (let i = 0; i < bigArr.length; i++) {
    bigArr[i] = (Math.random() * 100) | 0
  }
  const finder = new ClosestFinder(bigArr)
  const randomPairs = Array(n)
  for (let i = 0; i < randomPairs.length; i++) {
    const pos1 = (Math.random() * 100) | 0
    const pos2 = (Math.random() * 100) | 0
    randomPairs[i] = [bigArr[pos1], bigArr[pos2]]
  }

  console.time('build')
  for (let i = 0; i < randomPairs.length; i++) {
    finder.findClosest(randomPairs[i][0], randomPairs[i][1])
  }
  console.timeEnd('build')
}
