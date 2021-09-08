// JavaScript 采用词法作用域(lexical scoping)，也就是静态作用域。
// 函数的作用域在函数定义的时候就决定了。
// 而与词法作用域相对的是动态作用域，函数的作用域是在函数调用的时候才决定的。
var value = 1

function foo() {
  console.log(value)
}

function bar() {
  var value = 2
  foo()
}

bar()

// 结果是 ???

// 假设JavaScript采用动态作用域  结果会是2
// 假设JavaScript采用动态作用域，让我们分析下执行过程：

// 执行 foo 函数，依然是从 foo 函数内部查找是否有局部变量 value。
// 如果没有，就从调用函数的作用域，也就是 bar 函数内部查找 value 变量，所以结果会打印 2。

// 也许你会好奇什么语言是动态作用域？
// bash 就是动态作用域，
var scope = 'global scope'
function checkscope() {
  var scope = 'local scope'
  function f() {
    return scope
  }
  return f()
}
console.log(checkscope())
