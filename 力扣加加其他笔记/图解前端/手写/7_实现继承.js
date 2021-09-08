function extend(A, B) {
  // 过渡中间人
  function f() {}

  console.log(B.prototype, f.prototype)
  f.prototype = B.prototype
  A.prototype = new f()
  // console.log(A.prototype.constructor) // 会去原型链上找constructor属性 所以找到B
  A.prototype.constructor = A
}

function A(name) {
  this.name = name
}
function B(name) {
  this.name = name
}
extend(A, B)
B.prototype.say = function () {
  console.log('b say')
}
A.prototype.eat = function () {
  console.log('a eat')
}

const a = new A('a name')

console.log(a.name)
a.say()
a.eat()
