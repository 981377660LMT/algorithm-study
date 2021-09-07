// 1 <= nums.length <= 104

// 不妨假设数对前一个数不超过后一个数 类似两数之和的做法
// 统计较小值的不同个数
function findPairs(nums: number[], k: number): number {
  const big = new Set<number>()
  const small = new Set<number>()

  for (const num of nums) {
    if (big.has(num - k)) small.add(num - k)
    if (big.has(num + k)) small.add(num)
    big.add(num)
  }

  return small.size
}

console.log(findPairs([3, 1, 4, 1, 5], 2))
// 输出：2
// 解释：数组中有两个 2-diff 数对, (1, 3) 和 (3, 5)。
// 尽管数组中有两个1，但我们只应返回不同的数对的数量。
export default 2
