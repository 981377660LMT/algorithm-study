console.log(typeof null) // "object" - 'null' has "object" type in js (backward compatibility)
console.log(null instanceof Object) // false - 'null' is primitive and doesn't have 'instanceof' keyword
console.log(typeof 1) // "number" - one of js types
console.log(1 instanceof Number) // false - '1' is primitive and doesn't have 'instanceof' keyword
console.log(1 instanceof Object) // false - same as above
console.log(Number(1) instanceof Object) // false - Number(1) === 1 - same as above
console.log(new Number(1) instanceof Object) // true - 'new Number(1)' is object, so it's correct
console.log(typeof true) // "boolean" - one of js types
console.log(true instanceof Boolean) // false - 'true' is primitive and doesn't have 'instanceof' keyword
console.log(true instanceof Object) // false - same as above
console.log(Boolean(true) instanceof Object) // false - Boolean(true) === true - same as above
console.log(new Boolean(true) instanceof Object) // true - 'new Boolean(true)' is object, so it's correct
console.log([] instanceof Array) // true - '[]' is instanceof Array and Object
console.log([] instanceof Object) // true - '[]' is instanceof Array and Object
console.log((() => {}) instanceof Object) // true - if it's not a primitive it's object. So callback is instanceof object
