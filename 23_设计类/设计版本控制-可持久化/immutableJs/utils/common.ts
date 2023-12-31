/* eslint-disable max-len */
/* eslint-disable prefer-destructuring */

export function isArrayLike(value: any): boolean {
  if (Array.isArray(value) || typeof value === 'string') {
    return true
  }

  return (
    value &&
    typeof value === 'object' &&
    Number.isInteger(value.length) &&
    value.length >= 0 &&
    (value.length === 0
      ? // Only {length: 0} is considered Array-like.
        Object.keys(value).length === 1
      : // An object is only Array-like if it has a property where the last value
        // in the array-like may be found (which could be undefined).
        // eslint-disable-next-line no-prototype-builtins
        value.hasOwnProperty(value.length - 1))
  )
}

const toString = Object.prototype.toString

export default function isPlainObject(value: any): boolean {
  // The base prototype's toString deals with Argument objects and native namespaces like Math
  if (!value || typeof value !== 'object' || toString.call(value) !== '[object Object]') {
    return false
  }

  const proto = Object.getPrototypeOf(value)
  if (proto === null) {
    return true
  }

  // Iteratively going up the prototype chain is needed for cross-realm environments (differing contexts, iframes, etc)
  let parentProto = proto
  let nextProto = Object.getPrototypeOf(proto)
  while (nextProto !== null) {
    parentProto = nextProto
    nextProto = Object.getPrototypeOf(parentProto)
  }
  return parentProto === proto
}

/**
 * Contributes additional methods to a constructor
 */
export function mixin(ctor: any, methods: any) {
  const keyCopier = (key: PropertyKey) => {
    ctor.prototype[key] = methods[key]
  }
  Object.keys(methods).forEach(keyCopier)
  Object.getOwnPropertySymbols && Object.getOwnPropertySymbols(methods).forEach(keyCopier)
  return ctor
}

if (require.main === module) {
  console.log(isArrayLike({ length: 0 }))
  console.log(isArrayLike({ length: 1 }))
  console.log(isArrayLike(''))

  class A {
    a = 4
    b = 2
  }
  mixin(A, { a: 1, b: 2 })
  console.log(A.prototype)
}
