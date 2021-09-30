// 找出并返回所有从 beginWord 到 endWord 的 所有最短转换序列
// wordList 中的所有单词 互不相同

// BFS 可以帮助我们找到最短路径是多远，但是不能帮助我们得到那些节点是在最短路径上面。
// 而通过 BFS 遍历得到的信息进行 DFS 遍历回溯，可以帮助我们找到那些在最短路径上的节点。
const findLadders = (beginWord: string, endWord: string, wordList: string[]): string[][] => {
  const res: string[][] = []
  const wordSet = new Set(wordList)
  wordSet.delete(beginWord)
  if (!wordSet.has(endWord)) return res
  const { adjMap, levelMap } = bfs()
  dfs(beginWord, [beginWord], adjMap, levelMap)
  return res

  function bfs() {
    // 邻接表,为dfs做准备
    const adjMap = new Map<string, Set<string>>()
    // 距离表，相当于记录每个点的level
    // 树不需要是因为树的邻接表默认level+1，而图需要
    const levelMap = new Map<string, number>([[beginWord, 0]])
    const visited = new Set([beginWord])
    const queue: [string, number][] = [[beginWord, 0]]
    // bfs 操作
    while (queue.length) {
      const [word, dis] = queue.shift()!
      if (word === endWord) break

      for (let i = 0; i < word.length; i++) {
        for (let j = 0; j < 26; j++) {
          const next = word.slice(0, i) + String.fromCharCode(97 + j) + word.slice(i + 1)
          if (!wordSet.has(next) || word === next) continue // 不是单词表中的单词就忽略

          !adjMap.has(word) && adjMap.set(word, new Set())
          adjMap.get(word)!.add(next)

          if (visited.has(next)) continue
          levelMap.set(next, dis + 1)
          queue.push([next, dis + 1])
          visited.add(next)
        }
      }
    }

    return {
      adjMap,
      levelMap,
    }
  }

  // 通过levelMap可以保证dfs是沿着bfs走的最短路径
  function dfs(
    curWord: string,
    path: string[],
    adjMap: Map<string, Set<string>>,
    levelMap: Map<string, number>
  ) {
    if (curWord === endWord) {
      return res.push(path.slice())
    }

    if (adjMap.has(curWord)) {
      for (const next of adjMap.get(curWord)!) {
        if (levelMap.get(next) === levelMap.get(curWord)! + 1) {
          path.push(next)
          dfs(next, path, adjMap, levelMap)
          path.pop()
        }
      }
    }
  }
}

//
console.dir(findLadders('hot', 'dog', ['hot', 'dog', 'dot']), { depth: null })
// console.dir(findLadders('hot', 'dog', ['hot', 'dog']), { depth: null })

export {}
// 我们在 BFS 中，
// 就把每个节点的所有相邻节点保存到 adjMap 中，
// 就省去了 DFS 再去找相邻节点的时间
// 在 DFS 中，我们每次都根据节点的层数levelMap来进行深搜
