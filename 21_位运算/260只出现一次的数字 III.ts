/**
 * @param {number[]} nums
 * @return {number[]}
 * @description 给定一个整数数组 nums，其中恰好有两个元素只出现一次，其余所有元素均出现两次。 找出只出现一次的那两个元素。
 * 如何分开异或的两个数？
 * 右移MSB
 * 由于异或中每个1都只来源于一个数，那么我们可以根据MSB的1将两个数区分开来
   即向右位移MSB位，与1，一个必为1一个必为0
   将两个数分开后就等效于136只出现一次的数字1.js
 */
const singleNumber = (nums: number[]) => {
  const XOR = (arr: number[]) => arr.reduce((pre: number, cur: number) => pre ^ cur, 0)
  const twoNumXOR = XOR(nums)
  const MSB = twoNumXOR.toString(2).length - 1
  const oneNumXOR = XOR(nums.filter(num => ((num >> MSB) & 1) === 0))
  return [oneNumXOR, oneNumXOR ^ twoNumXOR]
}

// console.log(singleNumber([1, 1, 2, 2, 3, 4]))

// 面试题 17.19. 消失的两个数字
// 给定一个数组，包含从 1 到 N 所有的整数，但其中缺了两个数字。你能在 O(N) 时间内只用 O(1) 的空间找到它们吗？
function missingTwo(nums: number[]): number[] {
  const n = nums.length + 2
  let twoNumXOR = getTwoXOR(nums)

  const msb = twoNumXOR.toString(2).length - 1
  const oneNumXOR = getOneXOR(nums)
  return [oneNumXOR, oneNumXOR ^ twoNumXOR]

  function getTwoXOR(nums: number[]) {
    let res = 0
    for (const num of nums) {
      res ^= num
    }
    for (let i = 1; i <= n; i++) {
      res ^= i
    }
    return res
  }

  function getOneXOR(nums: number[]) {
    let res = 0
    for (const num of nums) {
      if ((num >> msb) & 1) continue
      res ^= num
    }
    for (let i = 1; i <= n; i++) {
      if ((i >> msb) & 1) continue
      res ^= i
    }
    return res
  }
}

console.log(missingTwo([1]))
