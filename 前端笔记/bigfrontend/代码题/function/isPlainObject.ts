export const isUndefined = (obj: any): obj is undefined => typeof obj === 'undefined'
export const isNil = (obj: any): obj is null | undefined => isUndefined(obj) || obj === null
export const isObject = (fn: any): fn is object => !isNil(fn) && typeof fn === 'object'

/**
 *
 * @param obj
 * @returns
 * @description
 * 纯对象：created by the {} object literal notation or constructed with new Object()
 * redux 从 4.0.0 开始在测试中使用了 isPlainObject
 * 先判断 obj 本身是否满足我们熟悉的合法对象概念；
   再判断 obj 的原型属性是不是 [Object: null prototype] {}
 */
export const isPlainObject = (obj: any): obj is object => {
  if (typeof obj !== 'object' || obj === null) return false

  let proto = obj
  while (Object.getPrototypeOf(proto) !== null) {
    proto = Object.getPrototypeOf(proto)
  }

  return Object.getPrototypeOf(obj) === proto
}

console.log(isPlainObject(Math)) // An intrinsic object ,true
console.log(isPlainObject({})) // An intrinsic object ,true
// console.log(typeof Date)  // function ,false
console.log(isPlainObject(Object.create({ 1: 'l' }))) // false

class AAA {}
class BBB extends AAA {}
console.log(isPlainObject(new BBB())) // false
