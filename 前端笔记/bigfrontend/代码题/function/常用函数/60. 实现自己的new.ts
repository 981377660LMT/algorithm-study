/**
 * @param {Function} constructor
 * @param {any[]} args - argument passed to the constructor
 * `myNew(constructor, ...args)` should return the same as `new constructor(...args)`
 * @description
 * 通常构造函数不返回值，但是如果它们想覆盖正常的对象创建过程，它们可以选择返回值。
 */
const myNew = (constructor: Function, ...args: any[]) => {
  // your code here
  const obj = Object.create(constructor.prototype)
  const res = constructor.apply(obj, args)
  if (typeof res === 'object') return res
  else return obj
}
