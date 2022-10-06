import { useStringHasher } from '../StringHasher'

function sumScores(s: string): number {
  const n = s.length
  const hasher = useStringHasher(s)

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
      if (hasher(start, start + mid) === hasher(0, mid)) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }

    return right
  }
}

console.log(sumScores('babab'))

export {}
