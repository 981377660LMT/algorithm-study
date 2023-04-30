const LOWERCASE = 'abcdefghijklmnopqrstuvwxyz'
const UPPERCASE = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
const DIGITS = '0123456789'
const MOD = 1e9 + 7
const EPS = 1e-8
const DIRS4 = [
  [-1, 0],
  [0, 1],
  [1, 0],
  [0, -1]
]
const DIRS8 = [
  [-1, 0],
  [-1, 1],
  [0, 1],
  [1, 1],
  [1, 0],
  [1, -1],
  [0, -1],
  [-1, -1]
]

function smallestBeautifulString(s: string, k: number): string {
  const n = s.length
  const res: number[] = []
  const ords = s.split('').map(c => c.charCodeAt(0))
  try {
    dfs(0, true, 0, -1)
  } catch (error) {}
  return res.map(c => String.fromCharCode(c)).join('')

  function dfs(pos: number, isLimit: boolean, pre1: number, pre2: number) {
    if (pos === n) {
      if (!isLimit) {
        throw new Error('ok')
      }
      return
    }
    const lower = isLimit ? ords[pos] : 97
    for (let cur = lower; cur < 97 + k; cur++) {
      if (cur === pre1 || cur === pre2) {
        continue
      }
      res.push(cur)
      dfs(pos + 1, isLimit && cur <= ords[pos], cur, pre1)
      res.pop()
    }
  }
}
// A = [ord(c) - ord('a') for c in s]
// n = len(A)
// if all(A[i] == k - 1 - i % 3 for i in range(n)):
//     return ''

const n = 5000
const k = 4
const ords = Array(n)
for (let i = 0; i < n; i++) {
  ords[i] = String.fromCharCode(97 + (k - 1 - (i % 3)))
}
const s = ords.join('')
console.time('a')
console.log(smallestBeautifulString(s, k))
console.timeEnd('a')
export {}
