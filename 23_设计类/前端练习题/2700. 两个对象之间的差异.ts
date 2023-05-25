// https://leetcode.cn/problems/differences-between-two-objects/

// 请你编写一个函数，它接收两个深度嵌套的对象或数组 obj1 和 obj2 ，并返回一个新对象表示它们之间差异。
// 该函数应该比较这两个对象的属性，并识别任何变化。
// 返回的对象应仅包含从 obj1 到 obj2 的值不同的键。
// 对于每个变化的键，值应表示为一个数组 [obj1 value, obj2 value] 。
// 不存在于一个对象中但存在于另一个对象中的键不应包含在返回的对象中。
// 在比较两个数组时，数组的索引被视为它们的键。
// 最终结果应是一个深度嵌套的对象，其中每个叶子的值都是一个差异数组。
// 你可以假设这两个对象都是 JSON.parse 的输出结果。

// !1. obj1 和 obj2 的类型不同
// !2. obj1 和 obj2 的类型相同 => 基本类型直接返回/非基本类型递归处理

function objDiff(obj1: any, obj2: any): any {
  if (type(obj1) !== type(obj2)) return [obj1, obj2]
  if (!isObject(obj1)) return obj1 === obj2 ? {} : [obj1, obj2]
  const diff: Record<string, unknown> = {}
  const sameKeys = Object.keys(obj1).filter(key => key in obj2)
  sameKeys.forEach(key => {
    const subDiff = objDiff(obj1[key], obj2[key])
    if (Object.keys(subDiff).length) diff[key] = subDiff
  })
  return diff
}

function type(obj: unknown): string {
  return Object.prototype.toString.call(obj).slice(8, -1)
}

function isObject(obj: unknown): obj is Record<string, unknown> {
  return typeof obj === 'object' && obj !== null
}

export {}
