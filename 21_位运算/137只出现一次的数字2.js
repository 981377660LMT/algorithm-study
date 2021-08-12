/**
 * @param {number[]} nums -2**31 <= nums[i] <= 2**31 - 1
 * @return {number}
 * @description 除了一个数字出现一次，其他都出现了三次，让我们找到出现一次的数
 * @summary 和的方法
 */
const singleNumber = function (nums) {
  const sum = arr => arr.reduce((pre, cur) => pre + cur, 0)
  return (3 * sum([...new Set(nums)]) - sum(nums)) / 2
}

// 考察每位1出现的次数是否为3的倍数的方法
const singleNumber = function (nums) {
  let res = 0
  for (let i = 0; i < 32; i++) {
    let count = 0
    const bit = 1 << i
    for (const num of nums) (num & bit) !== 0 && count++
    count % 3 !== 0 && (res |= bit)
  }
  return res
}
// https://leetcode-cn.com/problems/single-number-ii/solution/luo-ji-dian-lu-jiao-du-xiang-xi-fen-xi-gai-ti-si-l/
//第一时间应该想到的是找到一种逻辑操作，可以满足 1*1*1=0 且 0*1=1*0=1
// 出现0次为0，出现1次为1，出现2次的值无所谓，出现3次就又回到0
// var singleNumber = function (nums) {
//   const sum = arr => arr.reduce((pre, cur) => pre + cur, 0)
//   return (3 * sum([...new Set(nums)]) - sum(nums)) / 2
// }
console.log(singleNumber([1, 2, 3, 2, 1, 2, 1]))
