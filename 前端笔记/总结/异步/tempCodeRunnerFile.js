function wait() {
//   return new Promise(resolve => setTimeout(resolve, 10 * 100))
// }

// async function main() {
//   console.time()
//   const x = wait()
//   const y = wait()
//   const z = wait()
//   await x
//   await y
//   await z
//   console.timeEnd()
// }
// main()
// // Promise 内部代码块是按序执行。和普通的代码执行顺序没区别。只是在调用then()时 返回了代码块中resolve()返回的结果，其他没有任何buf加成。
// // 三个任务发起的时候没有await，可以认为是同时发起了三个异步。之后各自await任务的结果。结果按最高耗时计算，由于三个耗时一样。所以结果是 10 * 1000ms
// function wait() {
//   return new Promise(resolve => {
//     console.log(1)
//     resolve()
//     setTimeout(function () {
//       console.log(3)
//       // resolve();
//       console.log(4)
//     }, 3 * 1000)
//     console.log(2)
//   })
// }

// async function main() {
//   console.time()
//   const x = wait()
//   const y = wait()
//   const z = wait()
//   await x
//   await y
//   await z
//   console.timeEnd()
// }
// main()