import { LISDp, LISMaxSum } from '../../11_动态规划/lis最长上升子序列问题/LIS模板'
import { validMountainArray } from './941. 有效的山脉数组'

/**
 * 最大山脉子序列和.
 */
function maxSumOfMountainSubsequence(
  nums: number[],
  options?: {
    /** 左侧是否严格递增.默认为true. */
    leftStrict?: boolean
    /** 右侧是否严格递减.默认为true. */
    rightStrict?: boolean
    /** 是否允许某一侧为空.默认为false. */
    allowEmptySide?: boolean
  }
): number {
  const { leftStrict = true, rightStrict = true, allowEmptySide = false } = options || {}
  const n = nums.length
  const revNums = nums.slice().reverse()
  const preMax = LISMaxSum(nums, leftStrict)
  preMax.unshift(0)
  const sufMax = LISMaxSum(revNums, rightStrict).reverse()
  sufMax.push(0)

  let res = 0
  if (allowEmptySide) {
    for (let i = 0; i < n; i++) {
      res = Math.max(res, preMax[i + 1] + sufMax[i] - nums[i])
    }
  } else {
    const leftLen = LISDp(nums, leftStrict)
    const rightLen = LISDp(revNums, rightStrict).reverse()
    for (let i = 0; i < n; i++) {
      if (leftLen[i] < 2 || rightLen[i] < 2) continue
      res = Math.max(res, preMax[i + 1] + sufMax[i] - nums[i])
    }
  }

  return res
}

export { maxSumOfMountainSubsequence }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function check(
    nums: number[],
    options?: {
      leftStrict?: boolean
      rightStrict?: boolean
      allowEmptySide?: boolean
    }
  ): number {
    let res = 0
    for (let state = 0; state < 1 << nums.length; state++) {
      const arr = []
      for (let i = 0; i < nums.length; i++) {
        if ((state >> i) & 1) arr.push(nums[i])
      }
      if (validMountainArray(arr, options)) {
        res = Math.max(
          res,
          arr.reduce((a, b) => a + b, 0)
        )
      }
    }
    return res
  }

  for (let _ = 0; _ < 100; _++) {
    const nums = Array.from({ length: 18 }, () => ~~(Math.random() * 100))
    const leftStrict = Math.random() < 0.5
    const rightStrict = Math.random() < 0.5
    const allowEmptySide = Math.random() < 0.5
    const res1 = maxSumOfMountainSubsequence(nums, {
      leftStrict,
      rightStrict,
      allowEmptySide
    })
    const res2 = check(nums, {
      leftStrict,
      rightStrict,
      allowEmptySide
    })
    if (res1 !== res2) {
      console.log(nums)
      console.log(res1, res2)
      break
    }
  }

  console.log('ok!')
}
