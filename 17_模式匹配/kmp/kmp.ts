function kmp(s: string, p: string): number {
  const N = s.length
  const M = p.length
  const T = [0]
  for (let i = 1, len = 0; i < M; ) {
    if (p[i] === p[len]) T[i++] = ++len
    else if (len) len = T[len - 1]
    else T[i++] = 0
  }
  for (let i = 0, len = 0; i < N; ) {
    if (s[i] === p[len]) {
      len++
      i++
      if (len === M) return i - M
    } else if (len) len = T[len - 1]
    else i++
  }
  return -1
}

export { kmp }
