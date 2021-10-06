const sum = (a: number, b: number): number => {
  return a === 0 ? b : sum((a & b) << 1, a ^ b)
}

console.log(sum(1, 2))
