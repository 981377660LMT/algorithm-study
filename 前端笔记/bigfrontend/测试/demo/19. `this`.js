const obj = {
  a: 1,
  b: function () {
    console.log(this.a)
  },
  c() {
    console.log(this.a)
  },
  d: () => {
    console.log(this.a)
  },
  e: (function () {
    return () => {
      console.log(this.a)
    }
  })(),
  f: function () {
    return () => {
      console.log(this.a)
    }
  },
}

console.log(obj.a) // 1
obj.b() // 1
obj.b() //1
const b = obj.b
b() // undefined because reference is no longer available
obj.b.apply({ a: 2 }) // 2 {a:2} passed as this
obj.c() // 1
obj.d() // undefined arrow functions access lexical scope and a doesn't exist in lexical scope
obj.d() // undefined arrow functions access lexical scope and a doesn't exist in lexical scope
obj.d.apply({ a: 2 }) // undefined arrow functions access lexical scope and a doesn't exist in lexical scope
obj.e() // undefined it is similar to obj.d as its using an IIFE
obj.e() // undefined
obj.e.call({ a: 2 }) // undefined
obj.f()() // 1 because of closure lexical scope of inner arrow function would of object
obj.f()() // 1
obj.f().call({ a: 2 }) // 1
