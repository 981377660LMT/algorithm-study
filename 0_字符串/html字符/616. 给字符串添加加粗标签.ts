// 输入： s = "abcxyz123", words = ["abc","123"]
// 输出："<b>abc</b>xyz<b>123</b>"
function addBoldTag(s: string, words: string[]): string {
  const n = s.length
  const boldChar = Array<boolean>(n).fill(false)

  for (const word of words) {
    let start = 0
    while (true) {
      const position = s.indexOf(word, start)
      if (position === -1) break
      for (let i = 0; i < word.length; i++) {
        boldChar[position + i] = true
      }
      start = position + 1
    }
  }

  console.log(boldChar)
  const sb: string[] = []
  let i = 0
  while (i < n) {
    while (i < n && !boldChar[i]) {
      sb.push(s[i])
      i++
    }

    if (i === n) break

    sb.push('<b>')
    while (i < n && boldChar[i]) {
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
