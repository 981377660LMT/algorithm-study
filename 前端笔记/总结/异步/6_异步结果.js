var date = new Date()

console.log(1, new Date() - date)

setTimeout(() => {
  console.log(2, new Date() - date)
}, 500)

Promise.resolve().then(console.log(3, new Date() - date)) // 参数不是函数,直接执行
// Promise.resolve().then(() => console.log(3, new Date() - date)) //

while (new Date() - date < 1000) {}

console.log(4, new Date() - date)

// 1 0
// 3 7
// 4 1000
// 2 1001
