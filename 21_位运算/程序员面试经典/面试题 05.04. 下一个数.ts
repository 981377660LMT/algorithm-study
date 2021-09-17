// 给定一个正整数，找出与其二进制表达式中1的个数相同且大小最接近的那两个数
// num的范围在[1, 2147483647]之间；
// 如果找不到前一个或者后一个满足条件的正数，那么输出 -1。
function findClosedNumbers(num: number): number[] {
  if (num === 1) return [2, -1]
  if (num === 0x7fffffff) return [-1, -1]

  const getNextMax = (x: number) => {
    const getLowbit = (x: number) => x & -x
    const getToZero = (x: number) => x + getLowbit(x)
    const getFoo = (x: number) => x & ~getToZero(x)
    const lowbit = getLowbit(x)
    const toZero = getToZero(x)
    return ((getFoo(x) / lowbit) >> 1) | toZero
  }

  let max = getNextMax(num)
  let min = ~getNextMax(~num)
  if (max < 0) max = -1
  return [max, min]
}

console.log(findClosedNumbers(0b10))
// 输出：[4, 1] 或者（[0b100, 0b1]）
console.log(Number(0x7fffffff).toString(2).length)
