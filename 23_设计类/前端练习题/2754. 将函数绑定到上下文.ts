/* eslint-disable no-extend-native */
type Fn = (...args) => any

declare global {
  interface Function {
    bindPolyfill(obj: Record<any, any>): Fn
  }
}

// !把函数当成对象的方法执行.
Function.prototype.bindPolyfill = function (obj) {
  const fn = this
  const key = Symbol('key')
  return (...args) => {
    obj[key] = fn
    const res = obj[key](...args)
    delete obj[key]
    return res
  }
}

export {}
