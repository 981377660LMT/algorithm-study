const foo = [0]
if (foo) {
  console.log(foo == true) // false - when using '==' both sides convert to numbers: 0 == 1 -> false
} else {
  console.log(foo == false)
}

// true
// true
// true
// true
console.log(Boolean([0]))
console.log(Boolean([]))
console.log([0] == false)
console.log([] == false)
