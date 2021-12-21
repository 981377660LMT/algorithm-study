/**
 * @param {string} s
 * @return {string[]}
 * @description
 * 返回其通过重新排列组合后所有可能的回文字符串，并去除重复的组合。
 * @summary
 * 根据字符的个数构造回文串即可，每次在首尾两个位置看能否填充两个相同的字符(如果是字符串长是奇数中间位置特判一下)
 */
function generatePalindromes(s: string): string[] {
  const res: string[] = []
  const counter = new Map<string, number>()
  for (const char of s) {
    counter.set(char, (counter.get(char) || 0) + 1)
  }

  dfs(0, s.length - 1, Array(s.length).fill(''))
  return res

  // 因为直接覆盖path 所以回溯不需要重置 path[left] 和 path[right]
  function dfs(left: number, right: number, path: string[]): void {
    if (left > right) {
      res.push(path.join(''))
      return
    }

    for (const [char, frequency] of counter.entries()) {
      if (frequency < 1) continue
      if (left !== right && frequency < 2) continue

      counter.set(char, counter.get(char)! - (left < right ? 2 : 1))
      path[left] = char
      path[right] = char
      dfs(left + 1, right - 1, path)
      counter.set(char, counter.get(char)! + (left < right ? 2 : 1))
    }
  }
}

console.log(generatePalindromes('aabb'))
// 输出: ["abba", "baab"]
