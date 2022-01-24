Web Workers 的分类：Dedicated Worker、Shared Worker 和 Service Workers；

1. Chrome 的四个主要进程
   浏览器进程 （Browser Process）
   `渲染进程` （Renderer Process）：也称渲染引擎/浏览器内核。
   插件进程 （Plugin Process）
   GPU 进程 （GPU Process）
2. `渲染进程`中的每个线程
   GUI 渲染线程
   JavaScript 引擎线程
   事件触发线程
   定时触发器线程
   Http 异步请求线程

由于 JavaScript 引擎与 GUI 渲染线程是互斥的，如果 JavaScript 引擎执行了一些计算密集型或高延迟的任务，那么会导致 GUI 渲染线程被阻塞或拖慢。那么如何解决这个问题呢？ —— Web Workers。
在主线程运行的同时，Worker 线程在后台运行，两者互不干扰。等到 Worker 线程完成计算任务，再把结果返回给主线程。

3. Web Workers 的限制与能力
   你不能「直接在 worker 线程中操纵 DOM 元素，或使用 window 对象中的某些方法和属性。」
4. 主线程与 Web Workers 之间的通信
   主线程和 Worker 线程相互之间使用 postMessage() 方法来发送信息，并且通过 onmessage 这个事件处理器来接收信息。数据的交互方式为传递副本，而不是直接共享数据。
5. Web Workers 的分类
   - 专用线程 Dedicated Worker 只能为一个页面所使用
   - 共享线程 Shared Worker 则可以被多个页面所共享。
     例子:在 Shared Worker 的示例页面上有一个 「点赞」 按钮，每次点击时点赞数会加 1。首先你新开一个窗口，然后点击几次。然后新开另一个窗口继续点击，这时你会发现当前页面显示的点赞数是基于前一个页面的点赞数继续累加。
   - Service Workers
     Service workers 本质上充当 Web 应用程序与浏览器之间的代理服务器，也可以在网络可用时作为浏览器和网络间的代理。它们旨在（除其他之外）使得能够创建有效的`离线体验`，`拦截网络请求`并基于网络是否可用以及更新的资源是否驻留在服务器上来采取适当的动作。
     ![Service Worker](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/6/24/172e3d66dc411bbe~tplv-t2oaga2asx-watermark.awebp)
