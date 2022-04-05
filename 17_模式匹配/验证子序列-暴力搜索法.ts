/**
 *
 * @param s 源字符串
 * @param t 待搜索的字符串
 */
const isSubsequence = (s: string, t: string) => {
  if (s.length < t.length) return -1

  // 看[i,i+t.length-1]是否与t匹配
  for (let i = 0; i + t.length - 1 < s.length; i++) {
    let j = 0
    for (; j < t.length; j++) {
      if (s.charAt(i + j) !== t.charAt(j)) break
    }
    if (j === t.length) return i
  }

  return -1
}

console.log(isSubsequence('asdfghj', 'dfg'))

export {}
