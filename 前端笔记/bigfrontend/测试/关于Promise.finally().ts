// finally() 方法返回一个Promise。在promise结束时，无论结果是fulfilled或者是rejected，都会执行指定的回调函数。
// 这为在Promise是否成功完成后都需要执行的代码提供了一种方式。

// finally的onfinally函数不接收参数 不返回值(中途进来打酱油的 可以去掉)
// 硬要传的话，返回值对后面没影响

Promise.resolve(1)
  .finally(() => console.log(22))
  .then(v => console.log(v)) // 1
