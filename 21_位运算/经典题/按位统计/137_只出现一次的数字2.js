/**
 * @param {number[]} nums -2**31 <= nums[i] <= 2**31 - 1
 * @return {number}
 * @description 除了一个数字出现一次，其他都出现了三次，让我们找到出现一次的数
 * @summary 和的方法
 */

// 对每位统计 如果次数不为三的倍数 则 只出现了一次的元素的二进制在这位上是1
function singleNumber(nums) {
  let res = 0
  for (let i = 0; i < 32; i++) {
    let count = 0
    for (const num of nums) if ((num >> i) & 1) count++
    count % 3 !== 0 && (res |= 1 << i)
  }

  return res
}

console.log(singleNumber([1, 2, 3, 2, 1, 2, 1]))
export {}
