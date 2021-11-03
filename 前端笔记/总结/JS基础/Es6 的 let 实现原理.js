// 原始 es6 代码
var funcs = []
for (let i = 0; i < 10; i++) {
  funcs[i] = function () {
    console.log(i)
  }
}
funcs[0]() // 0

// babel 编译之后的 es5 代码（polyfill）
// 用闭包保存状态i
var _loop = function _loop(i) {
  funcs[i] = function () {
    console.log(i)
  }
}

for (var i = 0; i < 10; i++) {
  _loop(i)
}
funcs[0]() // 0

// let 是借助闭包和函数作用域来实现块级作用域的效果的
