// function foo(a, b, undefined, undefined) {
//   console.log(arguments.length)
// }
// console.log(foo.length)
// foo(1, 2, 3)
////////////////////////////////////////////////////////////////////
// if (
//   function foo() {
//     console.log('BFE')
//   }
// ) {
//   console.log('dev')
// }
// // ReferenceError: foo is not defined
// foo()
////////////////////////////////////////////////////////////////////
if ([]) {
  console.log([] == false) // true
}
console.log(undefined == false) // false
console.log(null == false) // false
console.log(undefined == 0) // false: undefined -> NaN
console.log(undefined < 0) // false: undefined -> NaN
console.log(undefined > 0) // false: undefined -> NaN
console.log(undefined <= 0) // false: undefined -> NaN
console.log(undefined >= 0) // false: undefined -> NaN
console.log(null == 0) // Special rule: null is not converted to 0 here
console.log(null < 0) // false: null -> 0
console.log(null > 0) // false: null -> 0
console.log(null <= 0) // true: null -> 0
console.log(null >= 0) // true: null -> 0
