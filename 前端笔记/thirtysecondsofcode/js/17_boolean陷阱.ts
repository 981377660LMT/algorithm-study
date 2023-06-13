// 在检查真值的数据时，Boolean ()函数特别有用，可能比双重否定(! !)操作更具可读性:

const right = new Boolean(false)
const y = Boolean(false)

if (right) {
  // This code is executed
  console.log(1, right.valueOf())
}

if (y) {
  // This code is not executed
  console.log(1, y.valueOf())
}

console.log(typeof right, typeof y)
export {}
