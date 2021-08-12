const GCD = (...nums: number[]): number => {
  const _GCD = (a: number, b: number): number => (b === 0 ? a : GCD(b, a % b))
  return nums.reduce(_GCD)
}

if (require.main === module) {
  console.log(GCD(3, 6, 8))
}

export { GCD }
