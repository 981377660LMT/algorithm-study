function longestDecomposition(text: string): number {
  let res = 0
  let l = 1

  while (l < text.length) {
    const sub1 = text.slice(0, l)
    const sub2 = text.slice(text.length - l)
    if (sub1 === sub2) {
      res += 2
      text = text.slice(l, text.length - l)
      l = 1
    } else {
      l++
    }
  }

  return text ? res + 1 : res
}

console.log(longestDecomposition('ghiabcdefhelloadamhelloabcdefghi'))
// 输出：7
// 解释：我们可以把字符串拆分成 "(ghi)(abcdef)(hello)(adam)(hello)(abcdef)(ghi)"。
