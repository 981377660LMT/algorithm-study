// https://leetcode.cn/problems/make-object-immutable/solution/typescript-bu-shi-yong-deep-proxy-de-shi-5fy7/
// 如果试图修改对象的键，则会产生以下错误消息： `Error Modifying: ${key}` 。
// 如果试图修改数组的索引，则会产生以下错误消息： `Error Modifying Index: ${index}` 。
// 如果试图调用会改变数组的方法，则会产生以下错误消息： `Error Calling Method: ${methodName}` 。
// 你可以假设只有以下方法能够改变数组： ['pop', 'push', 'shift', 'unshift', 'splice', 'sort', 'reverse'] 。

// !拦截对象的get方法来进行惰性代理，而不是初始化时遍历整个对象进行拦截

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

const propHandler: ProxyHandler<Obj> = {
  get(target, p) {
    const val = Reflect.get(target, p)
    if (!isObj(val)) return val
    return proxify(val)
  },
  set(target, prop) {
    if (Array.isArray(target)) {
      throw new Error(`Error Modifying Index: ${String(prop)}`)
    }
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
  }

  return new Proxy(obj, propHandler) as T
}

function isObj(obj: unknown): obj is Obj {
  return typeof obj === 'object' && obj !== null
}

export {}

if (require.main === module) {
  //   {"arr":[1,2,3]}
  // (obj) => { obj.arr.push(4); return 42; }
  const obj = makeImmutable({ arr: [1, 2, 3] })
  obj.arr.push(4)
}
