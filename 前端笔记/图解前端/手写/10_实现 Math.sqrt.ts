// 不用数学库，求 sqrt(2)精确到小数点后 10 位
const sqrt = (num: number, precision: number) => {
  if (num < 0) return NaN
  let l = 0
  let r = num
  let mid = (l + r) >> 1
  while (+Math.abs(num - mid ** 2).toFixed(precision) > Math.pow(10, -precision)) {
    if (mid ** 2 > num) r = mid
    else if (mid ** 2 < num) l = mid
    else return mid
    mid = (r + l) / 2
  }
  return mid
}

console.log(sqrt(2, 10))
console.log(sqrt(4, 10))
