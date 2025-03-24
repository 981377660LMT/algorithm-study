// finally() 方法返回一个Promise。在promise结束时，无论结果是fulfilled或者是rejected，都会执行指定的回调函数。
// 这为在Promise是否成功完成后都需要执行的代码提供了一种方式。

// 1. finally() never receive an argument
// !finally的onfinally函数不接收参数 (中途进来打酱油的)
// 2. Normal return value in finally won't make effect on promise object
// !硬要传的话，返回值对后面没影响
// 3. throw Error in finally()
// !Note: A throw (or returning a rejected promise) in the finally callback will reject the new promise with the rejection reason specified when calling throw.

Promise.resolve(1)
  .finally(() => {
    console.log(22)
    // return Promise.resolve(2)
    return Promise.reject(2)
  })
  .then(v => console.log(v)) // 1
