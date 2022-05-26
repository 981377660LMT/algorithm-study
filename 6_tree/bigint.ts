// JS中缺少显式整数类型常常令人困惑。
// 许多编程语言支持多种数字类型，如浮点型、双精度型、整数型和双精度型，但JS却不是这样。在JS中，按照IEEE 754-2008标准的定义，所有数字都以双精度64位浮点格式表示。
const big = BigInt(0)
console.log(big + 1n)
console.log(BigInt('666'))
console.log(BigInt(true))

// bigInt不能使用Math的方法比较
// 注意BigInt与bigint的区别
