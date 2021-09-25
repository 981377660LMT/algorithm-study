import { Deque } from './Deque'

/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 * `239. 滑动窗口最大值`
 * 单调双端队列结构解决滑动窗口问题：Java里的双端队列用LinkedList实现，我们这里用循环双端数组实现
 */
const maxSlidingWindow1 = function (nums: number[], k: number): number[] {
  const monoQueue = new Deque<number>(10 ** 5)
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

// 奇怪 296 ms
function goodMaxSlidingWindow(nums: number[], k: number): number[] {
  const queue = []
  for (let i = 0; i < k; i++) {
    while (queue.length && nums[i] >= nums[queue[queue.length - 1]]) {
      queue.pop()
    }
    queue.push(i)
  }

  const ret = [nums[queue[0]]]

  // 研究一下整理
  for (let i = k; i < nums.length; i++) {
    while (queue.length && nums[i] >= nums[queue[queue.length - 1]]) {
      queue.pop()
    }
    queue.push(i)

    if (queue[0] === i - k) queue.shift()
    ret.push(nums[queue[0]])
  }

  return ret
}

// 4888 ms
const badMaxSlidingWindow = function (nums: number[], k: number): number[] {
  const monoStack = []
  const res = []

  for (let i = 0; i < nums.length; i++) {
    while (monoStack.length && nums[monoStack[monoStack.length - 1]] < nums[i]) {
      monoStack.pop()
    }
    monoStack.push(i)

    // remove first element if it's outside the window
    if (i - k === monoStack[0]) monoStack.shift() // shift了nums.length次 总的复杂度O(n^2)
    // 需要添加了
    if (i >= k - 1) {
      res.push(nums[monoStack[0]])
    }
  }

  return res
}
export {}
