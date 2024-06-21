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
