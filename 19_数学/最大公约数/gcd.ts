const GCD = (...nums: number[]): number => {
  // console.log(nums)
  const _GCD = (a: number, b: number): number => (b === 0 ? a : GCD(b, a % b))
  return nums.reduce(_GCD)
}

if (require.main === module) {
  console.log(GCD(3, 6, 8))
  console.log(GCD(1, 1))
}

export { GCD }
