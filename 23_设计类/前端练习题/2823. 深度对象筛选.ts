/* eslint-disable max-len */
// 给定一个对象 obj 和一个函数 fn，返回一个经过筛选的对象 filteredObject。

import { isPrimitive } from '../../ts-utils/isPrimitive'

// 函数 deepFilter 应该在对象 obj 上执行深度筛选操作。深度筛选操作应该移除筛选函数 fn 输出为 false 的属性，
// 以及在键被移除后仍然存在的任何空对象或数组。

function deepFilter(obj: Record<string, any>, fn: (o: unknown) => boolean): Record<string, any> | undefined {
  if (isPrimitive(obj)) return fn(obj) ? obj : undefined

  if (Array.isArray(obj)) {
    const container: any[] = []
    obj.forEach(value => {
      const res = deepFilter(value, fn)
      if (res !== undefined) container.push(res)
    })
    return container.length ? container : undefined
  }

  const container: Record<string, any> = {}
  for (const key in obj) {
    if (!Object.prototype.hasOwnProperty.call(obj, key)) continue
    const remaining = deepFilter(obj[key], fn)
    if (remaining !== undefined) container[key] = remaining
  }
  return Object.keys(container).length ? container : undefined
}

if (require.main === module) {
  // {"a":1,"b":"2","c":3,"d":"4","e":5,"f":6,"g":{"a":1}}
  // (x) => typeof x === "string"
  console.log(deepFilter({ a: 1, b: '2', c: 3, d: '4', e: 5, f: 6, g: { a: 1 } }, (x: any) => typeof x === 'string'))
}
