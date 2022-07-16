console.log(foo()) // 3

function foo() {
  console.log(1)
}

console.log(foo()) // 3

var foo = 2

function foo() {
  console.log(3)
}

console.log(foo) // 2

foo()
// TypeError: foo is not a function

// 函数提升在变量提升之前
// 函数声明被提升时，声明和赋值两个步骤都会被提升，而普通变量却只能提升声明步骤，而不能提升赋值步骤。

// function foo  // => 声明一个function foo
// function foo
// var foo  // => 声明一个变量 foo
// foo = {  // function foo 初始化
//   console.log(1)
// }
// foo = {
//   console.log(3)
// }

// foo=2  // 变量初始化没有被提升，还在原位
