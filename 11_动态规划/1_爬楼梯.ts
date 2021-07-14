const climbStairs = (n: number): number => {
  if (n <= 1) return 1
  return climbStairs(n - 1) + climbStairs(n - 2)
}

console.log(climbStairs(6))

export {}
