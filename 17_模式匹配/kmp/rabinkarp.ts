function rabinkarp(s: string, p: string): number {
  const N = s.length
  const M = p.length
  const q = 1e9 + 7
  const D = maxCharCode(s) + 1
  let h = 1
  for (let i = 0; i < M - 1; i++) h = (h * D) % q
  let hash = 0
  let target = 0
  for (let i = 0; i < M; i++) {
    hash = (hash * D + code(s, i)) % q
    target = (target * D + code(p, i)) % q
  }
  for (let i = M; i <= N; i++) {
    if (check(i - M)) return i - M
    if (i === N) continue
    hash = ((hash - h * code(s, i - M)) * D + code(s, i)) % q
    if (hash < 0) hash += q
  }
  return -1

  function check(begin: number): boolean {
    if (hash !== target) return false
    for (let i = 0; i < M; i++) if (s[begin + i] !== p[i]) return false
    return true
  }
}

function maxCharCode(s: string): number {
  let D = 0
  for (let i = 0; i < s.length; i++) {
    D = Math.max(D, s.charCodeAt(i))
  }
  return D
}

function code(s: string, i: number): number {
  return s.charCodeAt(i)
}

export { rabinkarp }
