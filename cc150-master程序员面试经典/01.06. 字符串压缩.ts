function compressString(S: string): string {
  const trans = S.replace(/(\w)\1*/g, (match, group) => `${group}${match.length}`)
  console.log(trans)
  return trans.length < S.length ? trans : S
}

console.log(compressString('aabcccccaa'))
