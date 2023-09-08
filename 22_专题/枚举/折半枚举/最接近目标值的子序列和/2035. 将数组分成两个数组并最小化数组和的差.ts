import { groupSubsetSumBySize } from '../subsetSum/subsetSum'
import { twoSum } from '../twoSum'

/**
 * 2035. 将数组分成两个数组并最小化数组和的差
 * https://leetcode.cn/problems/partition-array-into-two-arrays-to-minimize-sum-difference/description/
 * 1 <= n <= 15.
   给你一个长度为 2 * n 的整数数组。你需要将 nums 分成 两个 长度为 n 的数组，分别求出两个数组的和，
   并 最小化 两个数组和之 差的绝对值 。nums 中每个元素都需要放入两个数组之一。
   https://leetcode-cn.com/problems/partition-array-into-two-arrays-to-minimize-sum-difference/solution/zui-jie-jin-mu-biao-zhi-de-zi-xu-lie-he-m0sq3/
   1.我们使用了一个小trick，也即将原数组中所有数变为两倍。这样可以保证我们的目标值sum/2是一个整数。
   2.枚举出前一半数和后一半数的全部选举情况后再拼接在一起，问题变成了从16组两个元素个数不超过C(7,15)的列表中找出和最接近原来总和一半的方案
 */
function minimumDifference(nums: number[]): number {
  nums = nums.map(v => v * 2)
  const target = Math.floor(nums.reduce((pre, cur) => pre + cur, 0) / 2)
  const mid = nums.length >>> 1
  const leftSum = groupSubsetSumBySize(nums.slice(0, mid))
  const rightSum = groupSubsetSumBySize(nums.slice(mid))

  let res = Infinity
  for (let leftCount = 0; leftCount <= mid; leftCount++) {
    const left = leftSum[leftCount].sort((a, b) => a - b)
    const right = rightSum[mid - leftCount].sort((a, b) => a - b)
    res = Math.min(res, twoSum(left, right, target))
  }
  return res
}

export {}

if (require.main === module) {
  // [2,-1,0,4,-2,-9]

  console.log(minimumDifference([2, -1, 0, 4, -2, -9]))
}
