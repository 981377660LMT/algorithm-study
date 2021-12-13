/**
 * @param {string[]} strs
 * @return {number}
 * 执行删除操作之后，最终得到的数组的行中的每个元素都是按字典序排列的。
 * 返回需要删除的最小个数。
 */
const minDeletionSize = (strs: string[]): number => {
  const n = strs[0].length
  // dp[i]表示以第i个字符结尾所能达到的最长字典序子序列,但是这里的针对对象不是单个字符串,而是多个字符串
  const dp = Array(n).fill(1)

  for (let i = 1; i < dp.length; i++) {
    nextPosition: for (let j = 0; j < i; j++) {
      for (const str of strs) {
        if (str.codePointAt(i)! < str.codePointAt(j)!) continue nextPosition
      }
      dp[i] = Math.max(dp[i], dp[j] + 1)
    }
  }

  // console.log(dp)
  return n - Math.max(...dp)
}

console.log(minDeletionSize(['babca', 'bbazb']))
// 输出：3
// 解释：
// 删除 0、1 和 4 这三列后，最终得到的数组是 A = ["bc", "az"]。
// 这两行是分别按字典序排列的（即，A[0][0] <= A[0][1] 且 A[1][0] <= A[1][1]）。
// 注意，A[0] > A[1] —— 数组 A 不一定是按字典序排列的。
export default 1
