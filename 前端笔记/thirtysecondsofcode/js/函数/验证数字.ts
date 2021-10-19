const isNumerical = (n: string) => {
  const num = parseFloat(n)
  return !Number.isNaN(num) && Number.isFinite(num) && Number(n) == n
}

isNumerical('10') // true
isNumerical('a') // false
