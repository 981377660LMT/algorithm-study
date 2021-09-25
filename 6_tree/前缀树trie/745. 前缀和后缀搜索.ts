class WordFilter {
  constructor(words: string[]) {}

  /**
   *
   * @param prefix
   * @param suffix
   * 返回词典中具有前缀 prefix 和后缀suffix 的单词的下标。
   * 如果存在不止一个满足要求的下标，返回其中 最大的下标 。
   * 如果不存在这样的单词，返回 -1 。
   */
  f(prefix: string, suffix: string): number {}

  static main() {
    const wordFilter = new WordFilter(['apple'])
    console.log(wordFilter.f('a', 'e')) // 0
  }
}

export {}
