1. 介绍下重绘和回流（Repaint & Reflow），以及如何进行优化

   1. 浏览器渲染机制
      浏览器采用流式布局模型（Flow Based Layout）
      浏览器会把 HTML 解析成 DOM，把 CSS 解析成 CSSOM，DOM 和 CSSOM 合并就产生了渲染树（Render Tree）。
      有了 RenderTree，我们就知道了所有节点的样式，然后计算他们在页面上的大小和位置，最后把节点绘制到页面上。
      由于浏览器使用流式布局，对 Render Tree 的计算通常只需要遍历一次就可以完成，但 table 及其内部元素除外，他们可能需要多次计算，通常要花 3 倍于同等元素的时间，这也是为什么要避免使用 table 布局的原因之一。
   2. 重绘：不会影响布局的变化例如 outline, visibility, color、background-color
   3. 回流：布局或者几何属性需要改变就称为回流
   4. 浏览器优化：现代浏览器大多都是通过队列机制来批量更新布局，浏览器会把修改操作放在队列中，至少一个浏览器刷新（即 16.6ms）才会清空队列，但当你获取布局信息的时候（offsetTop,width 等），队列中可能有会影响这些属性或方法返回值的操作，即使没有，浏览器也会强制清空队列，触发回流与重绘来确保返回正确的值。
   5. 减少回流与重绘
      使用 transform 替代 top
      使用 visibility 替换 display: none ，因为前者只会引起重绘，后者会引发回流（改变了布局
      CSS3 硬件加速（GPU 加速），使用 css3 硬件加速，可以让 transform、opacity、filters 这些动画不会引起回流重绘 。
      JS 批量修改 DOM

2. 分析比较 opacity: 0、visibility: hidden、display: none 优劣和适用场景
   display: none (不占空间，不能点击，文档回流)（场景，显示出原来这里不存在的结构）
   visibility: hidden（占据空间，不能点击，元素重绘）（场景：显示不会导致页面结构发生变动，不会撑开）
   opacity: 0（占据空间，可以点击，重绘，性能消耗较少）（场景：可以跟 transition 搭配）
3. 已知如下代码，如何修改才能让图片宽度为 300px ？注意下面代码不可修改

```HTML
<img src="1.jpg" style="width:480px!important;”>
```

max-width: 300px
transform: scale(0.625,0.625)

4. 如何解决移动端 Retina 屏 1px 像素问题
5. 如何用 css 或 js 实现多行文本溢出省略效果，考虑兼容性
   单行：
   overflow: hidden;
   text-overflow:ellipsis;
   white-space: nowrap;
   多行：
   display: -webkit-box;
   -webkit-box-orient: vertical;
   -webkit-line-clamp: 3; //行数
   overflow: hidden;
6. 如何实现骨架屏，说说你的思路
   https://juejin.cn/post/6943020826627145735
   核心思想就是：

   1. **puppeteer** 当 Puppeteer 连接到一个 Chromium 实例的时候会通过 puppeteer.launch 或 puppeteer.connect 创建一个 Browser 对象。这个时候你就会获得当前页面的 dom 结构。
   2. 获取你需要做骨架屏的 dom 元素的宽高，你还可以排除一些你不想做骨架屏的元素。
   3. 已知了宽高，你就可以去改她的背景颜色变成一个灰色的方框，看起来就会像一个骨架屏了

7. flex 计算问题

```HTML
<div class="container">
    <div class="left"></div>
    <div class="right"></div>
</div>

<style>
  * {
    padding: 0;
    margin: 0;
  }
  .container {
    width: 600px;
    height: 300px;
    display: flex;
  }
  .left {
    flex: 1 2 500px;
    background: red;
  }
  .right {
    flex: 2 1 400px;
    background: blue;
  }
</style>
```

设：left 的缩小比例是 2x，right 的缩小比例是 x
则：500 \* (1 - 2x) + 400 \* (1 - x) = 600

解得：x = 300 / 1400

left = 500 \* (1 - 2x) = 285.7px
right = 400 \* (1 - x) = 314.3px

8. CSS 会阻塞 DOM 解析吗
   css 加载**不会阻塞 DOM 树的解析**
   css 加载会阻塞 DOM 树的渲染
   css 加载会阻塞后面 js 语句的执行

9. 如何解决 a 标点击后 hover 事件失效的问题?
   改变 a 标签 css 属性的排列顺序
   只需要记住 LoVe HAte 原则就可以了：
   link→visited→hover→active

   ```CSS
      a:hover{
      color: green;
      text-decoration: none;
      }
      a:visited{ /* visited在hover后面，这样的话hover事件就失效了 */
      color: red;
      text-decoration: none;
      }

   ```

   正确的做法是将两个事件的位置调整一下。
   a:link：未访问时的样式，一般省略成 a a:visited：已经访问后的样式 a:hover：鼠标移上去时的样式 a:active：鼠标按下时的样式

10. flex 的兼容性怎样
    IE6~9 不支持，IE10~11 部分支持 flex 的 2012 版，但是需要-ms-前缀。
    其它的主流浏览器包括安卓和 IOS 基本上都支持了。
11. 你知道到哪里查看兼容性吗
    可以到 Can I use 上去查看，官网地址为：https://caniuse.com/
12. 移动端布局总结：
    移动端布局的方式主要使用 rem 和 flex，可以结合各自的优点，比如 flex 布局很灵活，但是字体的大小不好控制，我们可以使用 rem 和媒体查询控制字体的大小，媒体查询视口的大小，然后不同的上视口大小下设置设置 html 的 font-size。
13. rem 和 em 的区别
    em: font-size 时 以父级的字体大小为基准;长度单位时以当前字体大小为基准
    例父级 font-size: 14px，则子级 font-size: 1em;为 font-size: 14px;；若定义长度时，子级的字体大小如果为 14px，则子级 width: 2em;为 width: 28px。
    rem:以根元素的字体大小为基准。例如 html 的 font-size: 14px，则子级 1rem = 14px。
14. 在移动端中怎样初始化根元素的字体大小
    页面开头处引入下面这段代码，用于动态计算 font-size：
    (假设你需要的 1rem = 20px)

```JS
(function () {
  var html = document.documentElement;
  function onWindowResize() {
    html.style.fontSize = html.getBoundingClientRect().width / 20 + 'px';
  }
  window.addEventListener('resize', onWindowResize);
  onWindowResize();
})();

document.documentElement：获取document的根元素
html.getBoundingClientRect().width：获取html的宽度(窗口的宽度)
监听window的resize事件

一般还需要配合一个meta头：
<meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-sacle=1.0, maximum-scale=1.0, user-scalable=no" />

```

15. animation 有一个 steps()功能符知道吗？
    steps()功能符可以让动画不连续。
    和贝塞尔曲线(cubic-bezier()修饰符)一样，都可以作为 animation 的第三个属性值。和贝塞尔曲线的区别：贝塞尔曲线像是滑梯且有 4 个关键字(参数)，而 steps 像是楼梯坡道且只有 number 和 position 两个关键字。
