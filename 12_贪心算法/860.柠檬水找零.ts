/**
 * @param {number[]} bills
 * @return {boolean}
 */
const lemonadeChange = function (bills: number[]): boolean {
  let fives = 0
  let tens = 0

  for (const bill of bills) {
    if (bill === 5) fives++
    else if (bill === 10) {
      tens++
      fives--
    } else {
      // 赶快把十元对出去
      if (tens > 0) {
        tens--
        fives--
      } else fives -= 3
    }
    if (fives < 0 || tens < 0) return false
  }
  return true
}

console.log(lemonadeChange([5, 5, 5, 10, 20]))

export default 1
