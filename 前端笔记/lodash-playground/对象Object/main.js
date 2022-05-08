const _ = require('lodash')

// _.at(object, [paths])
// 创建一个数组，值来自 object 的paths路径相应的值。
const object = { a: [{ b: { c: 3 } }, 4] }

console.log(_.at(object, ['a[0].b.c', 'a[1]']))
