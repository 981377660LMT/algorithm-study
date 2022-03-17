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

2. setTimeout、Promise、Async/Await 的区别

- 在执行上下文栈的同步任务执行完后；
- 首先执行 Microtask 队列，按照队列先进先出的原则，一次执行完所有 Microtask 队列任务；
- 然后执行 Macrotask/Task 队列，**一次执行一个，一个执行完后，检测 Microtask 是否为空**；
- 为空则执行下一个 Macrotask/Task；
- 不为空则执行 Microtask

3. 在不用 promise.all()的情况下，在并发 100 次请求后，如何把浏览器在不同时间点获取到的数据按原先的数据逻辑顺序输出。
4. .then 或者 .catch 中 return 一个 error 对象并不会抛出错误，所以不会被后续的 .catch 捕获
5. requestAnimationFrame 属于宏任务还是微任务
6. 微任务包括
   MutationObserver、Promise.then()或 catch()、Promise 为基础开发的其它技术，比如 fetch API、V8 的垃圾回收过程、Node 独有的 process.nextTick。
   宏任务包括：**script** 、setTimeout、setInterval 、setImmediate 、I/O 、UI rendering。
   注意 ⚠️：在所有任务开始的时候，由于宏任务中包括了 script，所以浏览器会先执行一个宏任务，在这个过程中你看到的延迟任务(例如 setTimeout)将被放到下一轮宏任务中来执行。
7. promise 例题
   https://juejin.cn/post/6844904077537574919#heading-46

```JS
const promise = new Promise((resolve, reject) => {
  setTimeout(() => {
    console.log('timer')
    resolve('success')
  }, 1000)
})
const start = Date.now();
promise.then(res => {
  console.log(res, Date.now() - start)
})
promise.then(res => {
  console.log(res, Date.now() - start)
})

'timer'
'success' 1001
'success' 1002

Promise 的 .then 或者 .catch 可以被调用多次，但这里 Promise 构造函数只执行一次。或者说 promise 内部状态一经改变，并且有了一个值，那么后续每次调用 .then 或者 .catch **都会直接拿到该值**。  (共用一个值)
```

```JS
Promise.resolve().then(() => {
  return new Error('error!!!')
}).then(res => {
  console.log("then: ", res)
}).catch(err => {
  console.log("catch: ", err)
})

"then: " "Error: error!!!"
返回任意一个非 promise 的值都会被包裹成 promise 对象，因此这里的return new Error('error!!!')也被包裹成了return Promise.resolve(new Error('error!!!'))。
当然如果你想抛出一个错误的话，可以用下面👇两的任意一种：
return Promise.reject(new Error('error!!!'));
// or
throw new Error('error!!!')
```

```JS
const promise = Promise.resolve().then(() => {
  return promise;
})

promise.catch(console.err)
.then 或 .catch 返回的值不能是 promise 本身，否则会造成死循环。

Uncaught (in promise) TypeError: Chaining cycle detected for promise #<Promise>
```

```JS
Promise.resolve(1)
  .then(2)
  .then(Promise.resolve(3))
  .then(console.log)
.then 或者 .catch 的参数期望是函数，传入非函数则会发生值透传。
所以输出结果为：

1

```

**Promise.finally()**
其实你只要记住它三个很重要的知识点就可以了：

1. .finally()方法不管 Promise 对象最后的状态如何都会执行
2. .finally()方法的回调函数不接受任何的参数，也就是说你在.finally()函数中是没法知道 Promise 最终的状态是 resolved 还是 rejected 的
3. 它最终返回的默认会是一个上一次的 Promise 对象值，不过如果抛出的是一个异常则返回异常的 Promise 对象

```JS
Promise.resolve('1')
  .then(res => {
    console.log(res)
  })
  .finally(() => {
    console.log('finally')
  })

Promise.resolve('2')
  .finally(() => {
    console.log('finally2')
  	return '我是finally2返回的值'
  })
  .then(res => {
    console.log('finally2后面的then函数', res)
  })

'1'
'finally2'
'finally'
'finally2后面的then函数' '2'

因为每次执行完then 都返回新的微任务 入队 队列不断shift push 注意顺序
```

finally 中要是抛出的是一个异常是怎样的：

```JS
Promise.resolve('1')
  .finally(() => {
    console.log('finally1')
    throw new Error('我是finally中抛出的异常')
  })
  .then(res => {
    console.log('finally后面的then函数', res)
  })
  .catch(err => {
    console.log('捕获错误', err)
  })
'finally1'
'捕获错误' Error: 我是finally中抛出的异常
但是如果改为return new Error('我是finally中抛出的异常')，打印出来的就是'finally后面的then函数 1'
```

Promise.all 使用场景
一些游戏类的素材比较多的应用，打开网页时，预先加载需要用到的各种资源如图片、flash 以及各种静态文件。所有的都加载完后，我们再进行页面的初始化。
Promise.race 使用场景
我们可以用 race 给某个异步请求设置超时时间，并且在超时后执行相应的操作

async await
**紧跟着 await 后面的语句相当于放到了 new Promise 中，下一行及之后的语句相当于放在 Promise.then 中**

- 异步阻塞 与 异步并发
  阻塞和非阻塞的概念：
  阻塞：等待的过程中线程被挂起，CPU 资源被分配到其他地方去
  非阻塞：等待的过程 CPU 资源还在该线程中，线程还能做其他的事情
  同步与异步的概念：
  同步：同步时调用者主动等待调用过程，且能立即得到结果的。
  异步：程序发出调用之后，无法立即得到结果，需要额外的操作才能得到预期的结果是为异步。
