// function customSortString(order: string, s: string): string {
//   const weight = Array<number>(26).fill(-1)
//   for (let i = 0; i < order.length; i++) {
//     weight[order.codePointAt(i)! - 97] = i
//   }

//   return s
//     .split('')
//     .sort((s1, s2) => weight[s1.codePointAt(0)! - 97] - weight[s2.codePointAt(0)! - 97])
//     .join('')
// }

// 如果order中x在y之前出现，那么返回的字符串中x也应出现在y之前。
function customSortString(order: string, s: string): string {
  return s
    .split('')
    .sort((s1, s2) => order.indexOf(s1) - order.indexOf(s2))
    .join('')
}

console.log(customSortString('cba', 'abcd'))
