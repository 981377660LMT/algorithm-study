// 参考659：分割数组为连续子序列 注意区别
/**
 *
 * @param nums
 * @param k
 * 是否可以把这个数组划分成一些由 k 个连续数字组成的集合(必须等于k)
 * @summary
 * 贪心,先做最小的(可以不用考虑左边), 也可以先做最大的(不用考虑右边),统计数字个数,一旦数字不够划分,就false.
 * 一手顺子顺序出即可
 */
function isPossibleDivide(nums: number[], k: number): boolean {
  if (nums.length % k) return false
  nums.sort((a, b) => a - b)
  const counter = new Map<number, number>()
  nums.forEach(value => counter.set(value, (counter.get(value) || 0) + 1))

  for (const num of nums) {
    if (counter.get(num) === 0) continue // 用完了
    // 耗尽这一段
    for (let i = 0; i < k; i++) {
      const need = num + i
      if (!counter.has(need) || counter.get(need) === 0) return false
      counter.set(need, counter.get(need)! - 1)
    }
  }

  return true
}

if (require.main === module) {
  console.log(isPossibleDivide([1, 2, 3, 3, 4, 4, 5, 6], 4))
}

export { isPossibleDivide }
