/* eslint-disable no-inner-declarations */

import { bisectLeft } from '../../9_排序和搜索/二分/bisect'
import { sortSearch } from '../../9_排序和搜索/二分/sortSearch'

/** LIS求方案数. */
function countLIS(arr: ArrayLike<number>, strict = true): number {
  const lis: number[][] = []
  const countPreSum: number[][] = []
  for (let i = 0; i < arr.length; i++) {
    const num = arr[i]
    const target = strict ? num : num + 1
    const pos = sortSearch(0, lis.length, i => lis[i][lis[i].length - 1] >= target)
    let count = 1
    if (pos > 0) {
      const tmp1 = lis[pos - 1]
      const k = sortSearch(0, tmp1.length, k => tmp1[k] < target)
      const tmp2 = countPreSum[pos - 1]
      count = tmp2[tmp2.length - 1] - tmp2[k]
    }
    if (pos === lis.length) {
      lis.push([num])
      countPreSum.push([0, count])
    } else {
      lis[pos].push(num)
      const tmp = countPreSum[pos]
      tmp.push(tmp[tmp.length - 1] + count)
    }
  }

  const last = countPreSum[countPreSum.length - 1]
  return last[last.length - 1]
}

const MOD = 1e9 + 7

/** 求长为length的LIS个数. */
function countLISWithLength(arr: ArrayLike<number>, length: number, strict = true): number {
  const copy = Array(arr.length)
  const sorted = Array(arr.length)
  for (let i = 0; i < arr.length; i++) {
    copy[i] = arr[i]
    sorted[i] = arr[i]
  }
  sorted.sort((a, b) => a - b)
  for (let i = 0; i < arr.length; i++) {
    copy[i] = bisectLeft(sorted, arr[i]) + 2
  }

  const n = copy.length
  const add = (tree: number[], i: number, val: number): void => {
    for (; i < n + 2; i += i & -i) {
      tree[i] = (tree[i] + val) % MOD
    }
  }
  const sum = (tree: number[], i: number): number => {
    let res = 0
    for (; i > 0; i &= i - 1) res = (res + tree[i]) % MOD
    return res
  }

  const dp = Array(length + 1)
  for (let i = 0; i <= length; i++) {
    dp[i] = Array(n).fill(0)
  }

  for (let i = 1; i <= length; i++) {
    const tree = Array(n + 2).fill(0)
    const tmp1 = dp[i - 1]
    const tmp2 = dp[i]
    if (i === 1) {
      add(tree, 1, 1)
    }
    if (strict) {
      for (let j = 0; j < n; j++) {
        const v = copy[j]
        tmp2[j] = sum(tree, v - 1)
        add(tree, v, tmp1[j])
      }
    } else {
      for (let j = 0; j < n; j++) {
        const v = copy[j]
        tmp2[j] = sum(tree, v)
        add(tree, v, tmp1[j])
      }
    }
  }

  const last = dp[length]
  let res = 0
  for (let i = 0; i < n; i++) {
    res = (res + last[i]) % MOD
  }
  return res
}

export { countLIS, countLISWithLength }

if (require.main === module) {
  // https://leetcode.cn/problems/number-of-longest-increasing-subsequence/description/
  function findNumberOfLIS(nums: number[]): number {
    return countLIS(nums, true)
  }

  console.log(countLISWithLength([1, 1, 1, 1], 2, false))
}
