// 显然，JSON.parse() JSON.stringify()不支持JSON data type以外的数据类型。

// console.log(JSON.stringify({ a: 1n })) // TypeError: Do not know how to serialize a BigInt

// undefined会被忽略，或者转换为null。
console.log(JSON.stringify([undefined])) // "[null]"
console.log(JSON.stringify({ a: undefined })) // "{}"

// NaN和 Infinity也会被当作null处理。
console.log(JSON.stringify([NaN, Infinity])) // "[null,null]"
console.log(JSON.stringify({ a: NaN, b: Infinity })) // {"a":null,"b":null}
