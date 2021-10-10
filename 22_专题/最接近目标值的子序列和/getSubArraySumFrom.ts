/**
 * @description 计算nums全部子序列和
 * @summary 时间复杂度O(2^n) 小于取所有子集的复杂度O(2^n*n)
 */
function getSubArraySumFrom(nums: number[]): number[] {
  const n = nums.length
  const res = Array<number>(1 << n).fill(0)

  // 外层遍历数组每个元素，遍历到时，表示取该元素
  for (let i = 0; i < n; i++) {
    // 内层遍历从0到外层元素之间到每一个元素，表示能取到的元素，由于前面的结果已经计算过，因此可以直接累加
    for (let j = 0; j < 1 << i; j++) {
      res[(1 << i) + j] = res[j] + nums[i]
    }
  }

  return res
}

if (require.main === module) {
  console.log(getSubArraySumFrom([1, 2]))
}

export { getSubArraySumFrom }
