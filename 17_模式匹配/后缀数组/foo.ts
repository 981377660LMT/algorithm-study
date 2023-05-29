function compareSubstrNaive(s: string, a: number, b: number, c: number, d: number): -1 | 0 | 1 {
  while (a < b && c < d && s[a] === s[c]) {
    a++
    c++
  }
  if (a === b) return c === d ? 0 : -1
  if (c === d) return 1
  return s[a] < s[c] ? -1 : 1
}

function compareSubstrWithSlice(s: string, a: number, b: number, c: number, d: number): -1 | 0 | 1 {
  const s1 = s.slice(a, b)
  const s2 = s.slice(c, d)
  return s1.localeCompare(s2) as -1 | 0 | 1
}

const s = 'abcdeadsnfjbdws'.repeat(10000)
console.time('compareSubstrNaive')
for (let i = 0; i < 1e6; i++) {
  compareSubstrNaive(s, 0, 10000, 10000, 20000)
}

console.timeEnd('compareSubstrNaive')

console.time('compareSubstrWithSlice')
for (let i = 0; i < 1e6; i++) {
  compareSubstrWithSlice(s, 0, 10000, 10000, 20000)
}
console.timeEnd('compareSubstrWithSlice')
