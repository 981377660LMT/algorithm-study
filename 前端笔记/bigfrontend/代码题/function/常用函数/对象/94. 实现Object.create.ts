/**
 * @param {any} proto
 * @return {object}
 * 1.Object.create的第二个参数不用实现，忽略即可。
 * 2.如果传入的不是object，请throw error。
 * 3.请不要直接使用Object.create() 和 Object.setPrototypeOf()
 */
function myObjectCreate(proto: any): object {
  function Constructor() {}
  Constructor.prototype = proto.prototype || proto
  // @ts-ignore
  return new Constructor() // new 出来的对象的obj.__proto__为contructor.prototype
}

function myObjectCreate2(proto: any): object {
  if (typeof proto !== 'object' || proto === null) throw new Error('')
  const obj = {} as any
  obj.__proto__ = proto //  Object.setPrototypeOf(obj,proto)
  return obj
}

console.log(myObjectCreate(null))

const proto = {}
const a = myObjectCreate(proto)
expect(Object.getPrototypeOf(a)).toBe(proto)

// function Foo(){}
// Foo.prototype = null;
// console.log(new Foo().toString); //outputs function toString() { [native code] }
// 此时new的内部相当于Object.create(Foo) 而不是Object.create(Foo.prototype)

// function Foo(){}
// Foo.prototype = Object.create(null);
// console.log(new Foo().toString); //output undefined
