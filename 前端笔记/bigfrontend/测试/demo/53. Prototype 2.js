function F() {
  this.foo = 'bar'
}

const f = new F()
console.log(f.prototype)
console.log(f.__proto__)
console.log(f)
undefined
