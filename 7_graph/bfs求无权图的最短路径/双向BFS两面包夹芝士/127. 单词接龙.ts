import { biBfs } from './双向bfs'

function ladderLength(beginWord: string, endWord: string, wordList: string[]): number {
  const wordSet = new Set(wordList)
  if (!wordSet.has(endWord)) return 0

  const getNextState = (cur: string) => {
    const res: string[] = []

    for (let i = 0; i < cur.length; i++) {
      for (let j = 0; j < 26; j++) {
        const next = cur.slice(0, i) + String.fromCharCode(97 + j) + cur.slice(i + 1)
        if (wordSet.has(next) && cur !== next) res.push(next)
      }
    }

    return res
  }

  return biBfs(beginWord, endWord, getNextState) + 1
}

console.log(ladderLength('hit', 'cog', ['hot', 'dot', 'dog', 'lot', 'log', 'cog']))
