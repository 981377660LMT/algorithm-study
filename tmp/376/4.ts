// 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。

import { SortedListFastWithSum } from '../../22_专题/离线查询/根号分治/SortedList/SortedListWithSum'

// 你可以对数组执行 至多 k 次操作：

// 从数组中选择一个下标 i ，将 nums[i] 增加 或者 减少 1 。
// 最终数组的频率分数定义为数组中众数的 频率 。

// 请你返回你可以得到的 最大 频率分数。

// 众数指的是数组中出现次数最多的数。一个元素的频率指的是数组中这个元素的出现次数。

function distSum(sl: SortedListFastWithSum): (k: number) => number {
  return (k: number): number => {
    const pos = sl.bisectRight(k)
    const leftSum = k * pos - sl.sumSlice(0, pos)
    const rightSum = sl.sumSlice(pos, sl.length) - k * (sl.length - pos)
    return leftSum + rightSum
  }
}

function maxFrequencyScore(nums: number[], k: number): number {
  nums.sort((a, b) => a - b)
  const targets = Array(nums.length - 1).fill(0)
  for (let i = 0; i < nums.length - 1; ++i) {
    targets[i] = Math.floor((nums[i + 1] + nums[i]) / 2)
  }

  const sl = new SortedListFastWithSum()
  const D = distSum(sl)
  let left = 0
  let right = 0
  let res = 1

  for (let i = 0; i < targets.length; ++i) {
    const target = targets[i]
    let curSum = D(target)

    // 左端缩小窗口
    while (left <= right && curSum > k) {
      const leftMost = nums[left]
      curSum -= Math.abs(leftMost - target)
      sl.discard(leftMost)
    }

    // 向右扩大窗口
    while (right < nums.length) {
      const nextSum = curSum + (nums[right] - target)
      if (nextSum > k) {
        break
      }
      curSum = nextSum
      sl.add(nums[right])
      console.log(1, sl.length)
    }

    res = Math.max(res, sl.length)
  }

  return res
}
// [1,2,6,4]
// 3
// [1,4,4,2,4]
// 0
console.log(maxFrequencyScore([1, 2, 6, 4], 3))
