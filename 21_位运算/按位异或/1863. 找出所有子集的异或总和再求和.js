/**
 * @param {number[]} nums  1 <= nums.length <= 12
 * @return {number}
 * 时间复杂度：O(n*2**n)
 */
var subsetXORSum1 = function (nums) {
  let res = 0
  const n = nums.length

  // 遍历所有子集
  for (let i = 0; i < 1 << n; i++) {
    let xor = 0
    // 遍历每个元素
    for (let j = 0; j < n; j++) {
      if (i & (1 << j)) xor ^= nums[j]
    }
    res += xor
  }

  return res
}

// 按位考虑 子集中第i位对结果贡献位count*(1<<i) 其中count为子集中第i位元素为1的元素个数为奇数的子集数 有2**(n-1)个
// 时间复杂度：O(n)
// 或:
const subsetXORSum2 = function (nums) {
  let res = 0

  for (const num of nums) {
    res |= num
  }

  return res << (nums.length - 1)
}

console.log(subsetXORSum2([5, 1, 6]))
// 输出：28
// 解释：[5,1,6] 共有 8 个子集：
// - 空子集的异或总和是 0 。
// - [5] 的异或总和为 5 。
// - [1] 的异或总和为 1 。
// - [6] 的异或总和为 6 。
// - [5,1] 的异或总和为 5 XOR 1 = 4 。
// - [5,6] 的异或总和为 5 XOR 6 = 3 。
// - [1,6] 的异或总和为 1 XOR 6 = 7 。
// - [5,1,6] 的异或总和为 5 XOR 1 XOR 6 = 2 。
// 0 + 5 + 1 + 6 + 4 + 3 + 7 + 2 = 28

// 按位考虑取值情况
// 一个子集的异或总和中某位为 0 当且仅当子集内该位为 1 的元素数量为偶数
// 只有2**0可能是奇数
