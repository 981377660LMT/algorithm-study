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
function singleNumber(nums: number[]) {
  const xor = (arr: number[]) => arr.reduce((pre: number, cur: number) => pre ^ cur, 0)
  const xor2 = xor(nums)
  // 这里可以直接Lowbit
  const MSB = xor2.toString(2).length - 1
  const xor1 = xor(nums.filter(num => (num >> MSB) & 1))
  return [xor1, xor1 ^ xor2]
}
