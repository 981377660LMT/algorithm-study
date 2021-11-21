function countLetters(s: string): number {
  let res = 0
  let count = 1

  for (let cur = 0; cur < s.length; cur++) {
    // 'aaa'的情形：利用越界特性
    if (s[cur] !== s[cur + 1]) {
      res += (count * (count + 1)) / 2
      count = 1
    } else {
      count++
    }
  }

  return res
}

console.log(countLetters('aaaba'))
// 输出： 8
// 解释：
// 只含单一字母的子串分别是 "aaa"， "aa"， "a"， "b"。
// "aaa" 出现 1 次。
// "aa" 出现 2 次。
// "a" 出现 4 次。
// "b" 出现 1 次。
// 所以答案是 1 + 2 + 4 + 1 = 8。
