// 回溯超时
// function minimumDifference(nums: number[]): number {
//   const sum = nums.reduce((pre, cur) => pre + cur, 0)
//   let res = Infinity

//   nums.sort((a, b) => a - b)

//   const bt = (index: number, pathSum: number, visited: number, len: number): number => {
//     if (len === nums.length / 2) {
//       return (res = Math.min(res, Math.abs(sum - 2 * pathSum)))
//     }

//     for (let i = index; i < nums.length; i++) {
//       if (i !== index && nums[i] === nums[i - 1]) continue
//       if (visited & (1 << i)) continue
//       visited |= 1 << i
//       bt(i + 1, pathSum + nums[i], visited, len + 1)
//     }

//     return Infinity
//   }
//   bt(0, 0, 0, 0)

//   return res
// }

/**
 *
 * @param nums
 * @returns
 * 1 <= n <= 15
   给你一个长度为 2 * n 的整数数组。你需要将 nums 分成 两个 长度为 n 的数组，分别求出两个数组的和，
   并 最小化 两个数组和之 差的绝对值 。nums 中每个元素都需要放入两个数组之一。
   https://leetcode-cn.com/problems/partition-array-into-two-arrays-to-minimize-sum-difference/solution/zui-jie-jin-mu-biao-zhi-de-zi-xu-lie-he-m0sq3/
   @description
   1.我们使用了一个小trick，也即将原数组中所有数变为两倍。这样可以保证我们的目标值sum/2是一个整数。
   2.枚举出前一半数和后一半数的全部选举情况后再拼接在一起，问题变成了从16组两个元素个数不超过C(7,15)的列表中找出和最接近原来总和一半的方案
 */
function minimumDifference(nums: number[]): number {
  nums = nums.map(num => num * 2)
  const midIndex = nums.length / 2
  const leftSubArraySum = getSubArraySumFrom(nums.slice(0, midIndex))
  const rightSubArraySum = getSubArraySumFrom(nums.slice(midIndex))
  const target = nums.reduce((pre, cur) => pre + cur, 0) / 2

  let res = Infinity
  for (let leftCount = 0; leftCount <= midIndex; leftCount++) {
    const left = leftSubArraySum[leftCount].sort((a, b) => a - b)
    const right = rightSubArraySum[midIndex - leftCount].sort((a, b) => a - b)
    res = Math.min(res, twoSum(left, right, target))
  }
  return res

  /**
   * @description 计算nums的子序列和 下标表示由多少个数组成
   * @summary 时间复杂度O(2^n)
   */
  function getSubArraySumFrom(nums: number[]): number[][] {
    const n = nums.length
    const res = Array.from<number, number[]>({ length: nums.length + 1 }, () => [])
    for (let i = 0; i < 1 << n; i++) {
      const index = count(i)
      let sum = 0
      for (let j = 0; j < n; j++) {
        if (i & (1 << j)) sum += nums[j]
      }
      res[index].push(sum)
    }

    return res

    /**
     * @description
     * 二进制位1的个数
     */
    function count(num: number) {
      let res = 0
      while (num) {
        num &= num - 1
        res++
      }
      return res
    }
  }

  /**
   * @description
   * 单调不减的数组nums1和nums2分别找到两个数，其和与target的差最小 返回这个最小差值
   */
  function twoSum(nums1: number[], nums2: number[], target: number): number {
    console.log(nums1, nums2, target)
    let l = 0
    let r = nums2.length - 1
    let res = Infinity

    while (l < nums1.length && r > -1) {
      const sum = nums1[l] + nums2[r]
      res = Math.min(res, Math.abs(target - sum))
      if (sum === target) return 0
      else if (sum > target) r--
      else l++
    }

    return res
  }
}
// console.log(minimumDifference([2, -1, 0, 4, -2, -9]))
// console.log(minimumDifference([3, 9, 7, 3]))
// console.log(minimumDifference([-36, 36]))
console.log(minimumDifference([2, -1, 0, 4, -2, -9]))
