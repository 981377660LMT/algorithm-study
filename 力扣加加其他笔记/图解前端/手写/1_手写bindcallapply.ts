interface Function {
  myApply: (thisArg: any, argArray: any[]) => any
  myCall: (thisArg: any, ...argArray: any[]) => any
  myBind: (thisArg: any, ...argArray: any[]) => any
}

Function.prototype.myCall = function (thisArg: any, ...argArray: any[]): any {
  // 把函数挂到对象上调用
  thisArg.fn = this
  const result = thisArg.fn(...argArray)
  delete thisArg.fn
  return result
}

Function.prototype.myApply = function (thisArg: any, argArray: any[]): any {
  // 把函数挂到对象上调用
  thisArg.fn = this
  const result = thisArg.fn(...argArray)
  delete thisArg.fn
  return result
}

// 注意这里的实现
Function.prototype.myBind = function (thisArg: any, ...argArray: any[]): any {
  return (...restArgs: any[]) => this.call(thisArg, ...argArray, ...restArgs)
}

// test
const aaaa = {
  name: 'name of a',
}
function test(this: any, ...msg: any[]) {
  console.log(this.name)
  console.log(...msg)
}
const t = test.myBind(aaaa, 'hello')
t('world')
