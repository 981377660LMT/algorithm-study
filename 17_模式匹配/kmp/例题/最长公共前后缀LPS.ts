/**
 * @param {string} s
 * @return {string}
 * @description 找出在原字符串中既是 非空 前缀也是后缀（不包括原字符串自身）的最长字符串
 * O(n) 求最长公共前缀（LPS）(Longest Proper String) 是KMP的一个步骤
 * 关键是构造出lps数组
 */
function longestPrefix(s: string): string {
  const lps = getNext(s)
  return s.slice(0, lps[s.length - 1])
}

// 求next数组
function getNext(needle: string): number[] {
  // lps[i]表示[0,i]这一段字符串中lps的长度
  const next = Array<number>(needle.length).fill(0)
  let j = 0

  for (let i = 1; i < needle.length; i++) {
    while (j > 0 && needle[i] !== needle[j]) {
      //  前进到最长公共前缀结尾处
      j = next[j - 1]
    }
    if (needle[i] === needle[j]) j++
    next[i] = j
  }

  return next
}

if (require.main === module) {
  console.log(getNext('abcdabgabca'))
  console.log(longestPrefix('ababab'))
  // "abab"
}
