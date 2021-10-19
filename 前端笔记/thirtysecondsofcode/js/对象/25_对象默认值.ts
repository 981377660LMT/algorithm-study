// 为未定义的对象中的所有属性分配默认值。
console.log(defaults({ a: 1 }, { b: 2 }, { b: 6 }, { a: 3 })) // { a: 1, b: 2 })

function defaults(obj: Record<any, any>, ...defaults: Record<any, any>[]): any {
  return Object.assign({}, obj, ...defaults.reverse(), obj)
}
