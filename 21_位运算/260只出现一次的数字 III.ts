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
  const XOR = (arr: number[]) => arr.reduce((pre: number, cur: number) => pre ^ cur)
  const twoNumXOR = XOR(nums)
  const MSB = twoNumXOR.toString(2).length - 1
  const oneNumXOR = XOR(nums.filter(num => ((num >> MSB) & 1) === 0))
  return [oneNumXOR, oneNumXOR ^ twoNumXOR]
}

console.log(singleNumber([1, 1, 2, 2, 3, 4]))
