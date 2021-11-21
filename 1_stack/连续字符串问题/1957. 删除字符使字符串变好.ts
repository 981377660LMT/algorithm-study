// 一个字符串如果没有 三个连续 相同字符，那么它就是一个 好字符串 。
// 给你一个字符串 s ，请你从 s 删除 最少 的字符，使它变成一个 好字符串 。
function makeFancyString(s: string): string {
  const sb: string[] = []
  let sameCount = 0

  // 1. 更新sameCount
  // 2. 根据sameCount决定要不要加入sb
  for (const char of s) {
    if (sb.length > 0 && char === sb[sb.length - 1]) sameCount++
    else sameCount = 1
    if (sameCount <= 2) sb.push(char)
  }

  return sb.join('')
}

console.log(makeFancyString('leeetcode'))
