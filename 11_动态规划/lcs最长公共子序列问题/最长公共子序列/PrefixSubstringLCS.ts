/**
 * 查询s1[:a)和s2[b:c)的最长公共子序列长度.
 */
class PrefixSubstringLCS<S = number> {
  private readonly _col: number
  private readonly _dp1: Uint32Array

  constructor(arr1: ArrayLike<S>, arr2: ArrayLike<S>) {
    const row = arr1.length + 1
    const col = arr2.length + 1
    const dp1 = new Uint32Array(row * col)
    const dp2 = new Uint32Array(row * col)
    for (let c = 0; c < col; ++c) dp1[c] = c
    for (let r = 1; r < row; ++r) {
      for (let c = 1; c < col; ++c) {
        const pos1 = r * col + c
        const pos2 = r * col + c - 1
        const pos3 = (r - 1) * col + c
        if (arr1[r - 1] === arr2[c - 1]) {
          dp1[pos1] = dp2[pos2]
          dp2[pos1] = dp1[pos3]
        } else {
          let max = dp1[pos3]
          let min = dp2[pos2]
          if (max < min) {
            max ^= min
            min ^= max
            max ^= min
          }
          dp1[pos1] = max
          dp2[pos1] = min
        }
      }
    }

    this._col = col
    this._dp1 = dp1
  }

  /**
   * 查询 `arr1[:a)` 和 `arr2[b:c)` 的最长公共子序列长度.
   * @complexity `O(c - b)`
   */
  query(a: number, b: number, c: number): number {
    let res = 0
    const offset = a * this._col
    for (let i = b + 1; i < c + 1; ++i) {
      res += +(this._dp1[offset + i] <= b)
    }
    return res
  }
}

export { PrefixSubstringLCS }

if (require.main === module) {
  const s1 = 'abcb'
  const s2 = 'acb'
  const LCS = new PrefixSubstringLCS(s1, s2)
  console.log(LCS.query(3, 0, 3))
  console.log(LCS.query(4, 1, 3))
  console.log(LCS.query(2, 1, 2))

  // test 4000*4000 and 2e5 queries => 1s
  const n = 4000
  const q = 2e5
  const arr1 = Array.from({ length: n }, () => Math.floor(Math.random() * 1e9))
  const arr2 = Array.from({ length: n }, () => Math.floor(Math.random() * 1e9))
  const LCS2 = new PrefixSubstringLCS(arr1, arr2)
  const queries = Array.from({ length: q }, () => [
    n - Math.floor(Math.random()) * 100,
    Math.floor(Math.random()) * 100,
    n - Math.floor(Math.random()) * 100
  ])

  console.time('test')
  queries.forEach(([a, b, c]) => LCS2.query(a, b, c))
  console.timeEnd('test')
}
