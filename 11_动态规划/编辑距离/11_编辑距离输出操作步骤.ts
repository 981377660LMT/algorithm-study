type NodeType = 'INSERT' | 'REMOVE' | 'REPLACE' | 'SILENT'

class Node {
  private value: number
  private type: NodeType
  constructor(value: number, type: NodeType) {
    this.value = value
    this.type = type
  }
}

/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * @description 给 dp 数组增加额外的信息即可
 */
const minDistance = function (word1: string, word2: string) {
  const dpRow = word1.length + 1
  const dpCol = word2.length + 1
  // dp[i][j]表示w1的前i个字母要转换成w2的前j个字母所需的最少操作数
  const dp = Array.from({ length: dpRow }, () => Array(dpCol).fill(0))

  for (let j = 0; j < dpCol; j++) {
    dp[0][j] = j
  }

  for (let i = 0; i < dpRow; i++) {
    dp[i][0] = i
  }

  for (let i = 1; i < dpRow; i++) {
    for (let j = 1; j < dpCol; j++) {
      // 注意这个序号
      if (word1[i - 1] === word2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        dp[i][j] = Math.min(dp[i - 1][j], dp[i][j - 1], dp[i - 1][j - 1]) + 1
      }
    }
  }

  console.table(dp)
  return dp[dpRow - 1][dpCol - 1]
}

console.log(minDistance('horse', 'ros'))
// 3
export {}
