function numOfSubarrays(arr: number[]): number {
  const MOD = 10 ** 9 + 7
  let res = 0
  let sum = 0
  let odd = 0
  let even = 1 // (S0)
  for (const num of arr) {
    sum += num
    if (sum & 1) {
      odd++
      res += even % MOD
    } else {
      even++
      res += odd % MOD
    }
  }

  return res % MOD
}

console.log(numOfSubarrays([1, 3, 5]))
