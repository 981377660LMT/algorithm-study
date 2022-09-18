// 6186. 按位或最大的最小子数组长度
// !二分+st表

import { SparseTable } from '../../22_专题/RMQ问题/SparseTable'

function smallestSubarrays(nums: number[]): number[] {
  const n = nums.length
  const st = new SparseTable(nums, (a, b) => a | b)
  const res = []

  for (let i = 0; i < n; i++) {
    const or = st.query(i, n - 1)
    let left = i
    let right = n - 1
    while (left <= right) {
      const mid = Math.floor((left + right) / 2)
      if (st.query(i, mid) === or) right = mid - 1
      else left = mid + 1
    }

    res.push(left - i + 1)
  }

  return res
}

if (require.main === module) {
  console.log(smallestSubarrays([7]))
}

export {}
