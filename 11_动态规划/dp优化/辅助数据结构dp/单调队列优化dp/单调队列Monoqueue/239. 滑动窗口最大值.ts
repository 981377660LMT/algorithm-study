import { MaxDeque } from './MaxDequeue'

/**
 * `239. 滑动窗口最大值`
 */
function maxSlidingWindow(nums: number[], k: number): number[] {
  const queue = new MaxDeque<{ value: number }>()
  const res = []

  for (let i = 0; i < nums.length; i++) {
    queue.append({ value: nums[i] })
    if (i >= k) queue.popLeft()
    if (i >= k - 1) res.push(queue.max)
  }

  return res
}

export {}
