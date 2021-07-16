// 123 反转之后是 321。 900 反转之后是 9。
// 009这种形式需要转为9
const convert = (n: number) => parseInt(n.toString().split('').reverse().join(''))

console.log(convert(123), convert(900))
export {}
