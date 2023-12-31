/* eslint-disable consistent-return */
/* eslint-disable no-else-return */
/* eslint-disable no-inner-declarations */
/* eslint-disable max-len */

import { enumerateInterval } from '../../../区间问题/区间分解/enumerateInterval'

/**
 * 给定字符集信息和字符s,返回一个查询函数.用于查询s[start:end]间char的个数.
 * 当字符种类很少时,可以用一个counter数组实现区间哈希值的快速计算.
 */
function alphaPresum(stringOrOrds: ArrayLike<number> | string, sigma = 26, offset = 97): (start: number, end: number, ord: number) => number {
  const n = stringOrOrds.length
  if (typeof stringOrOrds === 'string') {
    const ords = Array(n)
    for (let i = 0; i < n; ++i) ords[i] = stringOrOrds.charCodeAt(i)
    stringOrOrds = ords
  }

  const preSum = new Uint32Array((n + 1) * sigma)
  for (let i = 1; i <= n; ++i) {
    const pos = i * sigma
    preSum.set(preSum.subarray(pos - sigma, pos), pos)
    preSum[pos + stringOrOrds[i - 1] - offset]++
  }

  const getCountOfSlice = (start: number, end: number, ord: number): number => {
    if (start < 0) start = 0
    if (end > n) end = n
    if (start >= end) return 0
    return preSum[end * sigma + ord - offset] - preSum[start * sigma + ord - offset]
  }
  return getCountOfSlice
}

export { alphaPresum }

if (require.main === module) {
  // https://leetcode.cn/problems/palindrome-rearrangement-queries/
  // 100129. 回文串重新排列查询
  // 给你一个长度为 偶数 n ，下标从 0 开始的字符串 s 。
  // 同时给你一个下标从 0 开始的二维整数数组 queries ，其中 queries[i] = [ai, bi, ci, di] 。
  // 对于每个查询 i ，你需要执行以下操作：
  // 将下标在范围 0 <= ai <= bi < n / 2 内的 子字符串 s[ai:bi] 中的字符重新排列。
  // 将下标在范围 n / 2 <= ci <= di < n 内的 子字符串 s[ci:di] 中的字符重新排列。
  // 对于每个查询，你的任务是判断执行操作后能否让 s 变成一个 回文串 。
  // 每个查询与其他查询都是 独立的 。
  // !请你返回一个下标从 0 开始的数组 answer ，如果第 i 个查询执行操作后，可以将 s 变为一个回文串，那么 answer[i] = true，否则为 false 。
  // 子字符串 指的是一个字符串中一段连续的字符序列。
  // !s[x:y] 表示 s 中从下标 x 到 y 且两个端点 都包含 的子字符串。
  //
  // !将后半部分区间反转，变为两个区间问题.

  function canMakePalindromeQueries(s: string, queries: number[][]): boolean[] {
    const n = s.length / 2
    const arr1 = Array<number>(n)
    for (let i = 0; i < n; i++) arr1[i] = s.charCodeAt(i) - 97
    const arr2 = Array<number>(n)
    for (let i = 0; i < n; i++) arr2[i] = s.charCodeAt(2 * n - i - 1) - 97
    const diffPreSum = new Uint32Array(n + 1)
    for (let i = 1; i <= n; i++) {
      diffPreSum[i] = diffPreSum[i - 1] + +(arr1[i - 1] !== arr2[i - 1])
    }

    const C1 = alphaPresum(arr1, 26, 0)
    const C2 = alphaPresum(arr2, 26, 0)

    const res = Array<boolean>(queries.length).fill(false)
    for (let qi = 0; qi < queries.length; qi++) {
      let { 0: l1, 1: r1, 2: l2, 3: r2 } = queries[qi]
      const start1 = l1
      const end1 = r1 + 1
      const start2 = 2 * n - r2 - 1
      const end2 = 2 * n - l2

      const hash1 = Array<number>(26)
      for (let i = 0; i < 26; i++) hash1[i] = C1(start1, end1, i)
      const hash2 = Array<number>(26)
      for (let i = 0; i < 26; i++) hash2[i] = C2(start2, end2, i)

      let ok = true
      const enumerate = enumerateInterval([{ start: start1, end: end1, value: 0 }], [{ start: start2, end: end2, value: 0 }])
      enumerate(0, n, (type, start, end, value1, value2) => {
        if (type === '00') {
          const diff = diffPreSum[end] - diffPreSum[start]
          if (diff) {
            ok = false
            return true
          }
        } else if (type === '10') {
          for (let i = 0; i < 26; i++) {
            const count = C2(start, end, i)
            if (count > hash1[i]) {
              ok = false
              return true
            } else {
              hash1[i] -= count
            }
          }
        } else if (type === '01') {
          for (let i = 0; i < 26; i++) {
            const count = C1(start, end, i)
            if (count > hash2[i]) {
              ok = false
              return true
            } else {
              hash2[i] -= count
            }
          }
        }
      })

      if (!ok) {
        res[qi] = false
        continue
      }

      for (let i = 0; i < 26; i++) {
        if (hash1[i] !== hash2[i]) {
          ok = false
          break
        }
      }
      res[qi] = ok
    }

    return res
  }

  // "abbcdecbba" queries = [[0,2,7,9]]
  console.log(canMakePalindromeQueries('abbcdecbba', [[0, 2, 7, 9]]))
}
