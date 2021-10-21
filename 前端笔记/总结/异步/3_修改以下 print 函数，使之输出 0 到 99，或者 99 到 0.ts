function print(n: number) {
  setTimeout(() => {
    console.log(n)
  }, Math.floor(Math.random() * 1000))
}
for (var i = 0; i < 100; i++) {
  print(i)
}
// 1.只能修改 setTimeout 到 Math.floor(Math.random() * 1000 的代码
// 2、不能修改 Math.floor(Math.random() * 1000

function print(n: number) {
  setTimeout(
    () => {
      console.log(n)
    },
    1,
    Math.floor(Math.random() * 1000)
  )
}
for (var i = 0; i < 100; i++) {
  print(i)
}
