// @ts-ignore
console.log(typeof a) // undefined
// undefined 不是未定义，两者有区别。尝试去读一个未定义的变量的值其实会直接Reference Error
// typeof 不能区分未定义，还是定义了但是没有值。两者都会都会返回undefined
// typeof 一个未定义的变量不会触发Reference Error
console.log(Number.EPSILON.toString(2))
console.log('0000000000000000000000000000000000000000000000000001'.length) // 52
