// 在检查真值的数据时，Boolean ()函数特别有用，可能比双重否定(! !)操作更具可读性:

const x = new Boolean(false)
const y = Boolean(false)

if (x) {
  // This code is executed
  console.log(1, x.valueOf())
}

if (y) {
  // This code is not executed
  console.log(1, y.valueOf())
}

console.log(typeof x, typeof y)
