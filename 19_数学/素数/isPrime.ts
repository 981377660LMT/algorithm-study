/**
 * @param {number} num - positive integer
 */
function isPrime(num: number) {
  if (num === 1) return false
  const sqrt = Math.sqrt(num)
  for (let i = 2; i <= sqrt; i++) {
    if (num % i === 0) return false
  }
  return true
}
