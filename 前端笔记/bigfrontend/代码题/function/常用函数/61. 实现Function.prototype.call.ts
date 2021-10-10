interface Function {
  mycall: (thisArg: any, ...args: any[]) => any
}

Function.prototype.mycall = function (thisArg: any, ...args: any[]) {
  thisArg = thisArg || globalThis // 不传也可
  thisArg = Object(thisArg) // transform primitive value
  const func = Symbol()
  thisArg[func] = this
  const res = thisArg[func](...args)
  delete thisArg[func]
  return res
}
