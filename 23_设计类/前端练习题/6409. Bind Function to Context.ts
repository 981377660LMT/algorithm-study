type Fn = (...args) => any

declare global {
  interface Function {
    bindPolyfill(obj: Record<any, any>): Fn
  }
}

Function.prototype.bindPolyfill = function (obj) {
  return (...args) => {
    const fn = this
    const key = Symbol()
    obj[key] = fn
    const res = obj[key](...args)
    delete obj[key]
    return res
  }
}

export {}
