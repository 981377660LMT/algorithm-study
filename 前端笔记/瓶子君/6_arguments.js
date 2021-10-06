// arguments.callee.name 指当前正在调用的函数的名称
// arguments.callee.caller.name 是指调用当前执行函数的函数的名称
function myFunc() {
  console.log(arguments)
  console.log(arguments.callee.name) // myFunc
  console.log(arguments.callee.caller.name)
}

;(i => myFunc(i, 89))(1)
