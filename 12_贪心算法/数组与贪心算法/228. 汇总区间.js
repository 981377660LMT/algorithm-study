/**
 * @param {number[]} nums
 * @return {string[]}
 */
var summaryRanges = function (nums) {
  const res = []
  nums.push(Infinity) // 哨兵元素
  let pre = nums[0]
  for (let i = 1; i < nums.length; i++) {
    if (nums[i] - nums[i - 1] === 1) continue
    if (nums[i - 1] === pre) {
      // 只有一个数
      res.push(pre.toString())
    } else {
      // 一段数
      res.push(`${pre}->${nums[i - 1]}`)
    }
    pre = nums[i]
  }

  return res
}

console.log(summaryRanges([0, 1, 2, 4, 5, 7]))
// 输出：["0->2","4->5","7"]
// 解释：区间范围是：
// [0,2] --> "0->2"
// [4,5] --> "4->5"
// [7,7] --> "7"
