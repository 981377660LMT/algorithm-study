/**
 * @param {string} s  1 <= s.length <= 25
 * @return {string[]}
 * 删除最小数量的无效括号，使得输入的字符串有效。
 * 返回所有可能的结果
 * @summary 因为求最小数量(最短路径) 所以用bfs
 * @description 复杂度n*2^n
 *
 *  n为字符串的长度。一个字符串最多可能有 2^n个子序列
 */
const removeInvalidParentheses = function (s: string): string[] {
  const isValid = (str: string): boolean => {
    let res = 0
    for (const char of str) {
      if (char === '(') res++
      else if (char === ')') res--
      if (res < 0) return false
    }
    return res === 0
  }

  const visited = new Set<string>()
  const queue: string[] = [s]

  while (queue.length) {
    const valid = queue.filter(isValid)
    if (valid.length) return valid

    const len = queue.length
    for (let i = 0; i < len; i++) {
      const item = queue[i]
      // 尝试去掉每一个括号的位置
      for (let i = 0; i < item.length; i++) {
        if (['(', ')'].includes(item[i])) {
          const next = item.slice(0, i) + item.slice(i + 1)
          if (!visited.has(next)) {
            queue.push(next)
            visited.add(next)
          }
        }
      }
    }
  }

  return []
}

console.log(removeInvalidParentheses('()())()'))

export {}
