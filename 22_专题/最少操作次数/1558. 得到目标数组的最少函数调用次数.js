/**
 * @param {number[]} nums
 * @return {number}
 * @description
 * 给你一个与 nums 大小相同且初始值全为 0 的数组 arr ，
 * 每次操作可以单点+1或者所有数x2
 * 请你调用以上函数得到整数数组 nums 。
 */
function minOperations(nums) {
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
