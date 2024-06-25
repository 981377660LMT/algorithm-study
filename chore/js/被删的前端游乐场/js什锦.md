js 什锦
https://godbasin.github.io/categories/#js%E4%BB%80%E9%94%A6

# js 判断某个位置是否特定元素

1. 常用坐标相关属性

```js
scrollHeight //获取对象的滚动高度
scrollLeft //设置或获取位于对象左边界和窗口中目前可见内容的最左端之间的距离
scrollTop //设置或获取位于对象最顶端和窗口中可见内容的最顶端之间的距离
scrollWidth //获取对象的滚动宽度
offsetHeight //获取对象相对于版面或由父坐标 offsetParent 属性指定的父坐标的高度
offsetLeft //获取对象相对于版面或由 offsetParent 属性指定的父坐标的计算左侧位置
offsetTop //获取对象相对于版面或由 offsetTop 属性指定的父坐标的计算顶端位置
event.clientX //相对文档的水平座标
event.clientY //相对文档的垂直座标
event.offsetX //相对容器的水平坐标
event.offsetY //相对容器的垂直坐标
document.documentElement.scrollTop 垂直方向滚动的值
event.clientX + document.documentElement.scrollTop //相对文档的水平座标+垂直方向滚动的量
```

2. 绑定鼠标事件

3. 获取鼠标坐标

4. 计算鼠标位置

5. 获取当前位置是否有特定元素
   给需要检测的元素绑定 id 或者自定义属性
   通过不断获取当前元素父元素，直至获取成功（通过自定义属性判断）或者元素为 body
   ```js
   function fnGetTable(oEl) {
     while (null !== oEl && $(oEl).attr('自定义属性') !== '特定属性值' && target.tagName !== 'BODY') {
       oEl = oEl.parentElement
     }
     return oEl
   }
   ```
6. 实例：下拉菜单点击外围自动关闭

# 闭包

闭包会使子函数保持其作用域链的所有变量及函数与内存中，内存消耗很大，所以不能滥用，并且在使用的时候尽量销毁父函数不再使用的变量哦。

# json 格式化

https://godbasin.github.io/2016/11/13/json-to-html-1-use-string-regular/
https://godbasin.github.io/2016/11/13/json-to-html-2-use-object/

1. 方法 1: 分析 JSON.stringify()后的字符串，使用正则把需要的格式匹配替换
2. 方法 2: 将 json 转化为 object，然后通过 js 判断数据类型进行格式化

# 代码调试

我们写代码的时候主要通过两种方式进行调试：代码中调试、浏览器中调试。

- alert()、console、debugger

```js
// 这样可以顺便把变量名也打印出来
console.log({ data, list, state })
```

- 浏览器(Chrome)调试打断点

1. 代码断点
2. 事件断点
   元素上事件断点：某一事件/某类事件，从 Elements > Event Listeners 中进行。
3. DOM 断点
   一般是 DOM 结构改变时触发。有时候我们需要监听某个 DOM 被修改情况，而不关心是哪行代码做的修改（也可能有多处都会对其做修改），可以直接在 DOM 上设置断点。
   在元素审查的 Elements 标签页中在某个元素上右键菜单里可以设置三种不同情况的断点：

   - 子节点修改
   - 自身属性修改
   - 自身节点被删除

4. XHR 断点
   右侧调试区有一个 XHR Breakpoints，点击+ 并输入 URL 包含的字符串，即可监听该 URL 的 Ajax 请求，输入内容就相当于 URL 的过滤器。
   一旦 XHR 调用触发时就会在 request.send()的地方中断。
5. 全局事件断点
   对事件发生时断点，不会限定在元素上，只要是事件发生，并且有 handler 就断点。
   还可以对 resize、ajax、setTimeout/setInterval 断点.

---

堆排序缺点：cache miss

- 单页应用：
  - 单页应用的好处是：
    页面的数据状态都能维持着。
    部分擦除重绘，比整个页面刷新的效果体验要好很多。
  - 单页应用也会有缺点：
    不利于 SEO（最大痛处）。
    请求等待时间长。

dom 事件捕获阶段：top -> bottom
dom 事件冒泡阶段：bottom -> top

事件委托(delegate)的思想
框架会帮我们用事件委托的方式处理掉，

为什么 javascript 是单线程的：
更多是因为对页面交互的同步处理。作为浏览器脚本语言，JavaScript 的主要用途是与用户互动，以及操作 DOM，若是`多线程会导致严重的同步问题`

---

html 加载

正常的网页加载流程是这样的：

1. 浏览器一边下载 HTML 网页，一边开始解析
2. 解析过程中，发现<script>标签
3. 暂停解析，网页渲染的控制权转交给 JavaScript 引擎
4. 如果<script>标签引用了外部脚本，就下载该脚本，否则就直接执行
5. 执行完毕，控制权交还渲染引擎，恢复往下解析 HTML 网页

**将 js 放在 body 的最后面，可以避免资源阻塞，同时使静态的 html 页面迅速显示。**
如果外部脚本加载时间很长（比如一直无法完成下载），就会造成网页长时间失去响应，浏览器就会呈现“假死”状态，这被称为“阻塞效应”。
`html 需要等 head 中所有的 js 和 css 加载完成后才会开始绘制，但是 html 不需要等待放在 body 最后的 js 下载执行就会开始绘制。`

**将 css 放在 head 里，可避免浏览器渲染的重复计算。**
经过上面的渲染过程，我们知道 Layout 的计算是比较消耗性能的，所以我们在开始计算 Render Tree 之前，就把所有的 css 文件拿到，`这样可减少 Repaint 和 Reflow。`
