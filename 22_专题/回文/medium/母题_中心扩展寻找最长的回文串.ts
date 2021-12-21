// 寻找l到r的最长回文子串(连续的)
const maxExpand = (s: string, l: number, r: number): string => {
  while (l >= 0 && r < s.length && s[l] === s[r]) {
    l--
    r++
  }
  return s.slice(l + 1, r)
}

if (require.main === module) {
  console.log(maxExpand('abbba', 2, 2))
  console.log(maxExpand('abccbd', 2, 3))
}

export { maxExpand }
