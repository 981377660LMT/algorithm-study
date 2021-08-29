/**
 * @param {number[]} digits
 * @return {number[]}
 * @description 给定一个由 整数 组成的 非空 数组所表示的非负整数，在该数的基础上加一。
 */
var plusOne = function (digits: number[]): number[] {
  const n = digits.length
  const res: number[] = []
  let carry = 1
  for (let index = n - 1; index >= 0; index--) {
    const sum = digits[index] + carry
    carry = ~~(sum / 10)
    res.push(sum % 10)
  }
  carry && res.push(1)
  return res.reverse()
}

console.log(plusOne([1, 2, 3]))
