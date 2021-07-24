// 设计一个算法，找出只含素因子 2，3，5 的第 n 小的数
// 1 通常被视为丑数。
const nthUglyNumber = (num: number) => {
  const isUgly = (n: number) => {
    while (n % 2 === 0) n /= 2
    while (n % 3 === 0) n /= 3
    while (n % 5 === 0) n /= 5
    return n === 1
  }
  let count = 0
  let search = 0

  while (count < num) {
    search++
    if (isUgly(search)) {
      count++
    }
  }

  return search
}

console.log(nthUglyNumber(1690))

export {}
