import { StringHasher2 } from '../StringHasher-new'

function sumScores(s: string): number {
  const n = s.length
  const H = new StringHasher2()
  const table = H.build(s)

  let res = 0
  for (let i = 1; i <= n; i++) {
    // if (s[n - i] !== s[0]) continue
    const count = countPre(i, n - i)
    res += count
  }

  return res

  function countPre(curLen: number, start: number): number {
    let left = 1
    let right = curLen
    while (left <= right) {
      const mid = (left + right) >> 1
      const h1 = H.query(table, start, start + mid)
      const h2 = H.query(table, 0, mid)
      if (h1[0] === h2[0] && h1[1] === h2[1]) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return right
  }
}

if (require.main === module) {
  console.log(sumScores('babab'))
}

export {}
