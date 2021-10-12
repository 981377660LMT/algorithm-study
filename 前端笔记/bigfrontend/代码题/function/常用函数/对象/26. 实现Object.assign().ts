// /**
//  * @param {any} target
//  * @param {any[]} sources
//  * @return {object}
//  */
function objectAssign(target: Record<PropertyKey, any>, ...sources: any[]): any {
  // your code here
  if (target == undefined) throw new Error('invalid target')
  if (typeof target !== 'object') {
    const constructor = Object.getPrototypeOf(target).constructor
    target = new constructor(target)
  }

  for (const source of sources) {
    if (source == undefined) continue
    Object.defineProperties(target, Object.getOwnPropertyDescriptors(source))
    for (const symbol of Object.getOwnPropertySymbols(source)) {
      target[symbol] = source[symbol]
    }
  }

  return target
}

// @ts-ignore
console.log(Object.getPrototypeOf(1).constructor)

// TypeError: Cannot convert undefined or null to object
console.log(Object.getOwnPropertyDescriptors(null))

// 实际上
// Object.assign() 处理的是可枚举属性，所以getters不会被复制，不可枚举属性被忽略。
