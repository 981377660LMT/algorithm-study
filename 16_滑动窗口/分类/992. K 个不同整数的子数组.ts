/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number}
 * @summary 题目转换成求：不超过k种元素的子数组个数 - 不超过k-1种元素的子数组个数(转化为水果成蓝问题)
 */
const subarraysWithKDistinct = function (nums: number[], k: number): number {
  const notMoreThan = (threshold: number): number => {
    let res = 0
    let left = 0
    const counter = new Map<number, number>()

    for (let right = 0; right < nums.length; right++) {
      const curNum = nums[right]
      counter.set(curNum, (counter.get(curNum) || 0) + 1)

      while (counter.size > threshold) {
        const preNum = nums[left]
        const count = counter.get(preNum)!
        if (count === 1) counter.delete(preNum)
        else counter.set(preNum, count - 1)
        left++
      }

      res += right - left + 1 // 以r结尾的子数组，多少个合理
    }

    return res
  }

  return notMoreThan(k) - notMoreThan(k - 1)
}

console.log(subarraysWithKDistinct([1, 2, 1, 2, 3], 2))
