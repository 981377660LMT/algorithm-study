const a = 1
console.log(a) // 1 - variable is defined and equal 1

var b
console.log(b) // undefined - variable is not defined yet
b = 2

console.log(c) // undefined - VAR is hoisting
var c = 3

console.log(d) // Error - LET in not hoisting -> ReferenceError
let d = 2
