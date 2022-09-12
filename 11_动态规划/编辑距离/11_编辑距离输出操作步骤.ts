/* eslint-disable no-param-reassign */
/* eslint-disable prefer-destructuring */
/* eslint-disable implicit-arrow-linebreak */

const INF = 2e15
type Option = 'INSERT' | 'REMOVE' | 'REPLACE' | 'SILENT'

/**
 * @param {string} word1
 * @param {string} word2
 * @return {number}
 * @description 给 dp 数组增加额外的信息，然后dp复原
 * 在做最优选择时，顺便把操作记录下来(pre数组)，然后就从结果反推具体操作。
 */
function minDistance(word1: string, word2: string): readonly [number, Option[]] {
  const n1 = word1.length
  const n2 = word2.length

  // dp[i][j]表示w1的前i个字母要转换成w2的前j个字母所需的最少操作数
  const dp = Array.from<unknown, number[]>({ length: n1 + 1 }, () => Array(n2 + 1).fill(INF))
  dp[0][0] = 0

  // pre[i][j]表示dp[i][j]是由哪种操作转移过来的
  const pre = Array.from<unknown, Option[]>({ length: n1 + 1 }, () => Array(n2 + 1).fill('SILENT'))

  for (let j = 1; j <= n2; j++) {
    dp[0][j] = j
    pre[0][j] = 'INSERT'
  }

  for (let i = 1; i <= n1; i++) {
    dp[i][0] = i
    pre[i][0] = 'REMOVE'
  }

  for (let i = 1; i <= n1; i++) {
    for (let j = 1; j <= n2; j++) {
      if (word1[i - 1] === word2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
        pre[i][j] = 'SILENT'
      } else {
        const insertNode = [dp[i][j - 1], 'INSERT'] as const
        const removeNode = [dp[i - 1][j], 'REMOVE'] as const
        const replaceNode = [dp[i - 1][j - 1], 'REPLACE'] as const
        let minNode: readonly [count: number, option: Option] = insertNode
        // 有一个严格小于 则取
        if (removeNode[0] < minNode[0]) {
          minNode = removeNode
        }
        if (replaceNode[0] < minNode[0]) {
          minNode = replaceNode
        }

        dp[i][j] = minNode[0] + 1
        pre[i][j] = minNode[1]
      }
    }
  }

  return [dp[n1][n2], getPath(n1, n2)] as const

  // !dp复原 从(i,j)开始倒推一直到(0,0) 读取pre数组里的父结点信息/操作信息
  function getPath(i: number, j: number): Option[] {
    const res: Option[] = []

    while (i > 0 || j > 0) {
      const type = pre[i][j]
      res.push(type)
      switch (type) {
        case 'REMOVE':
          i--
          break
        case 'INSERT':
          j--
          break
        case 'REPLACE':
          i--
          j--
          break
        case 'SILENT':
          i--
          j--
          break
        default:
          throw new Error('Invalid Option')
      }
    }

    return res.reverse()
  }
}

// [ 3, [ 'REPLACE', 'SILENT', 'REMOVE', 'SILENT', 'REMOVE' ] ]
console.log(minDistance('horse', 'ros'))

console.log(minDistance('easdfgh', 'ros'))

export {}
