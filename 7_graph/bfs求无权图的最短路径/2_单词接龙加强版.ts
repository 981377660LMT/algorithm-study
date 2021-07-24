// 找出并返回所有从 beginWord 到 endWord 的 所有最短转换序列
// wordList 中的所有单词 互不相同

// 先用 BFS 求出最短距离以及构建整个图的邻接表(注意广度优先走完必定是最短路径，因为想继续向下走不可能，已经被上面的节点封住了)
// 再用 DFS+回溯 根据邻接表(wordMap)和层数表(levelMap) 求出所有最短距离路径
const findLadders = (beginWord: string, endWord: string, wordList: string[]): string[][] => {
  const res: string[][] = []
  // 将wordList存储为Set，方便快速判断新单词是否在wordList内
  const wordSet = new Set<string>(wordList)
  wordSet.delete(beginWord)
  if (!wordSet.has(endWord)) return res
  const visited = new Set<string>(beginWord)
  const queue: [string, number][] = [[beginWord, 0]]
  // 邻接表,为dfs做准备
  const wordMap = new Map<string, Set<string>>()
  // 距离表，相当于记录每个点的level
  // 树不需要是因为树的邻接表默认level+1，而图需要
  const levelMap = new Map<string, number>([[beginWord, 0]])
  // 加速计算
  let canReach = false

  // bfs 操作
  while (queue.length) {
    const [word, dis] = queue.shift()!
    if (word === endWord) canReach = true

    for (let i = 0; i < word.length; i++) {
      for (let j = 0; j < 26; j++) {
        const next = word.slice(0, i) + String.fromCharCode(97 + j) + word.slice(i + 1)
        if (!wordSet.has(next)) continue // 不是单词表中的单词就忽略

        if (!wordMap.has(word)) {
          wordMap.set(word, new Set([next]))
        } else {
          wordMap.get(word)!.add(next)
        }

        if (visited.has(next)) continue
        levelMap.set(next, dis + 1)
        queue.push([next, dis + 1])
        visited.add(next)
      }
    }
  }

  if (!canReach) return []

  // 注意这里dfs不能用visited因为是回溯，要经过一个节点多次，通过level可以保证不循环
  const dfs = (beginWord: string, tmp: string[]) => {
    // console.log(beginWord, tmp)
    if (beginWord === endWord) {
      return res.push(tmp)
    }

    for (const next of wordMap.get(beginWord)!) {
      if (levelMap.get(next) === levelMap.get(beginWord)! + 1) {
        dfs(next, [...tmp, next])
      }
    }
  }
  dfs(beginWord, [beginWord])

  // console.log(levelMap, wordMap)

  return res
}

console.dir(findLadders('hot', 'dog', ['hot', 'dog', 'dot']), { depth: null })

export {}
