// 求 1+2+...+n ，要求不能使用乘除法、for、while、if、else、switch、case等关键字及条件判断语句（A?B:C）。
// function sumNums(n: number): number {
//   const arr = Array.from({ length: n }, () => Array(n + 1).fill(0)).flat()
//   return arr.length >> 1
// }

function sumNums(n: number): number {
  n && (n += sumNums(n - 1))
  return n
}
