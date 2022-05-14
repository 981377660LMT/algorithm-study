// 找到不相同的字符位置
// 判断是否可以通过一次操作得到
function isOneEditDistance(s: string, t: string): boolean {
  if (Math.abs(s.length - t.length) > 1) return false
  let i = 0
  let j = 0

  while (i < s.length && j < t.length) {
    if (s[i] === t[j]) {
      i++
      j++
    } else {
      // 删除一次/插入一次/替换一次
      return (
        s.slice(i + 1) === t.slice(j + 1) ||
        s.slice(i + 1) === t.slice(j) ||
        s.slice(i) === t.slice(j + 1)
      )
    }
  }

  // 前面全部相等，是否差距最后一个字符
  return Math.abs(s.length - t.length) === 1
}

console.log(isOneEditDistance('ab', 'acb'))
