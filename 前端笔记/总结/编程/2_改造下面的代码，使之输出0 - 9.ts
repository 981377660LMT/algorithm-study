for (var i = 0; i < 10; i++) {
  setTimeout(() => {
    console.log(i)
  }, 1000)
} // 10 10 10 ...
//////////////////////////////////////////////////////////
for (var i = 0; i < 10; i++) {
  setTimeout(console.log, 1000, i)
}
// 利用 setTimeout 函数的第三个参数，会作为回调函数的第一个参数传入

for (let i = 0; i < 10; i++) {
  setTimeout(() => {
    console.log(i)
  }, 1000)
}
// 利用 let 变量的特性 — 在每一次 for 循环的过程中，
// let 声明的变量会在当前的块级作用域里面（for 循环的 body 体，也即两个花括号之间的内容区域）创建一个文法环境（Lexical Environment），
// 该环境里面包括了当前 for 循环过程中的 i，
