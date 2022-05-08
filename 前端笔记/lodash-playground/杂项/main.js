const _ = require('lodash')

// iter
const iterator = _([1, 2, 3])
console.log(iterator.next())

// 转化 value 为属性路径的数组 。
console.log(_.toPath('a[0].b.c'))
