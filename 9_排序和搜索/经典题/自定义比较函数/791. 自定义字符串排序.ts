// 如果order中x在y之前出现，那么返回的字符串中x也应出现在y之前。
function customSortString(order: string, s: string): string {
  const rank = new Uint32Array(26)
  for (let i = 0; i < order.length; i++) {
    rank[order.codePointAt(i)! - 97] = i
  }

  return s
    .split('')
    .sort((s1, s2) => rank[s1.codePointAt(0)! - 97] - rank[s2.codePointAt(0)! - 97])
    .join('')
}

console.log(customSortString('cba', 'abcd'))
export {}
