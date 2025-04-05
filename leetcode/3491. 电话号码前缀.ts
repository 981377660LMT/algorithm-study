// 给定一个字符串数组 numbers 表示电话号码。
// 如果没有电话号码是任何其他电话号码的前缀，则返回 true；否则，返回 false。

export {}

class TrieNode {
  children: Map<string, TrieNode> = new Map()
  isEnd = false
  preCount = 0
}

function phonePrefix(numbers: string[]): boolean {
  const root = new TrieNode()
  for (const word of numbers) {
    let cur = root
    for (const c of word) {
      if (!cur.children.has(c)) {
        cur.children.set(c, new TrieNode())
      }
      cur = cur.children.get(c)!
      if (cur.isEnd) return false
      cur.preCount++
    }
    if (cur.preCount > 1) return false
    cur.isEnd = true
  }
  return true
}
