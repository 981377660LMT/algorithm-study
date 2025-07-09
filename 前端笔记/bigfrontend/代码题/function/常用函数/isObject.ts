import { isObject, isObjectLike, isArrayLikeObject, isPlainObject } from 'lodash'

// 综合对比示例
const testCases = [
  {}, // 普通对象
  [], // 数组
  'string', // 字符串
  function () {}, // 函数
  new Date(), // 日期对象
  /regex/, // 正则对象
  null, // null
  undefined, // undefined
  42, // 数字
  { length: 0 }, // 类数组对象
  // document.querySelectorAll('div'), // NodeList
  (function () {
    return arguments
  })(), // Arguments
  // Object.create(null) // 无原型对象
  new (class MyClass {})() // 自定义类实例
]

console.table(
  testCases.map(value => ({
    value: String(value),
    isObject: isObject(value),
    isObjectLike: isObjectLike(value),
    isArrayLikeObject: isArrayLikeObject(value),
    isPlainObject: isPlainObject(value)
  }))
)
