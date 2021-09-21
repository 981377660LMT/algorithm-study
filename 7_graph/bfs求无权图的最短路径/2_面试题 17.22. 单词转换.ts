function findLadders(beginWord: string, endWord: string, wordList: string[]): string[] {
  const res: string[] = []
  const wordSet = new Set(wordList)
  if (!wordSet.has(endWord)) return res
  const { adjMap, levelMap } = bfs()

  return dfs(beginWord, [beginWord]).next().value || []

  function bfs() {
    // 邻接表,为dfs做准备
    const adjMap = new Map<string, Set<string>>()
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

  // 通过level可以保证dfs是沿着bfs走的最短路径
  function* dfs(curWord: string, path: string[]): Generator<string[]> {
    if (curWord === endWord) {
      yield path
    }

    if (adjMap.has(curWord)) {
      for (const next of adjMap.get(curWord)!) {
        if (levelMap.get(next) === levelMap.get(curWord)! + 1) {
          path.push(next)
          yield* dfs(next, path)
          path.pop()
        }
      }
    }
  }
}

console.log(findLadders('hit', 'cog', ['hot', 'dot', 'dog', 'lot', 'log', 'cog']))
