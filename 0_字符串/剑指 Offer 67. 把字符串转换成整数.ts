function strToInt(s: string): number {
  const pattern = /^[-+]?\d+/g
  const match = s.trim().match(pattern)
  return match ? Math.max(Math.min(Number(match[0]), 2 ** 31 - 1), -(2 ** 31)) : 0
}

console.log(strToInt('   -42'))
console.log(strToInt('4193 with words'))
console.log(strToInt('words and 987'))
