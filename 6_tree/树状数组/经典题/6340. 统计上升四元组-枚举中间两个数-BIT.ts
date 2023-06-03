import { BITArray } from './BIT'

// !i1 i2 i3 i4 满足的数需要为 1 3 2 4 大小关系 (1324模式)
// 求满足条件的四元组的个数 n<=4000
function countQuadruplets(nums: number[]): number {
  const n = nums.length
  const leftSmaller = new BITArray(n + 10)

  let res = 0
  for (let i2 = 0; i2 < n; i2++) {
    const num2 = nums[i2]
    const rightBigger = new BITArray(n + 10)
    for (let i = i2 + 1; i < n; i++) {
      rightBigger.add(nums[i], 1)
    }

    for (let i3 = i2 + 1; i3 < n; i3++) {
      const num3 = nums[i3]
      rightBigger.add(num3, -1)
      if (num2 <= num3) {
        continue
      }

      const count1 = leftSmaller.query(num3 + 1) // 统计i2左侧严格小于num3的数的个数
      const count2 = rightBigger.queryRange(num2 + 1, rightBigger.length) // 统计i3右侧严格大于num2的数的个数
      res += count1 * count2
    }

    leftSmaller.add(num2, 1)
  }

  return res
}

console.log(countQuadruplets([1, 3, 2, 4, 5])) // 2

export {}
