import { Trie } from './实现trie/1_实现trie'

/**
 * @param {string[][]} board
 * @param {string[]} words
 * @return {string[]}
 * @description 找出所有同时在二维网格和字典中出现的单词。
 * @summary 我们需要对矩阵中每一项都进行深度优先遍历（DFS）。 递归的终点是
 * 超出边界
   递归路径上组成的单词不在 words 的前缀。
 */
const findWords = function (board: string[][], words: string[]) {
  const trie = new Trie()
  words.forEach(word => trie.insert(word))

  const r = board.length
  const c = board[0].length
  const next = [
    [-1, 0],
    [0, 1],
    [1, 0],
    [0, -1],
  ]

  const bt = (x: number, y: number, cur: string, res: Set<string>) => {
    // 1. 回溯终点
    if (trie.search(cur)) res.add(cur)

    if (!trie.startsWith(cur)) return

    // 2.回溯处理
    // 标记为visited
    const tmp = board[x][y]
    board[x][y] = '$'
    for (const [dx, dy] of next) {
      const nextRow = x + dx
      const nextColumn = y + dy
      // 在矩阵中
      if (nextRow >= 0 && nextRow < r && nextColumn >= 0 && nextColumn < c) {
        bt(nextRow, nextColumn, cur + board[nextRow][nextColumn], res)
      }
    }

    // 3. 回溯重置
    board[x][y] = tmp
  }

  // 4.每个点开始回溯
  const res = new Set<string>()
  for (let i = 0; i < r; i++) {
    for (let j = 0; j < c; j++) {
      bt(i, j, board[i][j], res)
    }
  }

  return [...res]
}

console.log(
  findWords(
    [
      ['o', 'a', 'a', 'n'],
      ['e', 't', 'a', 'e'],
      ['i', 'h', 'k', 'r'],
      ['i', 'f', 'l', 'v'],
    ],
    ['oath', 'pea', 'eat', 'rain']
  )
)

console.log(
  findWords(
    [
      ['o', 'a', 'b', 'n'],
      ['o', 't', 'a', 'e'],
      ['a', 'h', 'k', 'r'],
      ['a', 'f', 'l', 'v'],
    ],
    ['oa', 'oaa']
  )
)
