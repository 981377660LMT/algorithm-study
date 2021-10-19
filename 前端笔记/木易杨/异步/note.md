1. Promise 构造函数是同步执行还是异步执行，那么 then 方法呢

```JS
const promise = new Promise((resolve, reject) => {
  console.log(1)
  resolve()
  console.log(2)
})

promise.then(() => {
  console.log(3)
})

console.log(4)

```

执行结果是：1243
promise 构造函数是同步执行的，then 方法是异步执行的
