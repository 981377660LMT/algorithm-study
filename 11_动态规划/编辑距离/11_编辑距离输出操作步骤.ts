type Value = number
type NodeType = 'INSERT' | 'REMOVE' | 'REPLACE' | 'SILENT'
type Node = [Value, NodeType]

/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * @description 给 dp 数组增加额外的信息即可
 * 在做最优选择时，顺便把操作记录下来，然后就从结果反推具体操作。
 */
const minDistance = function (word1: string, word2: string) {
  const m = word1.length
  const n = word2.length

  // dp[i][j]表示w1的前i个字母要转换成w2的前j个字母所需的最少操作数
  const dp = Array.from<unknown, Node[]>({ length: m + 1 }, () =>
    Array(n + 1)
      .fill(0)
      .map(_ => [0, 'SILENT'])
  )

  for (let j = 1; j <= n; j++) {
    dp[0][j][0] = j
    dp[0][j][1] = 'INSERT'
  }

  for (let i = 1; i <= m; i++) {
    dp[i][0][0] = i
    dp[i][0][1] = 'REMOVE'
  }

  // console.table(dp)
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      // 注意这个序号
      if (word1[i - 1] === word2[j - 1]) {
        dp[i][j][0] = dp[i - 1][j - 1][0]
        dp[i][j][1] = 'SILENT'
      } else {
        const insertNode = [dp[i][j - 1][0], 'INSERT'] as Node
        const removeNode = [dp[i - 1][j][0], 'REMOVE'] as Node
        const replaceNode = [dp[i - 1][j - 1][0], 'REPLACE'] as Node
        let minNode = insertNode as Node
        // 有一个严格小于 则取
        if (removeNode[0] < minNode[0]) minNode = removeNode
        if (replaceNode[0] < minNode[0]) minNode = replaceNode
        dp[i][j][0] = minNode[0] + 1
        dp[i][j][1] = minNode[1]
      }
    }
  }

  return [dp[m][n], getPath(m, n)]

  function getPath(x: number, y: number): NodeType[] {
    const res: NodeType[] = []

    while (x > 0 || y > 0) {
      const type = dp[x][y][1]
      res.push(type)
      switch (type) {
        case 'REMOVE':
          x--
          break
        case 'INSERT':
          y--
          break
        case 'REPLACE':
        case 'SILENT':
          x--
          y--
          break
        default:
          throw new Error('Invalid Node')
      }
    }

    return res.reverse()
  }
}

console.log(minDistance('horse', 'ros'))
console.log(minDistance('easdfgh', 'ros'))

// 3
export {}
