// 输入： s = "abcxyz123", words = ["abc","123"]
// 输出："<b>abc</b>xyz<b>123</b>"
function addBoldTag(s: string, words: string[]): string {
  const n = s.length
  const boldChars = Array<boolean>(n).fill(false)

  // 找到需要bold的字符
  for (const word of words) {
    let start = 0
    while (true) {
      const position = s.indexOf(word, start)
      if (position === -1) break
      for (let i = 0; i < word.length; i++) {
        boldChars[position + i] = true
      }
      start = position + 1
    }
  }

  console.log(boldChars)
  const sb: string[] = []
  let i = 0
  while (i < n) {
    // 找到下一个bold位置
    while (i < n && !boldChars[i]) {
      sb.push(s[i])
      i++
    }

    if (i === n) break

    sb.push('<b>')
    while (i < n && boldChars[i]) {
      sb.push(s[i])
      i++
    }
    sb.push('</b>')
  }

  return sb.join('')
}

console.log(addBoldTag('aaabbcc', ['aaa', 'aab', 'bc']))
// 输出："<b>aaabbc</b>c"

console.log(addBoldTag('abcxyz123', ['abc', '123']))
// 输出："<b>abc</b>xyz<b>123</b>"
