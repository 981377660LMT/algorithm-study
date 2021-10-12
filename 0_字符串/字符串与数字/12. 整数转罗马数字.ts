/**
 * @param {number} num
 * @return {string}
 * @description 类似于n进制转换 不过那个一般是递归做法 这个需要迭代
 * 除 模 除 模 ...
 */
const intToRoman = function (num: number): string {
  const roman = ['M', 'CM', 'D', 'CD', 'C', 'XC', 'L', 'XL', 'X', 'IX', 'V', 'IV', 'I']
  const arab = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]

  let res: string[] = []
  let index = 0
  while (num) {
    const [div, mod] = [~~(num / arab[index]), num % arab[index]]
    for (let i = 0; i < div; i++) {
      res.push(roman[index])
    }
    num = mod
    index++
  }
  return res.join('')
}

console.log(intToRoman(9))

export default 1
