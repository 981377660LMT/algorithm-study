// 类似于395. 至少有 K 个重复字符的最长子串 逐步排除不合法
// 当一个字符串 s 包含的每一种字母的大写和小写形式 同时 出现在 s 中，就称这个字符串 s 是 美好 字符串。
// 请你返回 s 最长的 美好子字符串 。如果有多个答案，请你返回 最早 出现的一个。如果不存在美好子字符串，请你返回一个空字符串。
function longestNiceSubstring(s: string): string {
  if (s.length === 0) return s
  const set = new Set(s)

  for (const char of set) {
    if (!set.has(swapcase(char)))
      return s
        .split(char)
        .map(longestNiceSubstring)
        .sort((a, b) => b.length - a.length)[0]
  }

  return s
}

function swapcase(char: string): string {
  const codePoint = char.codePointAt(0)!
  if (codePoint >= 65 && codePoint <= 90) return String.fromCodePoint(codePoint + 32)
  if (codePoint >= 97 && codePoint <= 122) return String.fromCodePoint(codePoint - 32)
  return char
}

console.log(longestNiceSubstring('YazaAay'))
// 输出："aAa"
