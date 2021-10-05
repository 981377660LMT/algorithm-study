import { bisectLeft } from '../../../9_排序和搜索/7_二分搜索寻找最左插入位置'

/**
 * @param {number[]} target 1 <= target.length, arr.length <= 10**5 且 target 数组元素各不相同
 * @param {number[]} arr
 * @return {number}
 * 求最少插入操作次数，使得 target 成为 arr 的一个子序列。
 * 一个数组的 子序列 指的是删除原数组的某些元素（可能一个元素都不删除），同时不改变其余元素的相对顺序得到的数组
 * 由于这道题数据范围是 $10^5$，因此只能使用 $NlogN$ 的贪心才行。
 *
 * @summary
 * target 和 arr 的最长公共子序列长度为max，则最终答案为 n−max。
 * 为何能从 LCS 问题转化为 LIS 问题
 * 当其中一个数组元素各不相同时，最长公共子序列问题（LCS）可以转换为最长上升子序列问题（LIS）进行求解。
 * 根据target中互不相同，我们知道每个数字对应的坐标唯一;
   于是最长公共子序列等价于arr用target的坐标转换后构成最长的上升子序列.
   不管怎么样，公共子序列在target中必然是从左到右的，那么他们的坐标自然是从小到大的
 * 同时最长上升子序列问题（LIS）存在使用「维护单调序列 + 二分」的贪心解法，复杂度为 O(n\log{n})O(nlogn)
 */
const minOperations = (target: number[], arr: number[]): number => {
  const visited = new Set(target)
  // 实际不建议这样写 内存消耗很大
  const valueToIndex = new Map([...target.entries()].map(([i, v]) => [v, i]))
  arr = arr.filter(v => visited.has(v)).map(v => valueToIndex.get(v)!)
  if (arr.length <= 1) return target.length - arr.length

  const LIS: number[] = [arr[0]]
  for (let i = 1; i < arr.length; i++) {
    if (arr[i] > LIS[LIS.length - 1]) {
      LIS.push(arr[i])
    } else {
      LIS[bisectLeft(LIS, arr[i])] = arr[i]
    }
  }

  console.log(valueToIndex, arr)
  return target.length - LIS.length
}

console.log(minOperations([5, 1, 3], [9, 4, 2, 3, 4]))
// 输出：2
// 解释：你可以添加 5 和 1 ，使得 arr 变为 [5,9,4,1,2,3,4] ，target 为 arr 的子序列。

export default 1
