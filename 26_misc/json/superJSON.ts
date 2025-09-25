import superjson from 'superjson'

// 1. 创建一个包含多种数据类型的复杂对象
const complexObject = {
  now: new Date(),
  bigNumber: 12345678901234567890n, // BigInt
  mySet: new Set([1, 2, 'hello', 1]), // Set
  myMap: new Map([
    ['key1', 'value1'],
    ['key2', { a: 1 }]
  ]), // Map
  undef: undefined, // undefined
  regex: /test/gi // RegExp
}

console.log('Original Object:', complexObject)
console.log('----------------------------------------')

// 2. 使用标准 JSON 处理（会导致信息丢失）
// 注意：JSON.stringify 默认不支持 BigInt，会抛出错误。
// 为了比较，我们先将其转为字符串，但其他类型信息依然会丢失。
const standardJsonString = JSON.stringify({
  ...complexObject,
  bigNumber: complexObject.bigNumber.toString() // 手动转换 BigInt
})
const standardParsedObject = JSON.parse(standardJsonString)

console.log('After standard JSON.parse:', standardParsedObject)
/*
  你会看到：
  - `now` (Date) 变成了字符串。
  - `bigNumber` (BigInt) 变成了字符串。
  - `mySet` (Set) 变成了空对象 {}。
  - `myMap` (Map) 变成了空对象 {}。
  - `undef` (undefined) 属性直接丢失了。
  - `regex` (RegExp) 变成了空对象 {}。
*/
console.log('----------------------------------------')

// 3. 使用 superjson 处理（信息被完整保留）
const superJsonString = superjson.stringify(complexObject)

console.log('SuperJSON string:', superJsonString)
// 这个字符串包含了 superjson 用于恢复类型的元数据

const superJsonParsedObject = superjson.parse(superJsonString)

console.log('After superjson.parse:', superJsonParsedObject)
/*
  你会看到：
  - 所有的数据类型，包括 Date, BigInt, Set, Map, undefined, RegExp 都被正确地恢复了。
  - 对象的结构和类型与原始对象完全一致。
*/

// 验证类型是否正确恢复
console.log('Is date a Date object?', superJsonParsedObject.now instanceof Date) // true
console.log('Is bigNumber a BigInt?', typeof superJsonParsedObject.bigNumber === 'bigint') // true
console.log('Is mySet a Set object?', superJsonParsedObject.mySet instanceof Set) // true
