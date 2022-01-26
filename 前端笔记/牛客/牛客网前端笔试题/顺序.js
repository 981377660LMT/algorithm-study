async function foo() {
  console.log(2)
  console.log(await Promise.resolve(8))
  console.log(9)
  console.log(await Promise.resolve(18))
  console.log(19)
}
async function bar() {
  console.log(4)
  console.log(await 6)
  console.log(7)
}
console.log(1)
foo()
console.log(3)
bar()
console.log(5)

// 先同步，后异步，await后面全部作为任务被加入异步队列
