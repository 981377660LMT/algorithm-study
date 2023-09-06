// 实现split

function mySplit(s: string, separator: string): string[] {
  const res: string[] = []
  let last = 0
  while (last < s.length) {
    const next = s.indexOf(separator, last)
    if (next === -1) {
      res.push(s.slice(last))
      break
    }
    res.push(s.slice(last, next))
    last = next + separator.length
  }
  return res
}

export {}

if (require.main === module) {
  const arr = mySplit('a,b,c,d', ',')
  console.log(arr)
}
