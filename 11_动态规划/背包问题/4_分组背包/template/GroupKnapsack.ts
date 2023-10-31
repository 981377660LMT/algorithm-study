/* eslint-disable no-inner-declarations */

type Item = { value: number; weight: number }

/**
 * 分组背包，每组最多选一个.最大化价值.
 * @complexity O(n * m * maxCapacity)
 */
function groupKnapsackAtMostOne(groups: ArrayLike<ArrayLike<Item>>, maxCapacity: number): number {
  const dp = Array<number>(maxCapacity + 1).fill(0)
  const n = groups.length
  let maxJ = 0
  for (let i = 0; i < n; i++) {
    const group = groups[i]
    let curMax = 0
    for (let k = 0; k < group.length; k++) curMax = Math.max(curMax, group[k].weight)
    maxJ = Math.min(maxJ + curMax, maxCapacity)
    // 这里 j 的初始值可以优化成前 i 个组的每组最大重量之和（但不能超过 maxW）
    for (let j = maxJ; j >= 0; j--) {
      for (let k = 0; k < group.length; k++) {
        const item = group[k]
        if (item.weight <= j) {
          dp[j] = Math.max(dp[j], dp[j - item.weight] + item.value)
        }
      }
    }
  }

  return dp[maxCapacity]
}

/**
 * 分组背包，每组恰好选一个.最大化价值.
 * @returns dp[j] 表示从每组恰好选一个，能否凑成重量 j.
 * @complexity O(n * m * maxCapacity)
 */
function groupKnapsackExactOne(groups: ArrayLike<ArrayLike<number>>, maxCapacity: number): Uint8Array {
  const dp = new Uint8Array(maxCapacity + 1)
  dp[0] = 1
  const n = groups.length
  let maxJ = 0
  for (let i = 0; i < n; i++) {
    const group = groups[i]
    let curMax = 0
    for (let k = 0; k < group.length; k++) curMax = Math.max(curMax, group[k])
    maxJ = Math.min(maxJ + curMax, maxCapacity)
    // 这里 j 的初始值可以优化成前 i 个组的每组最大重量之和（但不能超过 maxW）
    for (let j = maxJ; j >= 0; j--) {
      let ok = false
      for (let k = 0; k < group.length; k++) {
        const weight = group[k]
        if (weight <= j && dp[j - weight]) {
          dp[j] = 1
          ok = true
          break
        }
      }
      if (!ok) dp[j] = 0
    }
  }
  return dp
}

export { groupKnapsackAtMostOne, groupKnapsackExactOne }

if (require.main === module) {
  // 1981. 最小化目标值与所选元素的差
  // https://leetcode.cn/problems/minimize-the-difference-between-target-and-chosen-elements/
  function minimizeTheDifference(mat: number[][], target: number): number {
    let maxSum = 0
    mat.forEach(row => {
      let max = 0
      for (let i = 0; i < row.length; i++) {
        max = Math.max(max, row[i])
      }
      maxSum += max
    })
    const ok = groupKnapsackExactOne(mat, maxSum)
    let res = 1e9
    for (let cand = 0; cand <= maxSum; cand++) {
      if (ok[cand]) {
        res = Math.min(res, Math.abs(cand - target))
      }
    }
    return res
  }

  // mat = [[1],[2],[3]],
  console.log(minimizeTheDifference([[1], [2], [3]], 100))
}
