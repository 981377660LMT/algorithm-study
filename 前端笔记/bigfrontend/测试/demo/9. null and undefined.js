console.log(JSON.stringify([1, 2, null, 3]))
console.log(JSON.stringify([1, 2, undefined, 3])) // "[1,2,null,3]"
console.log(null === undefined) // null -> 0 and undefined -> NaN, then NOT strictly equal
console.log(null == undefined) // Special rule: true -> Just Remember it
console.log(null == 0) // false Special rule: null is not converted to 0 here
console.log(null < 0)
console.log(null > 0)
console.log(null <= 0) // true
console.log(null >= 0) // true
console.log(undefined == 0) // false: undefined -> NaN
console.log(undefined < 0)
console.log(undefined > 0)
console.log(undefined <= 0)
console.log(undefined >= 0)
