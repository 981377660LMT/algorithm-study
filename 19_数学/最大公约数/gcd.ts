function gcd(...nums: number[]): number {
  // console.log(nums)
  const _gcd = (a: number, b: number): number => (b === 0 ? a : gcd(b, a % b))
  return nums.reduce(_gcd)
}

if (require.main === module) {
  console.log(gcd(3, 6, 8))
  console.log(gcd(1, 1))
}

export { gcd }
