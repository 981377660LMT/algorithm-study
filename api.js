// const a = [1, 23]
// a.splice(0)

// console.log(a)
console.log(-1 >> 1)
// 这个是八进制的数字，用2的radix会NaN
console.log(parseInt(0100, 2))

// 二进制字符串转十进制数字 100
console.log(parseInt('0100', 2))
// 去除字符串前面的0
console.log(parseInt('0100', 10).toString())

// 十进制数字转二进制字符串 '010'
console.log(Number(2).toString(2).padStart(3, 0))
console.log(Array.from([1, 2, 3], (v, k) => [v, k]))

// radix为0，则设置radix为默认值10
console.log(parseInt(1, 0))
parseInt('5/8/2017', 'javascript is such funny')

console.log([
  ...new Set([
    [1, 2, 3],
    [3, 2, 1],
  ]),
])

let a = 1
console.log(a++)

let state = 1
const newState = ++state % 4
console.log(newState)

console.assert(1 === 2, 'wrong')

const arr = []
arr.push(1, 2)
console.log(arr[4])

console.log(...arr.keys(), ...arr.entries())

console.log([1, 2, 3].slice(0, -1))

console.log(undefined + 1 || 1, undefined + 1)

console.log(Math.trunc('1.121323435'))
console.log(Math.trunc('-1.121323435'))
