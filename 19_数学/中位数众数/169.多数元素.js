/**
 * @param {number[]} nums
 * @return {number}
 * 给定一个大小为 n 的数组，找到其中的多数元素。多数元素是指在数组中出现次数 大于 ⌊ n/2 ⌋ 的元素。
 */
// var majorityElement = function (nums) {
//   return nums.sort()[Math.floor(nums.length / 2)]
// }

// 摩尔投票法：
// 核心就是对拼消耗。
// 从第一个数开始count=1，遇到相同的就加1，遇到不同的就减1，
// 减到0就重新换个数开始计数，总能找到最多的那个
// 假设你方人口超过总人口一半以上，并且能保证每个人口出去干仗都能一对一同归于尽。最后还有人活下来的国家就是胜利。
var majorityElement = function (nums) {
  const count = (nums, target) => {
    let res = 0
    const n = nums.length
    for (let i = 0; i < n; i++) {
      nums[i] === target && res++
    }
    return res
  }
  // 虚拟元素
  let candidate = 0
  let count = 0
  for (const num of nums) {
    if (candidate === num) count++
    else if (count === 0) {
      candidate = num
      count = 1
    } else {
      count--
    }
  }
  return count(nums, candidate) > nums.length >> 1 ? candidate : -1
}

console.log(majorityElement([2, 2, 1, 1, 1, 2, 2]))
