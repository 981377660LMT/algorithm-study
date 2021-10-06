// 使之输出 0 到 99，或者 99 到 0
// function print(n: number) {
//   setTimeout(() => {
//     console.log(n)
//   }, Math.floor(Math.random() * 1000))
// }

// 方法1, 利用setTimeout、setInterval的第三个参数,第三个以后的参数是作为第一个func()的参数传进去。
// function print(n: number) {
//   setTimeout(
//     a => {
//       console.log(n, a)
//     },
//     1,
//     Math.floor(Math.random() * 1000)
//   )
// }

// 方法二：
// 修改 setTimeout 第一个函数参数
function print(n: number) {
  setTimeout(
    (() => {
      console.log(n)
      return () => {}
    }).call(n),
    Math.floor(Math.random() * 1000)
  )
}
for (var i = 0; i < 100; i++) {
  print(i)
}
