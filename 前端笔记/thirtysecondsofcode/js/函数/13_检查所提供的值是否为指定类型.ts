// 1.确保该值不是未定义的或 null (否则报错)
// 2.值上的构造函数属性与类型进行比较 (回到原型链找constructor)
const is = <T>(type: T, val: any): val is T => {
  return ![undefined, null].includes(val) && val.constructor === type
}

is(Array, [1]) // true
is(ArrayBuffer, new ArrayBuffer(1)) // true
is(Map, new Map()) // true
is(RegExp, /./g) // true
is(Set, new Set()) // true
is(WeakMap, new WeakMap()) // true
is(WeakSet, new WeakSet()) // true
is(String, '') // true
is(String, new String('')) // true
is(Number, 1) // true
is(Number, new Number(1)) // true
is(Boolean, true) // true
is(Boolean, new Boolean(true)) // true

export {}

// @ts-ignore
console.log(NaN.constructor)
