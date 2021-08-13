/**
 * @param {string} s
 * @return {string}
 * 给你一个字符串 s ，请你去除字符串中重复的字母，使得每个字母只出现一次。需保证 返回结果的字典序最小（要求不能打乱其他字符的相对位置）。
 * 对于每一个字符，如果其对应的剩余出现次数大于 1，我们可以选择丢弃（也可以选择不丢弃），否则不可以丢弃。
 */
const removeDuplicateLetters = function (s: string): string {
  const n = s.length
  // 用集合取代数组的includes
  const visited = new Set<string>()
  const stack: string[] = []
  const remainCounter = new Map<string, number>()

  for (const letter of s) {
    remainCounter.set(letter, (remainCounter.get(letter) || 0) + 1)
  }

  for (let i = 0; i < n; i++) {
    if (!visited.has(s[i])) {
      while (
        stack.length &&
        stack[stack.length - 1] > s[i] &&
        // 再丢弃就没了
        remainCounter.get(stack[stack.length - 1])
      ) {
        const tmp = stack.pop()!
        visited.delete(tmp)
      }
      stack.push(s[i])
      visited.add(s[i])
    }
    remainCounter.set(s[i], remainCounter.get(s[i])! - 1)
  }

  return stack.join('')
}

console.log(removeDuplicateLetters('bcabc'))
console.log(removeDuplicateLetters('cbacdcbc'))
// "acdb"
