node 中事件循环的实现是**依靠的 libuv 引擎**
浏览器和 Node 环境下，microtask 任务队列的执行时机不同
Node 端，**微任务** 在事件循环的**各个阶段之间**执行
浏览器端，**微任务** 在事件循环的 **宏任务** 执行完之后执行
nodejs 事件循环六部

**定**时器检测阶段(timers)：本阶段执行 **timer** 的回调，即 setTimeout、setInterval 里面的回调函数。
I/O 事件**回**调阶段(I/O callbacks)：fs.readFile(path,cb)。
**闲**置阶段(idle, prepare)：仅系统内部使用。
**轮**询阶段(poll)：问操作系统准备好了没；检索新的 I/O 事件;执行与 I/O 相关的回调（几乎所有情况下，除了关闭的回调函数，那些由计时器和 setImmediate() 调度的之外），其余情况 node 将在适当的时候在此阻塞。
**检**查阶段(check)：**setImmediate**() 回调函数在这里执行
**关**闭事件回调阶段(close callback)：一些关闭的回调函数，如：socket.on(‘close’, …)。
`定 回 闲 轮 检 关`

node 的事件循环的阶段顺序为：
输入数据阶段(incoming data)->轮询阶段(poll)->检查阶段(check)->关闭事件回调阶段(close callback)->定时器检测阶段(timers)->I/O 事件回调阶段(I/O callbacks)->闲置阶段(idle, prepare)->轮询阶段…
日常开发中的绝大部分异步任务都是在 poll、check、timers 这 3 个阶段处

**区别**
浏览器只有两个队列:微任务，宏任务队列
NodeJs 有 6 个

nextTick 是推迟到任务队列末

- 浏览器和 Node.js 的事件循环机制有什么区别？
