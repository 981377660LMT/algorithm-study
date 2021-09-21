class BadMagicDictionary {
  private extendedSet: Set<string>

  constructor() {
    this.extendedSet = new Set()
  }

  /**
   *
   * @param dictionary 1 <= dictionary.length <= 100
   */
  buildDict(dictionary: string[]): void {
    dictionary.forEach(d => {
      for (let i = 0; i < d.length; i++) {
        for (let j = 0; j < 26; j++) {
          const next = d.slice(0, i) + String.fromCodePoint(j + 97) + d.slice(i + 1)
          if (next === d) continue
          this.extendedSet.add(next)
        }
      }
    })
  }

  /**
   *
   * @param searchWord
   * 请判定能否只将这个单词中一个字母换成另一个字母，使得所形成的新单词存在于你构建的字典中。
   */
  search(searchWord: string): boolean {
    return this.extendedSet.has(searchWord)
  }
}

type Length = number
// 优化后的
class GoodMagicDictionary {
  private root: Map<Length, string[]>

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
