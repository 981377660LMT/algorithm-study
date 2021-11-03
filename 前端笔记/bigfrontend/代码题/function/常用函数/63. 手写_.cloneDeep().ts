// 简单起见，该题目中你只需要支持：

// 基础数据类型(包括Symbol) 及其wrapper object。
// 简单Object（仅需处理可枚举属性）
// 数组Array

const memo = new WeakMap<object, any>()

function cloneDeep(obj: any) {
  if (obj === null || typeof obj !== 'object') return obj
  if (memo.has(obj)) return memo.get(obj)

  const res = Array.isArray(obj) ? [] : ({} as any)
  memo.set(obj, res)

  const keys = [...Object.getOwnPropertySymbols(obj), ...Object.keys(obj)]
  for (const key of keys) {
    const val = obj[key]
    res[key] = cloneDeep(val)
  }

  return res
}

console.log(cloneDeep({ a: { b: { c: 3 }, d: 4, e: Symbol() } }))
export {}

// Object.keys() 返回对象的可枚举属性
// Reflect.ownKeys() :
// 相当于
// Object.getOwnPropertyNames(target) concat(Object.getOwnPropertySymbols(target)
