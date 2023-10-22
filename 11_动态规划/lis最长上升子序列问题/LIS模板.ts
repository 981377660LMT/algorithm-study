import { FastSet } from '../../24_高级数据结构/珂朵莉树/FastSet'
import { SegmentTreeDynamic } from '../../6_tree/线段树/template/动态开点/SegmentTreeDynamicSparse'
import { bisectLeft, bisectRight } from '../../9_排序和搜索/二分/bisect'

function LIS(nums: ArrayLike<number>, isStrict = true): number {
  const n = nums.length
  if (n <= 1) return n

  const lis: number[] = []
  const bisect = isStrict ? bisectLeft : bisectRight
  for (let i = 0; i < n; i++) {
    const pos = bisect(lis, nums[i])
    if (pos === lis.length) {
      lis.push(nums[i])
    } else {
      lis[pos] = nums[i]
    }
  }

  return lis.length
}

/**
 * 求以每个元素结尾的LIS长度.
 */
function LISDp(nums: ArrayLike<number>, isStrict = true): Uint32Array {
  const n = nums.length
  if (n <= 1) return new Uint32Array(n).fill(1)

  const lis: number[] = []
  const dp = new Uint32Array(n)
  const bisect = isStrict ? bisectLeft : bisectRight
  for (let i = 0; i < n; i++) {
    const pos = bisect(lis, nums[i])
    if (pos === lis.length) {
      lis.push(nums[i])
      dp[i] = lis.length
    } else {
      lis[pos] = nums[i]
      dp[i] = pos + 1
    }
  }

  return dp
}

/**
 * 求LIS 返回(LIS,LIS的组成下标).
 */
function getLIS(nums: ArrayLike<number>, isStrict = true): [lis: number[], lisIndex: number[]] {
  const n = nums.length

  const lis: number[] = [] // lis[i] 表示长度为 i 的上升子序列的最小末尾值
  const dpIndex = new Uint32Array(n) // 每个元素对应的LIS长度
  const bisect = isStrict ? bisectLeft : bisectRight
  for (let i = 0; i < n; i++) {
    const pos = bisect(lis, nums[i])
    if (pos === lis.length) {
      lis.push(nums[i])
    } else {
      lis[pos] = nums[i]
    }
    dpIndex[i] = pos
  }

  const res: number[] = []
  const resIndex: number[] = []
  let j = lis.length - 1
  for (let i = n - 1; i >= 0; i--) {
    if (dpIndex[i] === j) {
      res.push(nums[i])
      resIndex.push(i)
      j -= 1
    }
  }

  return [res.reverse(), resIndex.reverse()]
}

/**
 * 求和最大的LIS.
 * @param nums nums[i]>=0.
 * @param isStrict 是否严格递增.默认为true.
 * @returns res[i] 表示以nums[i]结尾的LIS的最大和.
 */
function LISMaxSum(nums: ArrayLike<number>, isStrict = true): number[] {
  const n = nums.length
  if (n <= 1) return Array.from(nums)

  let max = 0
  for (let i = 0; i < n; i++) max = Math.max(max, nums[i])
  const dp = new SegmentTreeDynamic(0, max, () => 0, Math.max)
  const res = Array(n)

  for (let i = 0; i < n; i++) {
    const num = nums[i]
    const preMax = isStrict ? dp.query(0, num) : dp.query(0, num + 1)
    res[i] = preMax + num
    dp.update(num, res[i])
  }

  return res
}

/**
 * O(nloglogmax)求`严格递增`的LIS,要求所有元素范围在`[0, max]`内.
 * @param nums 数组.
 * @param max 数组的最大值.不超过1e9.
 */
function LIS2Strict(nums: ArrayLike<number>, max: number): number {
  max += 5
  const alive = new FastSet(max)
  let res = 0

  for (let i = 0; i < nums.length; i++) {
    const cur = nums[i]
    const next = alive.next(cur)
    if (next < max) {
      alive.erase(next)
      alive.insert(cur)
    } else {
      alive.insert(cur)
      res++
    }
  }
  return res
}

export { LIS, LISDp, getLIS, LISMaxSum }

if (require.main === module) {
  // https://leetcode.cn/problems/longest-increasing-subsequence/
  // eslint-disable-next-line no-inner-declarations
  function lengthOfLIS(nums: number[]): number {
    for (let i = 0; i < nums.length; i++) nums[i] += 1e4
    return LIS2Strict(nums, 2e4)
    return LIS(nums)
  }

  const arr = [1, 2, 2, 3, 3, 3, 1, 4, 3]
  console.log(LIS2Strict(arr, 10))
  console.log(LIS2Strict(arr, 10))
  console.log(getLIS(arr))
  console.log(LISMaxSum(arr, false))
  console.log(LISMaxSum([10, 9, 2, 5, 3, 7, 101, 18]))
}
