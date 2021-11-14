https://tsejx.github.io/javascript-guidebook/browser-object-model/binary-data-and-files/application

1. 抛开全局变量会成为 window 对象的属性不谈，定义全局变量与在 Window 对象上直接定义属性还是有一点差别：全局变量不能通过 delete 操作符删除，而直接在 window 对象上的定义的属性可以。
   这是因为，通过 var 语句添加的 window 属性有一个名为 [[Configurable]] 的特性，这个特性的值被设置为 false，因此这样定义的属性不可以通过 delete 操作符删除
2. 传统动画渲染的弊端( setTimeout 和 setInterval )
   - 定时器的**第二个时间参数只是指定了多久后将动画任务添加到浏览器的 UI 线程队列中**，`如果 UI 线程处于忙碌状态，那么动画不会立即执行。`
   - 动画的时间间隔不好确定，设置时间过长会使得动画不够平滑流畅，设置过短会令浏览器的重绘频率容易达到瓶颈（推荐最佳循环间隔是 17ms，因为大多数电脑的显示器刷新频率是 60Hz，1000ms/60）。
3. requestIdleCallback 与 时间切片
   一般浏览器的刷新率为 60HZ，即 1 秒钟刷新 60 次。1000ms / 60hz = 16.6，大概每过 16.6ms 浏览器会渲染一帧画面。
   在这段时间内，浏览器大体会做两件事：task 与 render。
   task -> requestAnimationFrame -> render -> requestIdleCallback
   如果 task 执行时间超过了 16.6ms（比如 task 中有个很耗时的 while 循环）。**那么这一帧就没有时间 render，页面直到下一帧 render 后才会更新。表现为页面卡顿一帧，或者说掉帧**。

   最好的办法是时间切片，把长时间 **task 分割为几个短时间 task**。
   为了解决掉帧造成的卡顿，React16 将递归的构建方式改为可中断的遍历。React16 就是基于 requestIdleCallbackAPI，实现了自己的 Fiber Reconciler。
   **以 5ms 的执行时间划分 task，每遍历完一个节点，就检查当前 task 是否已经执行了 5ms**。
   如果超过 5ms，则中断本次 task。
   `通过将 task 执行时间切分为一个个小段，减少长时间 task 造成无法 render 的情况，这就是时间切片`。

4. 二进制数据和文件 API
   历史上，JavaScript 无法处理二进制数据。如果一定要处理的话，只能使用 String.prototype.charCodeAt() 方法，逐个地将字节从文字编码转成二进制数据，还有一种办法是将二进制数据转成 Base64 编码，再进行处理。这两种方法不仅速度慢，而且容易出错。因此 ECMAScript 5 引入了 Blob 对象，允许直接操作二进制数据。

- Blob 对象：二进制数据基本对象，在它的基础上，又衍生出一系列相关的 API，用于操作文件
- File 对象：负责处理那些以文件形式存在的二进制数据，也就是操作本地文件
- FileList 对象：File 对象的网页表单接口
- FileReader 对象：负责将二进制数据读入内存内容
- URL 对象：用于对二进制数据生成 URL
- FormData 对象：读取页面表单项文件数据

5. 浏览器架构
6. 实现一个大文件上传和断点续传
