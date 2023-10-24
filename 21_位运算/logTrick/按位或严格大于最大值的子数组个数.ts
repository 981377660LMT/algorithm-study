import { SparseTable } from '../../22_专题/RMQ问题/SparseTable'
import { logTrick } from './logTrick'

/**
 * 求有多少个子数组满足：子数组的按位或严格大于子数组的最大值.
 * 1<=nums.length<=2e5, 1<=nums[i]<=1e9.
 */
function countSubarrayWithSumBitwiseOrGreaterThanMax(arr: ArrayLike<number>): number {
  let res = 0
  const st = new SparseTable(arr, () => 0, Math.max)

  // 右端点为right，左端点left的范围满足 pos1<=left<=pos2.
  // 求有多少个左端点left满足：子数组nums[left,right]的按位或严格于k.
  const query = (pos1: number, pos2: number, right: number, k: number): number => {
    if (st.query(pos2, right + 1) >= k) return 0
    let lo = right - pos2
    let hi = right - pos1
    while (lo <= hi) {
      const mid = (lo + hi) >>> 1
      if (st.query(right - mid, right + 1) < k) {
        lo = mid + 1
      } else {
        hi = mid - 1
      }
    }
    const pos = right - lo
    return pos2 - pos
  }

  logTrick(
    arr,
    (a, b) => a | b,
    (leftIntervals, right) => {
      for (let i = 0; i < leftIntervals.length; i++) {
        const { leftStart, leftEnd, value } = leftIntervals[i]
        res += query(leftStart, leftEnd - 1, right, value)
      }
    }
  )

  return res
}

if (require.main === module) {
  console.log(countSubarrayWithSumBitwiseOrGreaterThanMax([1, 2, 3, 4, 5])) // 7
}
