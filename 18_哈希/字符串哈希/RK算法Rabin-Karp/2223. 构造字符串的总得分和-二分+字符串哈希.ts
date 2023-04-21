import { SafeStringHasher } from '../StringHasher-new'

function sumScores(s: string): number {
  const n = s.length
  const H = new SafeStringHasher(s)

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
      const h1 = H.query(start, start + mid)
      const h2 = H.query(0, mid)
      if (h1 === h2) {
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
  const h = new SafeStringHasher('babab')
  console.log(h.query(0, 1) === h.query(2, 3))
}

export {}
