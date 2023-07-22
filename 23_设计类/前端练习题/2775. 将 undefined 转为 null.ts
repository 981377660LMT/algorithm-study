// 2775. 将 undefined 转为 null
// https://leetcode.cn/problems/undefined-to-null/
//
// 编写一个名为 undefinedToNull 的函数，该函数接受一个深层嵌套的对象或数组 obj ，
// 并创建该对象的副本，将其中的任何 undefined 值替换为 null 。
//
// 当使用 JSON.stringify() 将对象转换为 JSON 字符串时，undefined 值与 null 值的处理方式不同。
// 该函数有助于确保序列化数据不会出现意外错误。

function undefinedToNull(obj: Record<any, any>): unknown {
  if (obj === undefined) return null
  if (obj === null || typeof obj !== 'object') return obj
  if (Array.isArray(obj)) return obj.map(item => undefinedToNull(item))
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) obj[key] = undefinedToNull(obj[key])
  }
  return obj
}

/**
 * undefinedToNull({"a": undefined, "b": 3}) // {"a": null, "b": 3}
 * undefinedToNull([undefined, undefined]) // [null, null]
 */

export {}
