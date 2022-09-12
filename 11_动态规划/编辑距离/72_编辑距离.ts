/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * @description 不优化的diff算法O(n^3)的原因’
 * 传统Diff算法需要找到两个树的最小更新方式，所以需要[两两]对比每个叶子节点是否相同，
 * 对比就需要O(n^2)次了，
 * 再加上更新（移动、创建、删除）时需要遍历一次，所以是O(n^3)。
 * React认为：一个ReactElement的type不同，那么内容基本不会复用，所以直接删除节点，
 * 添加新节点，这是一个非常大的优化，大大减少了对比时间复杂度。
 * @description 你可以对一个单词进行如下三种操作:增删改
 * @summary 编辑距离
 */
function minDistance(word1: string, word2: string): number {
  const n1 = word1.length
  const n2 = word2.length
  // dp[i][j]表示w1的前i个字母要转换成w2的前j个字母所需的最少操作数
  const dp = Array.from({ length: n1 + 2 }, () => Array(n2 + 1).fill(0))

  for (let i = 0; i < n1 + 1; i++) {
    dp[i][0] = i
  }

  for (let j = 0; j < n2 + 1; j++) {
    dp[0][j] = j
  }

  for (let i = 1; i < n1 + 1; i++) {
    for (let j = 1; j < n2 + 1; j++) {
      // 注意这个序号
      if (word1[i - 1] === word2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        dp[i][j] = Math.min(dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1]) + 1
      }
    }
  }

  // console.table(dp)
  return dp[n1][n2]
}

console.log(minDistance('horse', 'ros'))
// 3
export {}
