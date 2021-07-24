// 将一个整数分解为若干质因数之乘积
// 你需要从小到大排列质因子
// 需要while循环每个可能的(质)因子不断将质因子推入数组中
const primeFactorize = (n: number): number[] => {
  if (n === 1) return []
  let currentNum = n
  const res: number[] = []

  // 对于每个质因子
  for (let primeFactor = 2; primeFactor ** 2 <= n; primeFactor++) {
    while (currentNum % primeFactor === 0) {
      currentNum = currentNum / primeFactor
      res.push(primeFactor)
    }
  }

  // 最后位质因子的情况
  if (currentNum !== 1) res.push(currentNum)

  return res
}

console.log(primeFactorize(17))
console.log(primeFactorize(34))
console.log(primeFactorize(24))
