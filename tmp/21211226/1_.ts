function isSameAfterReversals(num: number): boolean {
  if (num === 0) return true
  return !num.toString().endsWith('0')
}
console.log(isSameAfterReversals(526))
console.log(isSameAfterReversals(1800))
