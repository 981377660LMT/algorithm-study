/**
 * @param {string} s  1 <= s.length <= 25
 * @return {string[]}
 * 删除最小数量的无效括号，使得输入的字符串有效。
 * 返回所有可能的结果
 * @summary 因为求最小数量(最短路径) 所以用bfs
 */
function removeInvalidParentheses(s: string): string[] {
  const visited = new Set<string>()
  let queue: string[] = [s]

  while (queue.length > 0) {
    const valid = queue.filter(isValid)
    if (valid.length > 0) return valid

    const nextQueue: string[] = []
    const len = queue.length
    for (let _ = 0; _ < len; _++) {
      const cur = queue.pop()!
      for (let i = 0; i < cur.length; i++) {
        if (!['(', ')'].includes(cur[i])) continue
        const next = cur.slice(0, i) + cur.slice(i + 1)
        if (visited.has(next)) continue
        visited.add(next)
        nextQueue.push(next)
      }
    }

    queue = nextQueue
  }

  return []

  function isValid(str: string): boolean {
    let level = 0

    for (const char of str) {
      if (char === '(') level++
      else if (char === ')') level--
      if (level < 0) return false
    }

    return level === 0
  }
}

console.log(removeInvalidParentheses('()())()'))

export {}
