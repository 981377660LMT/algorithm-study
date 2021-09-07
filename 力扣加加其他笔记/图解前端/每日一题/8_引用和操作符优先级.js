var a = { x: 1 }
var b = a
a = a.x = { x: 1 }
console.log(a)
console.log(b)

// 操作符的的运算优先级问题, .的优先级高于赋值语句
