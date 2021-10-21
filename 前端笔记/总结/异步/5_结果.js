// 请问最后的输出结果是什么？  A 在所有同步任务执行完之前，任何的异步任务是不会执行的
console.log('A')
setTimeout(function () {
  console.log('B')
}, 0)
while (true) {}
