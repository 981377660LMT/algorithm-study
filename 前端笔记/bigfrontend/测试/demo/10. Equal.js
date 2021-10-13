console.log(0 == false)
console.log('' == false)
console.log([] == false) // true [] doesn't have number primitive, so toString() is called
console.log(undefined == false)
console.log(null == false)
console.log('1' == true)
console.log(1n == true)
console.log(' 1     ' == true)
