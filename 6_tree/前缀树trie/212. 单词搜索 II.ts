/**
 * @param {string[][]} board
 * @param {string[]} words
 * @return {string[]}
 * @description 找出所有同时在二维网格和字典中出现的单词。
 * @summary 我们需要对矩阵中每一项都进行深度优先遍历（DFS）。 递归的终点是
 * 超出边界
   递归路径上组成的单词不在 words 的前缀。
 */
const findWords = function (board: string[][], words: string[]) {}

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
