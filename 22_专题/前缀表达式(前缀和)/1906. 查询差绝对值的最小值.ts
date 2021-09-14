/**
 * @param {number[]} nums   1 <= nums[i] <= 100  2 <= nums.length <= 105
 * @param {number[][]} queries
 * @return {number[]}
 * 一个数组 a 的 差绝对值的最小值 定义为：0 <= i < j < a.length 且 a[i] != a[j] 的 |a[i] - a[j]| 的 最小值。
 * 如果 a 中所有元素都 相同 ，那么差绝对值的最小值为 -1 。
 */
const minDifference = function (nums: number[], queries: number[][]): number[] {
  const len = nums.length
  const maxVal = nums.reduce((pre, cur) => Math.max(pre, cur), -1)
  // 记录前 i个位置中数值 j出现的次数。
  const preArray = Array.from({ length: len + 1 }, () => Array(maxVal + 1).fill(0))
  // pre前缀和记录每个数在前index个数内出现的次数
  // ┌─────────┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
  // │ (index) │ 0 │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │
  // ├─────────┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
  // │    0    │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │
  // │    1    │ 0 │ 1 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │
  // │    2    │ 0 │ 1 │ 0 │ 1 │ 0 │ 0 │ 0 │ 0 │ 0 │
  // │    3    │ 0 │ 1 │ 0 │ 1 │ 1 │ 0 │ 0 │ 0 │ 0 │
  // │    4    │ 0 │ 1 │ 0 │ 1 │ 1 │ 0 │ 0 │ 0 │ 1 │
  // └─────────┴───┴───┴───┴───┴───┴───┴───┴───┴───┘

  for (let i = 1; i <= len; i++) {
    for (let j = 1; j <= maxVal; j++) {
      preArray[i][j] = preArray[i - 1][j]
    }
    preArray[i][nums[i - 1]]++
  }

  const res = Array<number>(queries.length).fill(-1)
  for (let i = 0; i < queries.length; i++) {
    const [l, r] = queries[i]
    let min = Infinity
    let pre = -Infinity

    // 判断区间内是否有不同数字出现
    // 这样做就相当于我们对query区间中的元素无重复地从小到大进行了一次遍历。

    for (let j = 1; j <= maxVal; j++) {
      // 谁出现在query区间内，preArray的两行相减即可，并且是无重复的遍历
      if (preArray[r + 1][j] - preArray[l][j]) {
        min = Math.min(min, j - pre)
        pre = j
      }
    }

    if (min !== Infinity) res[i] = min
  }
  console.table(preArray)

  return res
}

console.log(
  minDifference(
    [1, 3, 4, 8],
    [
      [0, 1],
      [1, 2],
      [2, 3],
      [0, 3],
    ]
  )
)
// 输出：[2,1,4,1]
// 解释：查询结果如下：
// - queries[0] = [0,1]：子数组是 [1,3] ，差绝对值的最小值为 |1-3| = 2 。
// - queries[1] = [1,2]：子数组是 [3,4] ，差绝对值的最小值为 |3-4| = 1 。
// - queries[2] = [2,3]：子数组是 [4,8] ，差绝对值的最小值为 |4-8| = 4 。
// - queries[3] = [0,3]：子数组是 [1,3,4,8] ，差的绝对值的最小值为 |3-4| = 1 。

// 如果题目值域很大，比如 $10^7$ ，那就不适合使用此算法了。
