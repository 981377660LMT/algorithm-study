// !注意会2000*2000会爆空间 大空间需要将数组压缩成数字 或者不要使用javascript

const COL = 1e6
function maxPalindromes(s: string, k: number): number {
  const n = s.length
  const intervals: number[] = []
  // 中心扩展法求出所有回文串
  const expand = (left: number, right: number) => {
    while (left >= 0 && right < n && s[left] === s[right]) {
      if (right - left + 1 >= k) {
        intervals.push(left * COL + right)
      }
      left--
      right++
    }
  }

  for (let i = 0; i < n; i++) {
    expand(i, i)
    expand(i, i + 1)
  }

  intervals.sort((a, b) => (a % COL) - (b % COL))
  let res = 0
  let preEnd = -1
  for (const num of intervals) {
    const [start, end] = [Math.floor(num / COL), num % COL]
    if (start > preEnd) {
      res++
      preEnd = end
    }
  }

  return res
}

// "fttfjofpnpfydwdwdnns"
// 2
console.log(maxPalindromes('fttfjofpnpfydwdwdnns', 2))
export {}
