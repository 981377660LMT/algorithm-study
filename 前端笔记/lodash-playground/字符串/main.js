const _ = require('lodash')

console.log(_.capitalize('asDDaa')) // 转换字符串string首字母为大写，剩下为小写。
console.log(_.camelCase('__FOO_BAR__'))
console.log(_.snakeCase('pythonStyle')) // 转换字符串string为 snake case.
console.log(_.kebabCase('fileName')) // 转换字符串string为 kebab case.

console.log(_.escape('fred, barney, & pebbles')) // 转义string中的 "&", "<", ">", '"', "'", 和 "`" 字符为HTML实体字符。
