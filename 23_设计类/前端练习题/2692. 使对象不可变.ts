// https://leetcode.cn/problems/make-object-immutable/solution/typescript-bu-shi-yong-deep-proxy-de-shi-5fy7/
// 如果试图修改对象的键，则会产生以下错误消息： `Error Modifying: ${key}` 。
// 如果试图修改数组的索引，则会产生以下错误消息： `Error Modifying Index: ${index}` 。
// 如果试图调用会改变数组的方法，则会产生以下错误消息： `Error Calling Method: ${methodName}` 。
// 你可以假设只有以下方法能够改变数组： ['pop', 'push', 'shift', 'unshift', 'splice', 'sort', 'reverse'] 。

// !deep watch

type Obj = Record<PropertyKey, unknown> | unknown[]
const MUTABLE_ARRAY_METHODS: (keyof [])[] = [
  'pop',
  'push',
  'shift',
  'unshift',
  'splice',
  'sort',
  'reverse'
]

function makeImmutable<T extends Obj>(obj: T): T {
  return proxify(obj)
}

/**
 * const obj = makeImmutable({x: 5});
 * obj.x = 6; // throws "Error Modifying x"
 */

const arrayHandler: ProxyHandler<unknown[]> = {
  set(_, prop) {
    throw new Error(`Error Modifying Index: ${String(prop)}`)
  }
}

const objectHandler: ProxyHandler<Record<PropertyKey, any>> = {
  set(_, prop) {
    throw new Error(`Error Modifying: ${String(prop)}`)
  }
}

const methodHandler: ProxyHandler<(...args: unknown[]) => unknown> = {
  apply(target) {
    throw new Error(`Error Calling Method: ${target.name}`)
  }
}

function proxify<T extends Obj>(obj: T): T {
  if (Array.isArray(obj)) {
    MUTABLE_ARRAY_METHODS.forEach(method => {
      obj[method] = new Proxy(obj[method], methodHandler)
    })
    obj.forEach((val, index) => {
      if (isObj(val)) {
        obj[index] = proxify(val)
      }
    })
    return new Proxy(obj, arrayHandler) as T
  }

  const keys = Object.keys(obj)
  keys.forEach(key => {
    const val = obj[key]
    if (isObj(val)) {
      obj[key] = proxify(val)
    }
  })
  return new Proxy(obj, objectHandler)
}

function isObj(obj: unknown): obj is Obj {
  return typeof obj === 'object' && obj !== null
}

export {}

if (require.main === module) {
  const obj = makeImmutable({ x: 5 })
  obj.x = 5 // throws "Error Modifying x"
}
