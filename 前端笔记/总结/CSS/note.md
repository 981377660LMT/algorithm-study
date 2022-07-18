1. flex 计算问题

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
   只需要记住 LoVe HAte (LVHA) 原则就可以了：
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

10. 你知道到哪里查看兼容性吗
    可以到 Can I use 上去查看，官网地址为：https://caniuse.com/
11. 移动端布局总结：
    移动端布局的方式主要使用 rem 和 flex，可以结合各自的优点，比如 flex 布局很灵活，但是字体的大小不好控制，我们可以使用 rem 和媒体查询控制字体的大小，媒体查询视口的大小，然后不同的上视口大小下设置设置 html 的 font-size。
12. rem 和 em 的区别
    em: font-size 时 以父级的字体大小为基准;长度单位时以当前字体大小为基准
    例父级 font-size: 14px，则子级 font-size: 1em;为 font-size: 14px;；若定义长度时，子级的字体大小如果为 14px，则子级 width: 2em;为 width: 28px。
    rem:以根元素的字体大小为基准。例如 html 的 font-size: 14px，则子级 1rem = 14px。

13. animation 有一个 steps()功能符知道吗？
    steps()功能符可以让动画不连续。
    和贝塞尔曲线(cubic-bezier()修饰符)一样，都可以作为 animation 的第三个属性值。和贝塞尔曲线的区别：贝塞尔曲线像是滑梯且有 4 个关键字(参数)，而 steps 像是楼梯坡道且只有 number 和 position 两个关键字。
14. 如何让<p>测试 空格</p>这两个词之间的空格变大
    通过给 p 标签设置 word-spacing，将这个属性设置成自己想要的值。
    将这个空格用一个 span 标签包裹起来，然后设置 span 标签的 letter-spacing 或者 word-spacing。
    **letter-spacing 添加字母(letter)之间的空白，而 word-spacing 添加每个单词(word)之间的空白。**
    letter-spacing 把中文之间的间隙也放大了，而 word-spacing 则不放大中文之间的间隙。

15. 脱离文档流是不是指该元素从 DOM 树中脱离
    并不会，DOM 树是 HTML 页面的层级结构，指的是元素与元素之间的关系，例如包裹我的是我的父级，与我并列的是我的兄弟级，类似这样的关系称之为层级结构。
    而文档流则类似于排队，我本应该在队伍中的，然而我脱离了队伍，但是我与我的父亲，兄弟，儿子的关系还在
16. 如何让 Chrome 浏览器支持小于 12px 的字体大小
    使用：-webkit-transform: scale(0.8);

17. 百分比 % 相对于谁？
    `百分比总是相对于父元素`，无论是设置 font-size 或 width 等。如果父元素的相应属性，经浏览器计算后，仍无绝对值，那么 % 的实际效果等同于 默认值，如 height: 100%
18. 对比块、内联和内联块盒子
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

19. 什么是层叠上下文
    层叠上下文是元素在 Z 轴上的层次关系集合并影响渲染顺序，设置 z-index 可改变 position 不为 static 的元素的层叠顺序
    z-index，常用来：

    **改善兼容性**
    解决遮挡问题
    解决滚动穿透问题

    **提升移动端体验**
    如通过-webkit-overflow-scrolling: touch 增加滚动回弹效果

    **性能优化**
    将频繁变化的内容单独一层放置

20. BEM
    块元素修饰符，全称是 Block Element Modifier
    定义了一种 CSS 的命名规则，用于解决命名冲突：
    `.block-name__element-name--modifier-name`
    .块名\_\_元素名--修饰符（元素名和修饰符都可为空）
    其中：

    Block：块，忽略结构和优先级，具有独立意义的实体
    Element：元素，块内部没有独立意义的实体
    Modifier：修饰符，标识块或元素的外观、行为、状态被修改
    含有修饰符的类名不可独立出现，通常跟在不含修饰符的类名后

21. css 的渲染层合成是什么 浏览器如何创建新的渲染层
    在 DOM 树中每个节点都会对应一个渲染对象（RenderObject），当它们的**渲染对象处于相同的坐标空间（z 轴空间）时，就会形成一个 RenderLayers，也就是渲染层**。渲染层将保证页面元素以正确的顺序堆叠，这时候就会出现层合成（composite），从而正确处理透明元素和重叠元素的显示。对于有位置重叠的元素的页面，这个过程尤其重要，因为一旦图层的合并顺序出错，将会导致元素显示异常。

    - 浏览器如何创建新的渲染层
      overflow 不为 visible
      有 CSS transform 属性且值不为 none
      有明确的定位属性（relative、fixed、sticky、absolute）
      opacity < 1

22. css 怎么开启硬件加速(GPU 加速)
    浏览器在处理下面的 css 的时候，会使用 GPU 渲染

- transform（当 3D 变换的样式出现时会使用 GPU 加速）
- opacity
- filter
- will-change
  采用 transform: translateZ(0)
  采用 transform: translate3d(0, 0, 0)
  使用 CSS 的 will-change 属性。 will-change 可以设置为 opacity、transform、top、left、bottom、right。

23. z-index 有什么需要注意的地方

    - **默认值，auto。 不创建新的局部层叠上下文**。 也就是说，该元素的层叠水平（真正的层叠顺序）和其父元素是一样的。
    - 整数。 可以是负数，0，正数。 这个会创建层叠上下文。
    - `只作用与定位元素`，非定位元素无效
