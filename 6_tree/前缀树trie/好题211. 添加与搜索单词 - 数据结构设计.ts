import { Trie, TrieNode } from './实现trie/1_实现trie'

class WordTrie extends Trie {
  // 因为需要匹配'.'，需要递归遍历节点
  override search(word: string): boolean {
    return this.match(this.root, word, 0)
  }

  // 用index比每次slice要好
  match(node: TrieNode, word: string, index: number): boolean {
    // 递归终点
    if (index === word.length) return node.isWord

    if (word[index] !== '.') {
      const next = node.children.get(word[index])
      if (!next) return false
      return this.match(node.children.get(word[index])!, word, index + 1)
    } else {
      // '.'时遍历所有孩子 如果一个true则为true
      for (const next of node.children.values()) {
        if (this.match(next, word, index + 1)) return true
      }
      return false
    }
  }
}

// word 中可能包含一些 '.' ，每个 . 都可以表示任何一个字母。
class WordDictionary {
  constructor(private trie: WordTrie) {}

  addWord(word: string): void {
    this.trie.insert(word)
  }

  search(word: string): boolean {
    return this.trie.search(word)
  }
}

const wordDict = new WordDictionary(new WordTrie())

wordDict.addWord('dad')
console.log(wordDict.search('pad'))
console.log(wordDict.search('.ad'))

export {}
