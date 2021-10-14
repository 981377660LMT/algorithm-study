class Foo {
  constructor() {
    this.name = 'a'
  }
}

let obj = {}
obj = Foo.call(obj)
// TypeError: Class constructor Foo cannot be invoked without 'new'
console.log(obj)
//////////////////////////////////////////////////////////////////////
// function Foo() {
//   this.name = 'a'
//   return this
// }

// let obj = {}
// obj = Foo.call(obj)
// console.log(obj)
// // { name: 'b' }

// ES6 classes should be only called with new, Constructor.call results in error.
