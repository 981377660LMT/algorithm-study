/**
 * @param {number[]} nums  正整数
 * @return {string}
 * @summary nums[0] 一定为 分子， nums[1]一定为分母
 * 保持被除数最大 除数最小
 */
var optimalDivision = function (nums: number[]): string {
  if (nums.length === 1) return nums[0].toString()
  if (nums.length === 2) return `${nums[0]}/${nums[1]}`
  const res = nums.join('/')
  const first = res.indexOf('/')
  console.log(first)

  return res.slice(0, first + 1) + '(' + res.slice(first + 1) + ')'
}

console.log(optimalDivision([1000, 100, 10, 2]))
// 输出: "1000/(100/10/2)"
export default 1
