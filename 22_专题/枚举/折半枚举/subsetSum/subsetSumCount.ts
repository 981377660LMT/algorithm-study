/**
 * 2^n个子集中, 有多少个子集的和在[floor, higher)之间.
 * O(2^{N/2}).2^25(3e7) -> 1.7s.
 */
function subsetSumCount(nums: number[], floor: number, higher: number): number {
  const n = nums.length
  const mid = n >>> 1
  const nums1 = nums.slice(0, mid)
  const nums2 = nums.slice(mid)
  const dp1 = subsetSumSorted(nums1)
  const dp2 = subsetSumSorted(nums2)
  return cal(higher) - cal(floor)

  function cal(limit: number): number {
    let res = 0
    let right = dp2.length
    for (let i = 0; i < dp1.length; i++) {
      while (right && dp1[i] + dp2[right - 1] >= limit) right--
      res += right
    }
    return res
  }

  function subsetSumSorted(arr: ArrayLike<number>): number[] {
    let dp = [0]
    for (let i = 0; i < arr.length; i++) {
      const tmp = Array<number>(dp.length)
      for (let j = 0; j < tmp.length; j++) {
        tmp[j] = dp[j] + arr[i]
      }

      const ndp = Array<number>(dp.length + tmp.length)
      let ptr1 = 0
      let ptr2 = 0
      let ptr3 = 0
      while (ptr1 < dp.length && ptr2 < tmp.length) {
        if (dp[ptr1] < tmp[ptr2]) ndp[ptr3++] = dp[ptr1++]
        else ndp[ptr3++] = tmp[ptr2++]
      }
      while (ptr1 < dp.length) ndp[ptr3++] = dp[ptr1++]
      while (ptr2 < tmp.length) ndp[ptr3++] = tmp[ptr2++]

      dp = ndp
    }

    return dp
  }
}

/**
 * 2^n个子集中, 有多少个子集的和在[floor, higher)之间, 求出大小为 0,1,...,n 时的子集的个数.
 * O(2^{N/2}).2^20(1e6) -> 300ms.
 */
function subsetSumCountBySize(nums: number[], floor: number, higher: number): number[] {
  const n = nums.length
  const mid = n >>> 1
  const nums1 = nums.slice(0, mid)
  const nums2 = nums.slice(mid)
  const dp1 = getDp(nums1)
  const dp2 = getDp(nums2)
  const count1 = cal(higher)
  const count2 = cal(floor)
  for (let i = 0; i < count1.length; i++) count1[i] -= count2[i]
  return count1

  function cal(limit: number): number[] {
    const res = Array<number>(n + 1).fill(0)
    for (let s1 = 0; s1 < dp1.length; s1++) {
      const cand1 = dp1[s1]
      for (let s2 = 0; s2 < dp2.length; s2++) {
        const cand2 = dp2[s2]
        let right = cand2.length
        for (let i = 0; i < cand1.length; i++) {
          while (right && cand1[i] + cand2[right - 1] >= limit) right--
          res[s1 + s2] += right
        }
      }
    }
    return res
  }

  function getDp(arr: ArrayLike<number>): number[][] {
    let dp = [0, 0] // [sum, size]
    for (let i = 0; i < arr.length; i++) {
      const tmp = Array<number>(dp.length)
      for (let j = 0; j < tmp.length; j += 2) {
        tmp[j] = dp[j] + arr[i]
        tmp[j + 1] = dp[j + 1] + 1
      }

      const ndp = Array<number>(dp.length + tmp.length)
      let ptr1 = 0
      let ptr2 = 0
      let ptr3 = 0
      while (ptr1 < dp.length && ptr2 < tmp.length) {
        if (dp[ptr1] < tmp[ptr2]) {
          ndp[ptr3++] = dp[ptr1++]
          ndp[ptr3++] = dp[ptr1++]
        } else {
          ndp[ptr3++] = tmp[ptr2++]
          ndp[ptr3++] = tmp[ptr2++]
        }
      }
      while (ptr1 < dp.length) {
        ndp[ptr3++] = dp[ptr1++]
        ndp[ptr3++] = dp[ptr1++]
      }
      while (ptr2 < tmp.length) {
        ndp[ptr3++] = tmp[ptr2++]
        ndp[ptr3++] = tmp[ptr2++]
      }

      dp = ndp
    }

    const res = Array<number[]>(arr.length + 1)
    for (let i = 0; i < res.length; i++) res[i] = []
    for (let i = 0; i < dp.length; i += 2) {
      res[dp[i + 1]].push(dp[i])
    }

    return res
  }
}

export { subsetSumCount, subsetSumCountBySize }

if (require.main === module) {
  const n = 46
  const nums = Array.from({ length: n }, () => Math.floor(Math.random() * 100))
  console.time('subsetSumCount')

  const a = subsetSumCount(nums, 0, ~~((n * 100) / 2))
  console.timeEnd('subsetSumCount')

  console.time('subsetSumCountBySize')
  const b = subsetSumCountBySize(nums, 0, ~~((n * 100) / 2)) // 46 -> 2s
  console.timeEnd('subsetSumCountBySize')
}
