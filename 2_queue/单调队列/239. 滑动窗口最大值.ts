import { ArrayDeque } from '../Deque'

/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * `239. 滑动窗口最大值`
 * 单调双端队列结构解决滑动窗口问题：Java里的双端队列用LinkedList实现，我们这里用循环双端数组实现
 */
const maxSlidingWindow1 = function (nums: number[], k: number): number[] {
  const monoQueue = new ArrayDeque<number>(10 ** 5)
  const res = []

  for (let i = 0; i < nums.length; i++) {
    while (monoQueue.length && nums[monoQueue.rear()!] < nums[i]) {
      monoQueue.pop()
    }

    monoQueue.push(i)

    // remove first element if it's outside the window
    if (i - k === monoQueue.front()) monoQueue.shift() // O(1)

    // 需要添加了
    if (i >= k - 1) {
      res.push(nums[monoQueue.front()!])
    }
  }

  return res
}

export {}
