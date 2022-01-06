type WordLength = number

// 优化后的
class GoodMagicDictionary {
  private root: Map<WordLength, string[]>

  constructor() {
    this.root = new Map()
  }

  /**
   *
   * @param dictionary 1 <= dictionary.length <= 100
   */
  buildDict(dictionary: string[]): void {
    dictionary.forEach(word => {
      const key = word.length
      !this.root.has(key) && this.root.set(key, [])
      this.root.get(key)!.push(word)
    })
  }

  /**
   *
   * @param searchWord
   * 请判定能否只将这个单词中一个字母换成另一个字母，使得所形成的新单词存在于你构建的字典中。
   */
  search(searchWord: string): boolean {
    const key = searchWord.length
    if (!this.root.has(key)) return false
    return this.root.get(key)!.some(word => {
      let diff = 0
      for (let i = 0; i < word.length; i++) {
        if (word[i] !== searchWord[i]) diff++
        if (diff > 1) return false
      }
      return diff === 1
    })
  }
}

const magicDictionary = new GoodMagicDictionary()
magicDictionary.buildDict(['hello', 'leetcode'])
console.log(magicDictionary.search('hello')) // 返回 False
console.log(magicDictionary.search('hhllo')) // 将第二个 'h' 替换为 'e' 可以匹配 "hello" ，所以返回 True
