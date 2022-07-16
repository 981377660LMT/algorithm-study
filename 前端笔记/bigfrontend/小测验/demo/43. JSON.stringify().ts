// attention that for console.log('"a"'), you should enter ""a""
// please refer to format guide

console.log(JSON.stringify(['false', false]))
console.log(JSON.stringify([NaN, null, Infinity, undefined]))
console.log(JSON.stringify({ a: null, b: NaN, c: undefined }))
// ["false",false]
// [null,null,null,null]
// {"a":null,"b":null}
