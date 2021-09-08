let timer = null
function mySetInterval(callback: (...args: any[]) => any, interval = 1000) {
  const run = () => {
    callback()
    timer = setTimeout(run, interval)
  }
  timer = setTimeout(run, interval)
}

// let a = mySettimeout(() => {
//   console.log(111)
// }, 1000)
// let cancel = mySettimeout(() => {
//   console.log(222)
// }, 1000)
// cancel()

// setInterval 缺点 与 setTimeout 的不同
// 每个 setTimeout 产生的任务会直接 push 到任务队列中；
// 而 setInterval 在每次把任务 push 到任务队列前，
// 都要进行一下判断(看上次的任务是否仍在队列中，如果有则不添加，没有则添加)。

// 定时器指定的时间间隔，表示的是何时将定时器的代码添加到消息队列，而不是何时执行代码。
// 所以真正何时执行代码的时间是不能保证的，取决于何时被主线程的事件循环取到，并执行。

// setInterval 有两个缺点：

// 使用 setInterval 时，某些间隔会被跳过；
// 可能多个定时器会连续执行；
for (var i = 0; i < 5; i++) {
  setTimeout(function () {
    console.log(i)
  }, 1000)
}
// 一秒后立即输出 5 个 5
