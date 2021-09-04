interface NestDict<V> {
  [key: string]: NestDict<V> | V
}

function isSameDict(
  item1: NestDict<number | string> | string | number,
  item2: NestDict<number | string> | string | number
): boolean {
  console.log(item1, item2)
  if (typeof item1 !== typeof item2) return false
  if (typeof item1 === 'string' || typeof item1 === 'number') return item1 === item2
  const keys1 = Object.keys(item1).sort()
  const keys2 = Object.keys(item2).sort()
  const len = Math.max(keys1.length, keys2.length)
  for (let i = 0; i < len; i++) {
    if (keys1[i] !== keys2[i]) return false
    if (!isSameDict(item1[keys1[i]], (item2 as NestDict<string | number>)[keys2[i]])) return false
  }
  return true
}

console.log(isSameDict({ a: { b: 1, c: '2' } }, { a: { c: '2', b: 1 } }))

export { isSameDict }
