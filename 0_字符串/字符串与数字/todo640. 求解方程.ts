/**
 * @param {string} equation
 * @return {string}
 * 怎么做到有负号就优先匹配负号
 */
const solveEquation = function (equation: string): string {
  const leftX = [...equation.matchAll(/(?<!=.*?)([-+]?\d*x)/g)]
    .map<number>(n => {
      if (n[1] === '+x' || n[1] === 'x') return 1
      else if (n[1] === '-x') return -1
      else return parseInt(n[1])
    })
    .reduce((pre, cur) => pre + cur, 0)

  const rightX = [...equation.matchAll(/(?<==.*?)([-+]?\d*x)/g)]
    .map<number>(n => {
      if (n[1] === '+x' || n[1] === 'x') return 1
      else if (n[1] === '-x') return -1
      else return parseInt(n[1])
    })
    .reduce((pre, cur) => pre + cur, 0)

  const leftNum = [...equation.matchAll(/(?<!=.*?)((?:\-\+)?\d+)[^x]/g)]
    .map(n => {
      console.log(n)
      return parseInt(n[1])
    })
    .reduce((pre, cur) => pre + cur, 0)

  const rightNum = [...equation.matchAll(/(?<==.*?)((?:\-\+)?\d+)[^x]/g)]
    .map(n => parseInt(n[1]))
    .reduce((pre, cur) => pre + cur, 0)

  console.log(leftX)
  console.log(rightX)
  console.log(leftNum)
  console.log(rightNum)
  if (leftX === rightX && leftNum === rightNum) return 'Infinite solutions'
  else if (leftX === rightX && leftNum !== rightNum) return 'No solution'
  else return `x=${(rightNum - leftNum) / (leftX - rightX)}`

  return ''
}

// console.log(solveEquation('x+5-3+x=6+x-2'))
console.log(solveEquation('x+5-3+x=6+x-2'))
console.log([].reduce(x => x, 0))
// 输出: "x=2"
// 如果方程没有解，请返回“No solution”。
// 如果方程有无限解，则返回“Infinite solutions”。
export default 1

console.log(parseInt('-x'))
console.log(parseInt('x'))
console.log(parseInt('2x'))
console.log(parseInt('-2x'))
