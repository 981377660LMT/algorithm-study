import { Trie } from './1_实现trie'

class WordTrie extends Trie {
  searchWord(word: string) {
    let rootP = this.root
    let res = ''

    for (const letter of word) {
      const next = rootP.children.get(letter)
      if (!next) return word
      res += letter
      rootP = next
      if (rootP.isWord) return res
    }

    return word
  }
}

/**
 * @param {string[]} dictionary
 * @param {string} sentence
 * @return {string}
 * @description 你需要将句子中的所有继承词用词根替换掉。如果继承词有许多可以形成它的词根，则用最短的词根替换它。
 */
const replaceWords = function (dictionary: string[], sentence: string): string {
  const res: string[] = []
  const trie = new WordTrie()

  dictionary.forEach(w => trie.insert(w))
  sentence.split(/\s+/g).forEach(w => res.push(trie.searchWord(w)))

  return res.join(' ')
}

console.log(replaceWords(['cat', 'bat', 'rat'], 'the cattle was rattled by the battery'))
// "the cat was rat by the bat"
