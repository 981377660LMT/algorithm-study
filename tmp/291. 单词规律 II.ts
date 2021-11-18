// pattern 里的每个`字母`和字符串 str 中每个 非空 `单词`之间，
// 是否存在着 双射 的对应规律
// https://leetcode-cn.com/problems/word-pattern-ii/comments/711636
function wordPatternMatch(pattern: string, s: string): boolean {
  const charToWord = new Map<string, string>()
  const visitedWord = new Set<string>()

  return bt(pattern, s)

  function bt(patterns: string, words: string): boolean {
    if (patterns.length === 0) return words.length === 0
    const curChar = patterns[0]
    // 之前可能添加过curChar映射的单词
    const mappingWordOfCurChar = charToWord.get(curChar)

    for (let i = 1; i <= words.length; i++) {
      const curWord = words.slice(0, i)

      // 同时不存在,curWord和curChar都没用过
      if (mappingWordOfCurChar == undefined && !visitedWord.has(curWord)) {
        charToWord.set(curChar, curWord)
        visitedWord.add(curWord)
        if (bt(patterns.slice(1), words.slice(i))) return true
        charToWord.delete(curChar)
        visitedWord.delete(curWord)
      } else if (mappingWordOfCurChar === curWord) {
        // 同时存在，双射
        if (bt(patterns.slice(1), words.slice(i))) return true
      }
    }

    return false
  }
}

console.log(wordPatternMatch('abab', 'redblueredblue'))
// 输出：true
// 解释：一种可能的映射如下：
// 'a' -> "red"
// 'b' -> "blue"
