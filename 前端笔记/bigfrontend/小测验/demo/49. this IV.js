var bar = 1

function foo() {
  return this.bar++
}

const a = {
  bar: 10,
  foo1: foo,
  foo2: function () {
    return foo()
  },
}

console.log(a.foo1.call()) //1 by default global this is passed using call so it will access global variable var
console.log(a.foo1()) // 10 context of object will be shared
console.log(a.foo2.call()) // 2 by default global this is passed using call so it will access global variable var
console.log(a.foo2()) // 3 since foo is called from function inside foo2 so object context is removed and it will have global this
