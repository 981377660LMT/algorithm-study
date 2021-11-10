var b = 10
;(function b() {
  b = 20
  console.log(b) // [Function b]
  console.log(window.b) // 10，不是20
})()

// 非匿名自执行函数，函数名只读。
// [Function: b]

// 简单改造，使之分别打印 10 和 20
// 打印 20 :
var b = 10
;(function () {
  b = 20
  console.log(b)
})()

// 打印 10:
var b = 10
;(function () {
  console.log(b)
  b = 20
})()
