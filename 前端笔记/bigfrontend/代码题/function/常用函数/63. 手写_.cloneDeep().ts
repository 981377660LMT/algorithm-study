// 深拷贝 deepcopy
// 简单起见，该题目中你只需要支持：
// 基础数据类型(包括Symbol) 及其wrapper object。
// 简单Object（仅需处理可枚举属性）
// 数组Array

import _ from 'lodash'

function deepcopy<T = any>(o: T, visited = new WeakMap()): T {
  if (!isObject(o)) return o
  if (visited.has(o)) return visited.get(o) // !如果已经访问过，则直接返回

  const res = Array.isArray(o) ? [] : ({} as any) // !暂时只支持数组和对象 不支持Set和Map等
  visited.set(o, res)

  const keys = [...Object.getOwnPropertySymbols(o), ...Object.getOwnPropertyNames(o)]
  keys.forEach(key => {
    const val = (o as Record<PropertyKey, any>)[key]
    res[key] = deepcopy(val, visited)
  })
  return res
}

function isObject(o: unknown): o is object {
  return typeof o === 'object' && o !== null
}

export { deepcopy }

if (require.main === module) {
  const o = {
    a: 1,
    b: { c: 2, d: { e: 3 } },
    f: [4, 5, 6],
    g: Symbol('g'),
    h: new Set([7, 8, 9])
  }

  // 循环引用
  const bad = { cycle: o }
  Object.assign(o, { cycle: bad })

  const res = deepcopy(o)
  console.log(res)
  console.log(_.cloneDeep(o))
}

// Object.keys() 返回对象的可枚举属性
// Reflect.ownKeys() :
// 相当于
// Object.getOwnPropertyNames(target) concat(Object.getOwnPropertySymbols(target)

// Object.keys 获取自身可枚举属性
// Object.getOwnPropertyNames 获取自身所有属性，除了Symbol
