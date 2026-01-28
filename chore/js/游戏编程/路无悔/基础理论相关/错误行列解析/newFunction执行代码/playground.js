// const a = [1, 23]
// a.splice(0)

// console.log(a)
console.log(-1 >> 1)
// 这个是八进制的数字，用2的radix会NaN
console.log(parseInt(0100, 2))

// 二进制字符串转十进制数字 100
console.log(0b0100)
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
    [3, 2, 1]
  ])
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

// 1. 运行时错误 (Runtime Error)
console.log('--- 1. 运行时错误 ---')
try {
  // @ts-ignore
  console.log(undefinedVariable) // 引用不存在的变量
} catch (e) {
  printErrorDetails(e)
}

// 2. 语法错误 (Syntax Error) - 通过 new Function 触发
console.log('\n--- 2. 语法错误 (Syntax Error) ---')
try {
  // 少了分号前的赋值，或者不合法的语法
  new Function('var a = ;')
} catch (e) {
  printErrorDetails(e)
}

// 3. 正则语法错误 (Syntax Error - Regex)
console.log('\n--- 3. 正则语法错误 ---')
try {
  new Function('var r = /abc') // 缺少闭合斜杠
} catch (e) {
  printErrorDetails(e)
}

function printErrorDetails(e) {
  console.log('类型 (name):', e.name)
  console.log('消息 (message):', e.message)

  // 某些浏览器特有属性 (非标准，Node环境通常为undefined)
  if (e.lineNumber) console.log('行号 (lineNumber - Firefox/Old):', e.lineNumber)
  if (e.columnNumber) console.log('列号 (columnNumber - Firefox/Old):', e.columnNumber)

  console.log('堆栈 (stack):')
  console.log(e.stack)
}
