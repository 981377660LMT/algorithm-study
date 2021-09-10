/**
 * @param {number[]} nums
 * @return {boolean}
 * @description
 * 所有 nums[seq[j]] 应当不是 全正 就是 全负
 * k > 1
 */
var circularArrayLoop = function (nums: number[]): boolean {
  const n = nums.length
  // 注意这里索引模n的处理
  const getNextPos = (index: number) => (((index + nums[index]) % n) + n) % n

  for (let i = 0; i < n; i++) {
    let slow = i
    let fast = getNextPos(slow)
    while (nums[fast] * nums[i] > 0 && nums[getNextPos(fast)] * nums[i] > 0) {
      if (fast === slow) {
        // 只有一个点的情况
        if (slow === getNextPos(slow)) break
        else return true
      }
      slow = getNextPos(slow)
      fast = getNextPos(getNextPos(fast))
    }
  }

  return false
}

console.log(circularArrayLoop([2, -1, 1, 2, 2]))
console.log(circularArrayLoop([-2, 1, -1, -2, -2]))
