// 找出并返回所有从 beginWord 到 endWord 的 最短距离
// wordList 中的所有单词 互不相同

// 先用 BFS 求出最短距离
// 再用 DFS 求出最短距离路径
const findLadders = (beginWord: string, endWord: string, wordList: string[]): number => {
  // 充当visited的作用
  const wordSet = new Set<string>(wordList)
  if (!wordSet.has(endWord)) return 0
  const queue: [string, number][] = [[beginWord, 1]]

  while (queue.length) {
    const [word, dis] = queue.shift()!
    if (word === endWord) return dis
    for (let i = 0; i < word.length; i++) {
      for (let j = 0; j < 26; j++) {
        const target = word.slice(0, i) + String.fromCharCode(97 + j) + word.slice(i + 1)
        if (wordSet.has(target)) {
          queue.push([target, dis + 1])
          wordSet.delete(target)
        }
      }
    }
  }

  return 0
}

console.dir(findLadders('hit', 'cog', ['hot', 'dot', 'dog', 'lot', 'log', 'cog']), { depth: null })

export {}
