import { BinaryTree } from '../../Tree'

/**
 *
 * @param s
 * @returns
 * 输入字符串中只包含 '(', ')', '-'(负号) 和 '0' ~ '9'
 * 空树由 "" 而非"()"表示。
 */
function str2tree(s: string): BinaryTree | null {
  if (!s) return null

  return dfs(s)

  function dfs(str: string): BinaryTree | null {
    if (!str) return null
    const firstBracket = str.indexOf('(')
    const rootVal = firstBracket === -1 ? str : str.slice(0, firstBracket)
    const [left, right] = extractBracket(str)
    const root = new BinaryTree(Number(rootVal))
    root.left = dfs(left.slice(1, -1))
    root.right = dfs(right.slice(1, -1))
    return root
  }

  function extractBracket(str: string) {
    const res: number[] = []
    const bracket = new Set(['(', ')'])
    let level = 0

    for (let i = 0; i < str.length; i++) {
      if (!bracket.has(str[i])) continue

      if (level === 0) res.push(i)

      if (str[i] === '(') level++
      else level--

      if (level === 0) res.push(i)
    }

    const left = str.slice(res[0] || 0, res[1] + 1 || 0)
    const right = str.slice(res[2] || 0, res[3] + 1 || 0)

    return [left, right]
  }
}

console.log(str2tree('4(2(3)(1))(6(5))'))
