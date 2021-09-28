class TrieNode {
  val: string
  isWord: boolean
  children: TrieNode[]
  constructor(val: string) {
    this.val = val
    this.isWord = false
    this.children = []
  }
}

class MagicDictionary {
  private root: TrieNode
  constructor() {
    this.root = new TrieNode('')
  }

  buildDict(dictionary: string[]): void {
    for (const word of dictionary) {
      let root = this.root
      for (const char of word) {
        const key = char.codePointAt(0)! - 97
        if (!root.children[key]) root.children[key] = new TrieNode(char)
        root = root.children[key]
      }
      root.isWord = true
    }
  }

  search(searchWord: string): boolean {
    const arr = searchWord.split('')

    // 对每个位置替换成另外字母 看是否存在这个单词
    for (let i = 0; i < arr.length; i++) {
      for (let ascii = 0; ascii < 26; ascii++) {
        const char = String.fromCodePoint(ascii + 97)
        if (char === arr[i]) continue
        const tmp = arr[i]
        arr[i] = char
        if (this.helper(arr.join(''), this.root)) {
          return true
        }
        arr[i] = tmp
      }
    }

    return false
  }

  private helper(str: string, root: TrieNode): boolean {
    for (const char of str) {
      const key = char.codePointAt(0)! - 97
      if (!root.children[key]) return false
      root = root.children[key]
    }
    return root.isWord
  }
}
const magicDictionary = new MagicDictionary()
magicDictionary.buildDict(['hello', 'leetcode'])
console.log(magicDictionary.search('hello')) // 返回 False
console.log(magicDictionary.search('hhllo')) // 将第二个 'h' 替换为 'e' 可以匹配 "hello" ，所以返回 True

export {}
