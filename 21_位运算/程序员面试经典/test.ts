const getNextMax = (x: number) => {
  const getLowbit = (x: number) => x & -x
  const getToZero = (x: number) => x + getLowbit(x)
  const getFoo = (x: number) => x & ~getToZero(x)
  console.log(getFoo(6))
  const lowbit = getLowbit(x)
  const toZero = getToZero(x)
  return ((getFoo(x) / lowbit) >> 1) | toZero
}

const testNum = 0b101101110101010110
// console.log(lowbit(testNum).toString(2))
// console.log(toZero(testNum).toString(2))
// console.log(kkk(testNum).toString(2))
console.log(getNextMax(0b101))

// 取反是怎么运算的
console.log(~2)
console.log(Number(-3).toString(2))
console.log(-3 & 2)
