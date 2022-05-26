/**
 * @param {number[]} nums
 * @return {number}  请你返回将 arr 变成 nums 的最少函数调用次数。
 * 本题给定了两种操作：
   让序列中某个数加 1；
   让序列中所有数全体乘以 2。
   询问你需要操作多少次，才能得到目标数组。
   @summary
   我们可以采用逆向思维，从目标数组转化为初始数组，支持两种操作：
   让序列中某个数减 1(单个算)；
   让序列中所有数全体除以 2(一起算取最大)。
 */
function minOperations(nums: number[]): number {
  let add = 0
  let maxMulti = 0
  for (let num of nums) {
    let multi = 0

    while (num) {
      if (num & 1) {
        add++
        num--
      }

      if (num >= 2) {
        multi++
        num /= 2
      }
    }

    maxMulti = Math.max(maxMulti, multi)
  }

  return add + maxMulti
}

console.log(minOperations([1, 5]))

export {}
