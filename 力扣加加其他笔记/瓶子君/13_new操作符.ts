const myNew = (fn: (...args: any[]) => any, ...args: any[]) => {
  const obj = Object.create(fn.prototype)
  const res = fn.apply(obj, args)
  return res instanceof Object ? res : obj
}

export {}
