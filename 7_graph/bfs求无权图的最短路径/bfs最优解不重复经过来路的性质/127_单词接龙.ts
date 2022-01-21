// 找出并返回所有从 beginWord 到 endWord 的 最短距离
// wordList 中的所有单词 互不相同

const ladderLength = (beginWord: string, endWord: string, wordList: string[]): number => {
  // 充当visited的作用
  const wordSet = new Set(wordList)
  if (!wordSet.has(endWord)) return 0
  const queue: [string, number][] = [[beginWord, 1]]

  while (queue.length > 0) {
    const [word, dis] = queue.shift()!
    if (word === endWord) return dis
    for (const next of getNextWords(word, wordSet)) {
      queue.push([next, dis + 1])
      wordSet.delete(next) // 这里可以删除的原因是 之后的才是接近答案的(层数更接近答案)
    }
  }

  return 0

  function* getNextWords(curWord: string, wordSet: Set<string>) {
    for (let i = 0; i < curWord.length; i++) {
      for (let j = 0; j < 26; j++) {
        const next = curWord.slice(0, i) + String.fromCharCode(97 + j) + curWord.slice(i + 1)
        if (wordSet.has(next) && curWord !== next) yield next
      }
    }
  }
}

console.dir(ladderLength('hit', 'cog', ['hot', 'dot', 'dog', 'lot', 'log', 'cog']), { depth: null })

export {}
