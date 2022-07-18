setImmediate(function () {
  console.log(1) // 8
}, 0)

// setTimeout的优先级高于setIImmediate
setTimeout(function () {
  console.log(2) // 7
}, 0)

new Promise(function (resolve) {
  console.log(3) // 1
  resolve()
  console.log(4) // 2
}).then(function () {
  console.log(5) // 6
})

console.log(6) // 3
// 优先级process.nextTick 高于 Promise

process.nextTick(function () {
  console.log(7) // 5
})

console.log(8) // 4

// 3
// 4
// 6
// 8
// 7
// 5
// 2
// 1

// macro-task: script (整体代码)，setTimeout, setInterval, setImmediate, I/O, UI rendering.
// micro-task: process.nextTick, Promise(原生)，Object.observe，MutationObserver
// 除了script整体代码，micro-task的任务优先级高于macro-task的任务优先级。
// 其中，script(整体代码) ，可以理解为待执行的所有代码。
