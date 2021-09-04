function nearestPalindromic(n: string): string {
  const num = Number(n)
  const reverseStr = (s: string) => s.split('').reverse().join('')
  if (num <= 10 || parseInt(reverseStr(n)) === 1) return (num - 1).toString()
  if (n === '11') return '9'
  //   奇数长度 3 分为 2 和 1
  //   偶数长度 4 分为 2 和 2
  const pre = n.slice(0, (n.length + 1) >> 1)
  const tmp = [Number(pre) - 1, Number(pre), Number(pre) + 1].map(String)
  // 这里有点问题 思路大致这样了
  const candidates = tmp.map(str => str + reverseStr(str).slice(n.length % 2))

  let min = Infinity
  let res = ''
  console.log(candidates)
  for (const candidate of candidates) {
    const diff = Math.abs(Number(candidate) - num)
    if (diff === 0) continue
    if (diff < min) {
      min = diff
      res = candidate
    }
  }
  return res
}

console.log(nearestPalindromic('123'))
