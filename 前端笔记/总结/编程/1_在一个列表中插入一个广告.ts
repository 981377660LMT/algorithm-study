export {}
// 请把俩个数组 [A1, A2, B1, B2, C1, C2, D1, D2] 和 [A, B, C, D]，合并为 [A1, A2, A, B1, B2, B, C1, C2, C, D1, D2, D]
function merge(arr1: string[], arr2: string[]): string[] {
  // const pad = String.fromCodePoint(0xffff)
  // console.log(pad)
  // return [...arr1, ...arr2.map(v => v + pad)].sort().map(v => v.replace(pad, ''))
  return [...arr1, ...arr2].sort(
    (a, b) => a.codePointAt(0)! - b.codePointAt(0)! || b.length - a.length
    // a.codePointAt(1)! - b.codePointAt(1)!
  )
}

console.log(merge(['A1', 'A2', 'B1', 'B2', 'C1', 'C2', 'D1', 'D2'], ['A', 'B', 'C', 'D']))
console.log('$'.codePointAt(0), String(1).codePointAt(0))
// tring.fromCharCode 用于从unicode码点返回对应字符,但码点大于0xFFFF的不行
// 所以String.fromCodePoint可以识别码点大于0xFFFF的字符
console.log(''.codePointAt(0))
