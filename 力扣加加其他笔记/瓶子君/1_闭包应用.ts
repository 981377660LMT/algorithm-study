// setTimeout 传参
//通过闭包可以实现传参效果
const myfunc = (param: number) => {
  return function () {
    console.log(param)
  }
}
const f1 = myfunc(1)
setTimeout(f1, 1000)
// 既然基本类型变量存储在栈中，栈中数据在函数执行完成后就会被自动销毁，那执行函数之后为什么闭包还能引用到函数内的变量？
// 闭包中的变量没有保存在栈中，而是保存到了堆中
