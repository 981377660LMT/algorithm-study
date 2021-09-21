import { bisectLeft } from '../../../9_排序和搜索/二分api/7_二分搜索寻找最左插入位置'
import { bisectInsort } from '../../../9_排序和搜索/二分api/7_二分搜索插入元素'

/**
 *
 * @param nums 在一维数组中，求解和不超过 K 的最大连续子数组之和
 * @param k
 * @link https://leetcode-cn.com/problems/max-sum-of-rectangle-no-larger-than-k/solution/gong-shui-san-xie-you-hua-mei-ju-de-ji-b-dh8s/
 * // 预处理出「前缀和」，然后枚举子数组的左端点，然后通过「二分」来求解其右端点的位置。
   // java里的有序集合treeset基于红黑树实现 都是O(logn)
   // python里的有序集合SortedList基于红黑树实现 都是O(logn)
   // js里需要手动维护一个有序数组 插入是O(n)  其余是O(logn)
 */
const maxSubArraySum = (nums: number[], k: number) => {
  let res = -Infinity
  const sum = Array(nums.length + 1).fill(0)
  for (let i = 1; i <= nums.length; i++) {
    sum[i] = nums[i - 1] + sum[i - 1]
  }

  // 注意是从1开始
  const treeSet: number[] = [0]
  for (let r = 1; r <= nums.length; r++) {
    // 通过前缀和计算 right
    const pre = sum[r]
    // 通过二分找 left:大于等于 pre - k 的最小数
    // 比如 pre = 10， k = 3，就要找大于等于 7 的最小数，这个数不能大于 10。
    const leftIndex = bisectLeft(treeSet, pre - k)
    // 避免数组索引越界
    if (leftIndex < treeSet.length) {
      res = Math.max(res, pre - treeSet[leftIndex])
    }
    bisectInsort(treeSet, pre)
  }

  return res === -Infinity ? -1 : res
}

console.log(maxSubArraySum([1, 2, 3, 4], 8)) // 7
console.log(maxSubArraySum([1, -2, 3, 4], 5)) // 2

// 对于没有负权值的一维数组，我们可以枚举左端点 i，
// 同时利用前缀和的「单调递增」特性，
// 通过「二分」找到符合 sum[j] \leqslant k + sum[i - 1]sum[j]⩽k+sum[i−1]
// 条件的最大值 sum[j]，从而求解出答案。

// 但是如果是含有负权值的话，前缀和将会丢失「单调递增」的特性，
// 我们也就无法使用枚举 i 并结合「二分」查找 j 的做法。

// 这时候需要将过程反过来处理：我们从左到右枚举 j，
// 并使用「有序集合」结构维护遍历过的位置，
// 找到符合 sum[j] - k \leqslant sum[i]sum[j]−k⩽sum[i] 条件的最小值 sum[i]，
// 从而求解出答案。
