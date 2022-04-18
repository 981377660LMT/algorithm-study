// 如果一个单词 word 符合下面两个条件，那么它就可以算作谜底：
// 单词 word 中包含谜面 puzzle 的第一个字母。
// 单词 word 中的每一个字母都可以在谜面 puzzle 中找到。
// 返回一个答案数组 answer，数组中的每个元素 answer[i] 是在给出的单词列表 words 中可以作为字谜迷面 puzzles[i] 所对应的谜底的单词数目。

// puzzles[i].length == 7
// 1 <= words.length <= 10^5
// 1 <= puzzles.length <= 10^4
function findNumOfValidWords(words: string[], puzzles: string[]): number[] {
  const res: number[] = []
  const wordCounter = new Map<number, number>()
  words.forEach(word => {
    const zippedWord = getZippedWord(word)
    wordCounter.set(zippedWord, (wordCounter.get(zippedWord) || 0) + 1)
  })

  for (const puzzle of puzzles) {
    let count = 0
    for (const subset of getSubset(puzzle.slice(1))) {
      let mask = 1 << (puzzle[0].codePointAt(0)! - 97)
      for (const char of subset) {
        mask |= 1 << (char.codePointAt(0)! - 97)
      }
      count += wordCounter.get(mask) || 0
    }
    res.push(count)
  }

  return res

  function getZippedWord(word: string): number {
    let mask = 0
    for (const char of word) {
      mask |= 1 << (char.codePointAt(0)! - 97)
    }
    return mask
  }

  function getSubset(arr: string): string[][] {
    const res: string[][] = []
    const n = 1 << arr.length

    for (let i = 0; i < n; i++) {
      const tmp: string[] = []
      for (let j = 0; j < arr.length; j++) {
        if (i & (1 << j)) tmp.push(arr[j])
      }
      res.push(tmp)
    }

    return res
  }
}

console.log(
  findNumOfValidWords(
    ['aaaa', 'asas', 'able', 'ability', 'actt', 'actor', 'access'],
    ['aboveyz', 'abrodyz', 'abslute', 'absoryz', 'actresz', 'gaswxyz']
  )
)
// 对于每个 puzzle 没有必要遍历所以 words，只用找符合条件的 words 出现了多少次就行了(遍历小的)
// 1.只关注是否出现=>或运算压缩单词
// 26位int  int中的每一位取0和1表示字符是否出现过
// "aabb" 可以用 11 表示，"accc" 可以用 101 表示
// 不同的单词可能映射成同一个数字，比如 "aabbb" 和 "ab" 都映射成了 11。这就是状态压缩。
// 2.匹配
// word 状态压缩后的数字 和 puzzle[0] + subset(puzzle[1:N - 1]) 状态压缩后的数字相等
// 求出puzzle[0] + subset(puzzle[1:N - 1]) 对应的二进制数字之后，累加 hashmap 中该二进制数字出现的次数，就是该 puzzle 对应的 word 有多少。
