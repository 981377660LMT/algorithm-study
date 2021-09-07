function judgeSquareSum(c: number): boolean {
  const isSquare = (num: number) => (~~Math.sqrt(num)) ** 2 === num
  const sqrt = ~~Math.sqrt(c)
  for (let i = 0; i <= sqrt; i++) {
    if (isSquare(c - i ** 2)) return true
  }
  return false
}

console.log(judgeSquareSum(5))

export default 2
