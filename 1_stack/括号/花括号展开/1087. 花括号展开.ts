// # 如果你熟悉 Shell 编程，那么一定了解过花括号展开，它可以用来生成任意字符串。

import { isalpha } from '../../../0_字符串/string'

// # 1 <= S.length <= 50
// # 你可以假设题目中不存在嵌套的花括号

// print(Solution().expand("{a,b}c{d,e}f"))
// # 输出：["acdf","acef","bcdf","bcef"]
function expand(s: string): string[] {
  const groups: string[][] = []
  extractWordGroups()

  const res: string[] = []
  bt(0, [])

  return res.sort()

  function extractWordGroups(): void {
    let isExpandMode = false
    let wordBuffer: string[] = []

    for (const char of s) {
      if (char === '{') {
        isExpandMode = true
      } else if (char === '}') {
        groups.push(wordBuffer)
        wordBuffer = []
        isExpandMode = false
      } else if (isExpandMode) {
        if (isalpha(char)) wordBuffer.push(char)
      } else {
        if (isalpha(char)) groups.push([char])
      }
    }
  }

  function bt(index: number, path: string[]): void {
    if (index === groups.length) {
      res.push(path.join(''))
      return
    }

    for (const char of groups[index]) {
      path.push(char)
      bt(index + 1, path)
      path.pop()
    }
  }
}

console.log(expand('{a,b}c{d,e}f'))
