class ValidWordAbbr {
  private record: Map<string, Set<string>>
  constructor(dictionary: string[]) {
    this.record = new Map()
    for (const word of dictionary) {
      if (word.length <= 2) continue
      const abbr = this.getAbbr(word)
      !this.record.has(abbr) && this.record.set(abbr, new Set())
      this.record.get(abbr)!.add(word)
    }
  }

  isUnique(word: string): boolean {
    const n = word.length
    if (n <= 2) return true
    const abbr = this.getAbbr(word)
    if (!this.record.has(abbr)) return true // 没有任何其他单词的 缩写 与该单词 word 的 缩写 相同
    return this.record.get(abbr)!.size === 1 && this.record.get(abbr)!.has(word) // 所有 缩写 与该单词 word 的 缩写 相同的单词都与 word 相同
  }

  private getAbbr(word: string) {
    if (word.length <= 2) return word
    return `${word[0]}${word.length - 2}${word[word.length - 1]}`
  }
}

export {}
// 单词的 缩写 需要遵循 <起始字母><中间字母数><结尾字母> 这样的格式
// dog --> d1g 因为第一个字母 'd' 和最后一个字母 'g' 之间有 1 个字母
// internationalization --> i18n 因为第一个字母 'i' 和最后一个字母 'n' 之间有 18 个字母
// it --> it 单词只有两个字符，它就是它自身的 缩写

// 如果满足下述任意一个条件，返回 true
// 字典 dictionary 中没有任何其他单词的 缩写 与该单词 word 的 缩写 相同
// 字典 dictionary 中的所有 缩写 与该单词 word 的 缩写 相同的单词都与 word 相同
