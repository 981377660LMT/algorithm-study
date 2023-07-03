/* eslint-disable no-inner-declarations */

/**
 * `O(nm/32)`求最长公共子序列.
 * 5e4x5e4 => 1s.
 */
function getLCS<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): T[] {
  const n = arr1.length
  const m = arr2.length
  const copy1 = new Uint32Array(n)
  const copy2 = new Uint32Array(m)

  const id = new Map<T, number>()
  for (let i = 0; i < n; i++) {
    const item = arr1[i]
    const curId = id.get(item)
    if (curId !== void 0) {
      copy1[i] = curId
    } else {
      copy1[i] = id.size
      id.set(item, id.size)
    }
  }
  for (let i = 0; i < m; i++) {
    const item = arr2[i]
    const curId = id.get(item)
    if (curId !== void 0) {
      copy2[i] = curId
    } else {
      copy2[i] = id.size
      id.set(item, id.size)
    }
  }
  const rid: T[] = Array(id.size)
  id.forEach((v, k) => {
    rid[v] = k
  })

  const sets: _BS[] = Array(id.size)
  for (let i = 0; i < id.size; ++i) sets[i] = new _BS(n + 1)
  const dp: _BS[] = Array(m + 1)
  for (let i = 0; i <= m; ++i) dp[i] = new _BS(n + 1)
  let tmp = new _BS(n + 1)
  for (let i = 0; i < n; ++i) sets[copy1[i]].set(i + 1)
  for (let i = 1; i <= m; ++i) {
    dp[i].setAll(dp[i - 1])
    tmp.setAll(dp[i])
    tmp.orAll(sets[copy2[i - 1]])
    dp[i].shift()
    dp[i].set(0)
    dp[i].setAll(_BS.minus(tmp, dp[i]))
    dp[i].xorAll(tmp)
    dp[i].andAll(tmp)
  }

  let i = n
  let j = m
  const res: T[] = []
  while (i > 0 && j > 0) {
    if (copy1[i - 1] === copy2[j - 1]) {
      res.push(rid[copy1[i - 1]])
      i--
      j--
    } else if (!dp[j].get(i)) {
      i--
    } else {
      j--
    }
  }
  return res.reverse()
}

class _BS {
  static minus(a: _BS, b: _BS): _BS {
    let last = 0
    const res = new _BS(a._data)
    for (let i = 0; i < a._data.length; ++i) {
      const cur = a._data[i] < b._data[i] + last
      res._data[i] = a._data[i] - b._data[i] - last
      last = +cur
    }
    return res
  }

  private readonly _data: Uint32Array

  constructor(nOrData: number | Uint32Array) {
    if (typeof nOrData === 'number') {
      this._data = new Uint32Array((nOrData >>> 5) + 1)
    } else {
      this._data = nOrData.slice()
    }
  }

  set(i: number): void {
    this._data[i >>> 5] |= 1 << (i & 31)
  }

  get(i: number): boolean {
    return !!(this._data[i >>> 5] & (1 << (i & 31)))
  }

  shift(): void {
    let last = 0
    for (let i = 0; i < this._data.length; ++i) {
      const cur = this._data[i] >>> 31
      this._data[i] <<= 1
      this._data[i] |= last
      last = cur
    }
  }

  setAll(bs: _BS): void {
    this._data.set(bs._data)
  }

  orAll(bs: _BS): void {
    for (let i = 0; i < this._data.length; ++i) {
      this._data[i] |= bs._data[i]
    }
  }

  xorAll(bs: _BS): void {
    for (let i = 0; i < this._data.length; ++i) {
      this._data[i] ^= bs._data[i]
    }
  }

  andAll(bs: _BS): void {
    for (let i = 0; i < this._data.length; ++i) {
      this._data[i] &= bs._data[i]
    }
  }
}

export { getLCS }

if (require.main === module) {
  // https://leetcode.cn/problems/longest-common-subsequence/
  function longestCommonSubsequence(text1: string, text2: string): number {
    return getLCS(text1, text2).length
  }

  const n1 = 5e4
  const n2 = 5e4
  const text1 = Array(n1)
    .fill(0)
    .map(() => String.fromCharCode(Math.floor(Math.random() * 26) + 97))
    .join('')
  const text2 = Array(n2)
    .fill(0)
    .map(() => String.fromCharCode(Math.floor(Math.random() * 26) + 97))
    .join('')
  console.time('getLCS')
  console.log(longestCommonSubsequence(text1, text2))
  console.timeEnd('getLCS')
}
