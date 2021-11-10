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
16. 如何让<p>测试 空格</p>这两个词之间的空格变大
    通过给 p 标签设置 word-spacing，将这个属性设置成自己想要的值。
    将这个空格用一个 span 标签包裹起来，然后设置 span 标签的 letter-spacing 或者 word-spacing。
    **letter-spacing 添加字母(letter)之间的空白，而 word-spacing 添加每个单词(word)之间的空白。**
    letter-spacing 把中文之间的间隙也放大了，而 word-spacing 则不放大中文之间的间隙。
17. 如何解决 inline-block 空白问题？
    给父级设置 font-size: 0

```HTML
<style>
.sub {
  background: hotpink;
  display: inline-block;
  /* 给父级设置font-size: 0 */
}
</style>
<body>
  <div class="super">
    <div class="sub">
      孩子
    </div>
    <div class="sub">
      孩子
    </div>
    <div class="sub">
      孩子
    </div>
  </div>
</body>

```

18. 脱离文档流是不是指该元素从 DOM 树中脱离
    并不会，DOM 树是 HTML 页面的层级结构，指的是元素与元素之间的关系，例如包裹我的是我的父级，与我并列的是我的兄弟级，类似这样的关系称之为层级结构。
    而文档流则类似于排队，我本应该在队伍中的，然而我脱离了队伍，但是我与我的父亲，兄弟，儿子的关系还在
19. 如何让 Chrome 浏览器支持小于 12px 的字体大小
    使用：-webkit-transform: scale(0.8);
20. 空元素(单标签)元素有哪些？
<br />
<hr />
<input />
<img />
<link />
<meta>

21. b 与 strong 的区别以及 i 和 em 的区别？
    被<b>和<strong>包裹的文字会被加粗
    被<i>和<em>包裹的文字会以斜体的方式呈现
    <b>标签和<i>标签都是「自然样式标签」，都只是在样式上加粗和变斜，并没有什么实际的意义。并且据了解，这两种标签在 HTML4.01 中已经不被推荐使用了。<strong>标签和<em>的话是「语义样式标签」。就像是<h1>、<h2>一样都有自己的语义。<em>表示一般的强调文本，而<strong>表示更强的强调文本。另外在使用阅读设备的时候，<strong>会重读(这点呆呆也没有实践过所以不太敢保证)。
22. http-equiv
    equivalent(元信息元素的附加属性)
    content-security-policy
    它允许页面作者定义当前页的内容策略。 内容策略主要指定允许的服务器源和脚本端点，这有助于防止跨站点脚本攻击。
    refresh
    这个属性指定:
    如果 content 只包含一个正整数，则为重新载入页面的时间间隔(秒)；
    如果 content 包含一个正整数，并且后面跟着字符串 ';url=' 和一个合法的 URL，则是重定向到指定链接的时间间隔(秒)

23. 获取元素的属性
    **注:** 以下四个属性只能读取 不能对元素进行修改;
    1.offsetWidth 获取元素的实际宽度 包含 border 和 padding 在内
    2.offsetHeight 获取元素的实际高度 包含 border 和 padding 在内
    3.offsetLeft 元素定位之后相对于参照物父容器的偏移
    4.offsetTop 元素定位之后相对于参照物父容器的偏移
24. 选择器 https://github.com/LinDaiDai/niubility-coding-js/blob/master/CSS/CSS%E7%9A%84%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86.md
    子代选择器 ： 只选择直系的后代 .class_1 > p { }
    **同级元素通用选择器** 用法：选择器 1~选择器 2{ }
    俩个选择器之间需要有相同的父级
    选择器 2 必须处于选择器 1 的后面
    选择具有相同的父级，并且加载顺序处于后面的内容
    **相邻兄弟选择器** 用法： 选择器 1+选择器 2{ }
    俩个选择器必须是兄弟关系(也就是要有同一个父级)
    俩个选择器必须是紧挨着的
    选择的是相连接的后面的兄弟
    **可以看到选择器 +和~的区别就是+只针对一项元素,而~可能是多项的。**
25. 如何去除 CSS 注释？
    用正则/\/_[\s\S]_?\*\//全局匹配注释，替换为空字符串
    用工程化工具，如 cssnano 来去除注释
26. 伪类 伪元素
    伪类选择器 a:hover
    伪元素选择器 a::before
