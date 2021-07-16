const foo = (arr: string[]): { count: number; value: string } => {
  if (arr.length === 0) return { count: 0, value: '' }

  const memo = new Map<string, number>()
  let max = { count: 0, value: '' }
  arr.forEach(str => {
    memo.has(str) ? memo.set(str, memo.get(str)! + 1) : memo.set(str, 1)
    if (memo.get(str)! > max.count) {
      max = { count: memo.get(str)!, value: str }
    }
  })

  return max
}

console.log(foo(['192.168.1.1', '192.118.2.1', '192.168.1.1']))
export {}
