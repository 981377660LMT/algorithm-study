// 给定一个长度为 n 的 非空 整数数组，每次操作将会使 n - 1 个元素增加 1。找出让数组所有元素相等的最小操作次数。
// 等价于每次将一个数减1
/**
 * @param {number[]} nums
 * @return {number}
 */
var minMoves = function (nums) {
  const min = Math.min.apply(null, nums)
  return nums.reduce((pre, cur) => pre + cur - min, 0)
}