27. 什么是继承
    CSS 属性分为非继承属性和 继承属性

    常见的继承属性：
    字体 font 系列
    文本 text-align text-ident line-height letter-spacing
    颜色 color
    列表 list-style
    可见性 visibility
    光标 cursor

    容易被误认为继承属性的 非继承属性：
    透明度 opacity
    背景 background 系列

28. 何计算 CSS 选择器的优先级？
    |选择器|权重|
    -----|-----|-----
    |style=""| 1000|
    |ID | 100|
    |类、伪类、属性| 10|
    |元素、伪元素| 1|
    |关系、通配符| 0|
29. 百分比 % 相对于谁？
    百分比总是相对于父元素，无论是设置 font-size 或 width 等。如果父元素的相应属性，经浏览器计算后，仍无绝对值，那么 % 的实际效果等同于 默认值，如 height: 100%
30. 对比块、内联和内联块盒子
    块盒子：display:block
    换行
    width 和 height 生效
    竖直方向 padding 和 margin 生效

    内联盒子：display:inline
    不换行
    width 和 height **无效**
    竖直方向 padding 和 margin **无效**

    内联块盒子：display:inline-block
    **不换行**
    width 和 height 生效
    竖直方向 padding 和 margin 生效

31. 什么是层叠上下文
    层叠上下文是元素在 Z 轴上的层次关系集合并影响渲染顺序，设置 z-index 可改变 position 不为 static 的元素的层叠顺序
    z-index，常用来：

    **改善兼容性**
    解决遮挡问题
    解决滚动穿透问题

    **提升移动端体验**
    如通过-webkit-overflow-scrolling: touch 增加滚动回弹效果

    **性能优化**
    将频繁变化的内容单独一层放置

32. BEM
    块元素修饰符，全称是 Block Element Modifier
    定义了一种 CSS 的命名规则，用于解决命名冲突：
    `.block-name__element-name--modifier-name`
    .块名\_\_元素名--修饰符（元素名和修饰符都可为空）
    其中：

    Block：块，忽略结构和优先级，具有独立意义的实体
    Element：元素，块内部没有独立意义的实体
    Modifier：修饰符，标识块或元素的外观、行为、状态被修改
    含有修饰符的类名不可独立出现，通常跟在不含修饰符的类名后

33. css 的渲染层合成是什么 浏览器如何创建新的渲染层
    在 DOM 树中每个节点都会对应一个渲染对象（RenderObject），当它们的**渲染对象处于相同的坐标空间（z 轴空间）时，就会形成一个 RenderLayers，也就是渲染层**。渲染层将保证页面元素以正确的顺序堆叠，这时候就会出现层合成（composite），从而正确处理透明元素和重叠元素的显示。对于有位置重叠的元素的页面，这个过程尤其重要，因为一旦图层的合并顺序出错，将会导致元素显示异常。

    - 浏览器如何创建新的渲染层
      overflow 不为 visible
      有 CSS transform 属性且值不为 none
      有明确的定位属性（relative、fixed、sticky、absolute）
      opacity < 1

34. css 优先级是怎么计算的
    1.  第一优先级：!important 会覆盖页面内任何位置的元素样式
    2.  内联样式，如 style="color: green"，权值为 1000
    3.  ID 选择器，如#app，权值为 0100
    4.  类、伪类、属性选择器，如.foo, :first-child, div[class="foo"]，权值为 0010
    5.  标签、伪元素选择器，如 div::first-line，权值为 0001
    6.  通配符、子类选择器、兄弟选择器，如\*, >, +，权值为 0000
    7.  继承的样式没有权值
35. css 怎么开启硬件加速(GPU 加速)
    浏览器在处理下面的 css 的时候，会使用 GPU 渲染

- transform（当 3D 变换的样式出现时会使用 GPU 加速）
- opacity
- filter
- will-change
  采用 transform: translateZ(0)
  采用 transform: translate3d(0, 0, 0)
  使用 CSS 的 will-change 属性。 will-change 可以设置为 opacity、transform、top、left、bottom、right。

36. 透明度 opacity 和 rgba 的区别
    最大的不同是 opacity 作用于元素，以及元素内的所有内容的透明度，而 rgba()只作用于元素的颜色或其背景色。
37. background-color 属性可以覆盖 background-image 属性吗？
    当元素本身设置了 background-image 属性时，**如果设置了 background-color，图片不会被覆盖，**background-color 会在 image 底层；如果设置的是 background，那么图片会被颜色给覆盖掉。
