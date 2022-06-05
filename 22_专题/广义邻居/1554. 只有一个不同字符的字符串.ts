// 给定一个字符串列表 dict ，其中所有字符串的长度都相同。
// 当存在两个字符串在相同索引处只有一个字符不同时，返回 True ，否则返回 False 。
// dict 中的字符数小于或等于 10^5 。
// 你可以以 O(n*m) 的复杂度解决问题吗？其中 n 是列表 dict 的长度，m 是字符串的长度。

//思路： 檢查每個位置i=0~N-1，把所有字串的第i位抽掉放進 set，若 set 裡的個數比原先字典少，表示有重複
function differByOne(dict: string[]): boolean {
  const visited = new Set<string>()

  for (const word of dict) {
    for (let i = 0; i < word.length; i++) {
      const slice = word.slice(0, i) + '*' + word.slice(i + 1)
      if (visited.has(slice)) return true
      visited.add(slice)
    }
  }

  return false
}

console.log(differByOne(['abcd', 'acbd', 'aacd']))
// 输出：true
// 解释：字符串 "abcd" 和 "aacd" 只在索引 1 处有一个不同的字符。
