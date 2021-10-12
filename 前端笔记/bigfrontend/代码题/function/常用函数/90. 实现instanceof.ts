// from @angular/core
interface Class<T = any> extends Function {
  new (...args: any[]): T
}

function myInstanceOf(obj: any, target: Class): boolean {
  if (obj == null || typeof obj !== 'object') return false
  if (typeof target !== 'function') return false

  while (obj != null) {
    // getPrototype of in place of __proto__
    const proto = Object.getPrototypeOf(obj)
    if (proto === target.prototype) return true
    obj = proto
  }

  return false
}

class A {}
class B extends A {}

const b = new B()
myInstanceOf(b, B) // true
myInstanceOf(b, A) // true
myInstanceOf(b, Object) // true

// function C() {}
// myInstanceOf(b, C) // false
// C.prototype = B.prototype
// myInstanceOf(b, C) // true
// C.prototype = {}
// myInstanceOf(b, C) // false
