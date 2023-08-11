/* eslint-disable @typescript-eslint/ban-ts-comment */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable max-len */

// 遍历 类数组，普通对象，Map，Set，Iterable

// TODO:
// 1.用函数重载代替extends+联合类型判断
//   import { forEach } from 'lodash'
// !2. 从 Iterable 中除去 string 类型
// !3. extends/重载 的顺序要越具体的放越前面，越抽象的放越后面

type CallbackFnLike = (value: any, key: any, obj: any) => void

function forEach<T extends ArrayLike<any> | Record<any, any> | Map<any, any> | Set<any> | Iterable<any>>(
  obj: Exclude<T, string>,
  callbackFn: T extends Set<infer U>
    ? (value: U, key: U, set: T) => void
    : T extends Map<infer K, infer V>
    ? (value: V, key: K, map: T) => void
    : T extends ArrayLike<infer U>
    ? (value: U, key: number, array: T) => void
    : T extends Iterable<infer U>
    ? (value: U, key: U, iterable: T) => void
    : T extends Record<infer K, infer V>
    ? (value: V, key: K, record: T) => void
    : 'error: obj must be one of ArrayLike, Record, Map, Set, Iterable'
): void {
  if (obj === null || typeof obj === 'undefined') return

  if (typeof obj !== 'object') obj = [obj] as any

  if (isArray(obj)) {
    for (let i = 0; i < obj.length; i++) {
      ;(callbackFn as CallbackFnLike).call(null, obj[i], i, obj)
    }
  } else if (isMap(obj) || isSet(obj)) {
    obj.forEach((value, key) => {
      ;(callbackFn as CallbackFnLike).call(null, value, key, obj)
    })
  } else if (isIterable(obj)) {
    for (const value of obj) {
      ;(callbackFn as CallbackFnLike).call(null, value, value, obj)
    }
  } else {
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        ;(callbackFn as CallbackFnLike).call(null, obj[key], key, obj)
      }
    }
  }
}

function isArray(o: unknown): o is Array<any> {
  return Array.isArray(o)
}

function isMap(o: unknown): o is Map<any, any> {
  return o instanceof Map
}

function isSet(o: unknown): o is Set<any> {
  return o instanceof Set
}

/**
 * @see {@link https://stackoverflow.com/a/32538867}
 */
function isIterable(o: unknown): o is Iterable<any> {
  // @ts-ignore
  return o != null && typeof o[Symbol.iterator] === 'function'
}

export { forEach }

if (require.main === module) {
  forEach([1, 2, 3] as const, (value, key, obj) => {
    console.log(value, key, obj)
  })

  forEach({ a: 1, b: 2, c: 3 }, (value, key, obj) => {
    console.log(value, key, obj)
  })

  forEach(
    new Map([
      ['a', 1],
      ['b', 2],
      ['c', 3]
    ]),
    (value, key, obj) => {
      console.log(value, key, obj)
    }
  )

  forEach(new Set([1, 2, 3]), (value, key, obj) => {
    console.log(value, key, obj)
  })

  forEach(new URL('http://www.baidu.com'), (value, key, obj) => {
    console.log(value, key, obj)
  })
}
