// 选出一个包含字母 E 数量与字母 F 数量之差(E-F)最大的子串(子数组)
// 等价于求最大子数组和
const { readFileSync } = require('fs')
const iter = readlines()
const input = () => iter.next().value
function* readlines(path = 0) {
  const lines = readFileSync(path)
    .toString()
    .trim()
    .split(/\r\n|\r|\n/)

  yield* lines
}

const n = Number(input())
const nums = String(input())
  .split('')
  .map(char => (char === 'E' ? 1 : -1))

console.log(kanade(nums, true))

function kanade(nums, getMax = true) {
  if (nums.length === 0) return 0
  if (nums.length === 1) return nums[0]

  let res = getMax ? -Infinity : Infinity
  let sum = 0
  for (const num of nums) {
    if (getMax) {
      sum = Math.max(sum + num, num)
      res = Math.max(res, sum)
    } else {
      sum = Math.min(sum + num, num)
      res = Math.min(res, sum)
    }
  }

  return res
}
