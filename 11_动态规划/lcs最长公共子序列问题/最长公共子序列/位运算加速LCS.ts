/* eslint-disable no-inner-declarations */
// 位运算加速LCS(最长公共子序列), 这也是目前跑的最快的版本
// https://wenku.baidu.com/view/ed99e4f77c1cfad6195fa776.html?_wkts_=1688318339264
// https://atcoder.jp/contests/dp/tasks/dp_f
// https://www.cnblogs.com/-Wallace-/p/bit-lcs.html
// https://rsk0315.hatenablog.com/entry/2022/12/30/180216
// https://github.com/hqztrue/LeetCodeSolutions/blob/master/0501-0600/516.%20Longest%20Palindromic%20Subsequence.pdf
// https://github.com/hqztrue/LeetCodeSolutions/blob/master/1101-1200/1143.%20Longest%20Common%20Subsequence%200ms.cpp
// https://loj.ac/s?problemDisplayId=6564&status=Accepted

/**
 * 位运算加速LCS(最长公共子序列).
 * @complex `O(nm/32)`.1e5x1e5 => 1.5s.
 */
function LCSBit<T>(arr1: ArrayLike<T>, arr2: ArrayLike<T>): number {
  const n = arr1.length
  const m = arr2.length
  const copy1 = new Uint32Array(n)
  const copy2 = new Uint32Array(m)

  const id = new Map<T, number>()
  for (let i = 0; i < n; ++i) {
    const item = arr1[i]
    const curId = id.get(item)
    if (curId !== void 0) {
      copy1[i] = curId
    } else {
      copy1[i] = id.size
      id.set(item, id.size)
    }
  }
  for (let i = 0; i < m; ++i) {
    const item = arr2[i]
    const curId = id.get(item)
    if (curId !== void 0) {
      copy2[i] = curId
    } else {
      copy2[i] = id.size
      id.set(item, id.size)
    }
  }

  const f = Array<_BS>(id.size)
  for (let i = 0; i < id.size; ++i) f[i] = new _BS(n)
  for (let i = 0; i < n; ++i) f[copy1[i]].set(i)
  const dp = new _BS(n)
  for (let i = 0; i < m; ++i) dp.run(f[copy2[i]])
  return dp.count()
}

class _BS {
  private static _onesCount32(uint32: number): number {
    uint32 -= (uint32 >>> 1) & 0x55555555
    uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
    return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
  }

  private readonly _data: Uint32Array

  constructor(n: number) {
    this._data = new Uint32Array(1 + ((n / 31) | 0))
  }

  set(i: number): void {
    this._data[(i / 31) | 0] |= 1 << i % 31
  }

  count(): number {
    let res = 0
    for (let i = 0; i < this._data.length; ++i) {
      res += _BS._onesCount32(this._data[i])
    }
    return res
  }

  run(o: _BS): void {
    let c = 1
    for (let i = 0; i < this._data.length; ++i) {
      const v = this._data[i]
      let x = v
      const y = v | o._data[i]
      x += x + c + (~y & 0x7fffffff) // ((1 << 31) >>> 0) - 1
      this._data[i] = x & y
      c = x >>> 31
    }
  }
}

export { LCSBit }

if (require.main === module) {
  // https://leetcode.cn/problems/longest-common-subsequence/
  function longestCommonSubsequence(text1: string, text2: string): number {
    return LCSBit(text1, text2)
  }

  const n1 = 1e5
  const n2 = 1e5
  const arr1 = Array.from({ length: n1 }, () => ~~(Math.random() * n2))
  const arr2 = Array.from({ length: n1 }, () => ~~(Math.random() * n2))
  console.time('LCSBit')
  console.log(LCSBit(arr1, arr2))
  console.timeEnd('LCSBit')
}
