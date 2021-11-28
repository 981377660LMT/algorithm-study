/**
 * @link https://leetcode-cn.com/problems/max-black-square-lcci/solution/20xing-dai-ma-by-suibianfahui/
 * @param {number[][]} matrix matrix.length == matrix[0].length <= 200
 * @return {number[]}
 * 找出 4 条边皆为黑色像素的最大子方阵。
 * 返回一个数组 [r, c, size] ，其中 r, c 分别代表子方阵左上角的行号和列号，size 是子方阵的边长
 * 若 r 相同，返回 c 最小的子方阵
 * 乍一看，这道题和 221. 最大正方形，实则不然。 这道题是可以空心的，只要边长是部分都是 0 即可，这就完全不同了。
 * @summary
 * 倒序遍历获取每个节点的最大黑边行和列
 * 计算最大方阵:
 * 递减size遍历, 如果找到一个符合条件的(右顶点和下顶点的边长度>=cursize)就break
 * 注意一个优化操作是循环遍历到当前res size+1为止, 因为如果等于当前res size的话一定不满足要求了
 * // 垂线dp
 */
var findSquare = function (matrix: number[][]): number[] {
  let res: number[] = []
  const m = matrix.length
  const n = matrix[0].length
  const countDown = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))
  const countRight = Array.from<number, number[]>({ length: m }, () => Array(n).fill(0))

  // 这里可以用python defaultdict的思想简化
  for (let i = m - 1; i >= 0; i--) {
    for (let j = n - 1; j >= 0; j--) {
      if (matrix[i][j] === 0) {
        countDown[i][j] = i + 1 >= m ? 1 : countDown[i + 1][j] + 1
        countRight[i][j] = j + 1 >= n ? 1 : countRight[i][j + 1] + 1
      }
    }
  }

  console.table(countDown)
  console.table(countRight)

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (matrix[i][j] !== 0) continue
      const curSize = Math.min(countDown[i][j], countRight[i][j])
      const maxSize = res[2] || 0
      for (let size = curSize; size > maxSize; size--) {
        console.log(i, j, curSize, maxSize, 6666)
        // 候选大到小找到一个就终止
        // counyDown 保证每列下面都有size个0
        // countRight 保证每行右边都有szie个0
        if (countRight[i + size - 1][j] >= size && countDown[i][j + size - 1] >= size) {
          res = [i, j, size]
          break
        }
      }
    }
  }

  return res[2] ? res : []
}

console.log(
  findSquare([
    [1, 0, 1],
    [0, 0, 1],
    [0, 0, 1],
  ])
)
// 输出: [1,0,2]
// 解释: 输入中 0 代表黑色，1 代表白色，左下角为满足条件的最大子方阵
export {}
