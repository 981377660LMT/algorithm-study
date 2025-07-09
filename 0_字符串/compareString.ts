function compareString(a: string, b: string): 0 | -1 | 1 {
  if (a === b) return 0
  return a < b ? -1 : 1
}

export { compareString }

if (typeof require !== 'undefined' && typeof module !== 'undefined' && require.main === module) {
  const s1 = 'M'
  const s2 = 'b'
  console.log(s1.localeCompare(s2)) // 1
  console.log(s1.codePointAt(0), s2.codePointAt(0)) // 77 98
  console.log(compareString(s1, s2)) // -1
}
