## 四大布局

1. flex
2. grid
3. 移动端适配布局 rem vw
4. 响应式布局 一套代码 pc 端移动端 @media

此外还有个性布局

大体上 CSS 的学习可以朝着三个方向进行研究，分别是：**特效**、**工程化**、**布局**。比如说特效：可以学习动画、3D、渐变等；工程化：可以学习 sass、less、postcss 等；布局：可以学习 flex、grid、响应式等。

4.  尺寸与位置
    **css 盒模型**
    注意:padding 不能负值 margin 可以负值
    背景会平铺到非 margin 的区域 ，例如边框
    margin 传递问题(**子元素上下方向的 margin 会传递给父级元素**)：解决=>触发 BFC,父元素设置 overflow 不是 visible 如 overflow: hidden
    margin 上下叠加问题:都有 margin 取绝对值较大者(最好解决方案:flex/grid 布局)

    **块级盒子(block)**四个特性:
    占一行，支持所有样式，不写宽度与父同宽，矩形区域
    布局主力
    **内联盒子(inline)**四个特性:
    不占一行，不支持高宽等，不写宽度与内容同宽，不一定矩形区域
    内联盒子间隙怎么消除: 将盒子 font-size 设置为 0
    不用内联盒子布局，而是用来修饰

    **自适应盒模型**
    自适应盒模型是指当盒子不设置宽度时，盒模型相关组成部分的处理方式是如何的
    当我们固定了子元素宽度也就相当于子元素的 content 区域，如果这个时候再给子元素添加 padding、margin 或者 border，盒模型会自动向四周扩散，最终导致溢出父元素。那就有很大的一部分人直接算出左右间距，然后用父盒子的宽度，减去间距然后设置子盒子的宽度。其实，如果了解盒模型的特性，大可不必如此做。
    **利用盒模型特性：不设置子盒子的宽度，直接设置子盒子的 margin**

    **标准盒模型与怪异(IE)盒模型**
    即 contentbox(默认,content) 与 borderbox(content+padding+border)
    又要加百分比又要加 padding contentbox 会溢出 而 borderbox 不会
    应用举例：下部固定的导航标签栏
    **浮动**
    上下排列变左右排列；脱离正常文档流(z 方向层级不同)
    清除浮动 3 种方法

    - clear:both 适用于相邻元素
    - BFC/伪元素::after 3 句话 适用于解决父元素高度塌陷问题
      **定位**

    static 默认
    relative 正常位置+偏移量，**相对自身左上角进行偏移**，加了 relative 也不影响其他元素布局
    absolute 宽度由内容决定
    sticky
    fixed
    **BFC** 即 Block Formatting Contexts(块级格式化上下文),指一个独立的渲染区域或者说是一个隔离的独立容器。
    满足以下条件之一，即可触发 BFC

    - 浮动元素， float 的值不是 none
    - 定位元素， position 的值不是 static 或者 relative
    - **overflow 的值不是 visible** 遇事不决 overflow:hidden
    - display 的值是 inline-block、flex
      &emsp;&emsp;下面的 box 盒子就是一个 BFC 独立容器：

               ```css
               .box {
                 width: 100px;
                 height: 100px;
                 overflow: hidden; /* 触发了BFC，形成独立盒子 */
               }
               ```

      BFC 特性：

      1. 内部的 Box 会在垂直方向上一个接一个的放置；
      2. 垂直方向上的距离由 margin 决定；（解决外边距重叠问题）
      3. bfc 的区域不会与 float 的元素区域重叠；（防止浮动文字环绕）
      4. 计算 bfc 的高度时，浮动元素也参与计算；（清除浮动）
      5. bfc 就是页面上的一个独立容器，容器里面的子元素不会影响外面元素；

      在现代布局 flex 和 grid 中，是默认自带 BFC 规范的，所以可以解决非 BFC 盒子的一些问题，这就是为什么 flex 和 grid 能成为更好的布局方式原因之一。
      BFC 解决 margin 传递与叠加问题：将 div 包裹在`overflow: hidden;`的 section 里
      BFC 清除浮动：父容器加`overflow: hidden;`

    **标签默认样式**
    一些 HTML 标签在浏览器中会有默认样式，例如：body 标签会有 margin:8px；ul 标签会有 margin:16px 0;及 padding-left:40px。
    由于 Reset CSS 相对“暴力”，不管你有没有用，统统重置成一样的效果，且影响的范围很大，所以更加“平和”的一种方式
    **Normalize CSS** 可以看成是一种 Reset CSS 的替代方案。创造 Normalize CSS 有下面这几个目的：

    - 保护有用的浏览器默认样式而不是完全去掉它们
    - 一般化的样式：为大部分 HTML 元素提供
    - 修复浏览器自身的 bug 并保证各浏览器的一致性
    - 优化 CSS 可用性：用一些小技巧
    - 解释代码：用注释和详细的文档来

1.  flex 布局
    轮播图标签:不定项居中布局

```HTML
<style>
      .box {
            width: 300px;
            height: 150px;
            background: skyblue;
            display: flex;
            justify-content: center;
            align-items: flex-end;
        }

        .box div {
            width: 30px;
            height: 30px;
            background: pink;
            border-radius: 50%;
            margin:5px;
        }
</style>
```

下部固定的导航标签栏：均分列布局

```HTML
<style>
        .main{
            height:200px;
            background:skyblue;
            display: flex;
            justify-content: space-between;
            align-items: flex-end;
            padding:0 20px;
        }
        .main div{
            width:30px;
            height:30px;
            background:pink;
            border-radius: 50%;
        }
</style>
```

头部导航栏：子项分组布局(有左有右)

```HTML

<style>
  .main {
        height: 200px;
        background: skyblue;
        display: flex;
        align-items: center;
      }
      .main div {
        width: 50px;
        height: 100px;
        background: pink;
        margin-right: 10px;
      }
      /* 右边距自适应分组，右边全部挤到角落 */
      .main div:nth-of-type(3) {
        margin-right: auto;
      }
      .main div:nth-of-type(6) {
        margin-right: auto;
      }
</style>
```

flex 子项 默认值
flex:flex-grow flex-shrink flex-basis(主轴上的);
0 1 auto(默认值为 auto，即项目的本来大小)

flex:1 表示 flex:1 1 0%

**flex-grow:1 表示占满剩余的所有空间**
**flex-shrink : 1 表示自动收缩，跟容器大小相同**
flex-basis 属性定义了在分配剩余空间之前，项目占据的主轴空间（main size）

多行 item 栏:(左右)等高布局

```HTML
<style>
      .main {
        width: 500px;
        background: skyblue;
        display: flex;
        justify-content: space-between;
      }
      .main div {
        width: 100px;
        background: pink;
      }
</style>
```

后台管理界面:两列与三列布局(定宽+自适应)
flex-grow: 1;

```HTML
<style>
      .main {
        height: 100vh;
        background: skyblue;
        display: flex;
      }
      .col1 {
        width: 200px;
        background: pink;
      }
      .col2 {
        flex-grow: 1;
        background: springgreen;
      }
      .col3 {
        width: 100px;
        background: tomato;
      }
</style>
```

管理系统粘性页脚适配:stickyFooter 布局
内容不满一屏 footer 也在最底端
flex-grow: 1;

```HTML

<style>
        .main {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
      }
      .main .header {
        height: 100px;
        background: pink;
      }
      /* 适配剩余空间 */
      .main .content {
        flex-grow: 1;
      }
      .main .footer {
        height: 100px;
        background: skyblue;
      }
</style>
```

移动端溢出隐藏菜单:溢出项布局
flex-shrink: 0;

```HTML

<style>
      .main {
        height: 100px;
        background: skyblue;
        display: flex;
        align-items: center;
      }
      .main div {
        width: 100px;
        height: 80px;
        background: pink;
        margin-right: 10px;
        /* 不收缩 */
        flex-shrink: 0;
      }
</style>
```

综合案例:
Swiper 轮播图
知乎导航

6. grid 网格布局
   **grid 容器**

**定义网格与** fr(fractional unit) 单位
grid-template-columns
grid-template-rows

```HTML
 <style>
      .container {
        display: grid;
        width: 300px;
        height: 300px;
        background-color: rgb(36, 73, 128);
        grid-template-columns: 50px 1fr 1fr;
        /* fr 平分 */
        grid-template-rows: 2fr 1fr;
      }

      .container > div {
        background-color: rgb(235, 227, 228);
      }
    </style>
```

**合并网格**
grid-template-areas 与 grid-area：注册与使用
grid-area:行起点/列起点/行长/列长
grid-area: 1/1 / span 1 / span 3;

```HTML
<style>
      .container {
        display: grid;
        width: 300px;
        height: 300px;
        background-color: rgb(36, 73, 128);
        grid-template:
          'a1 a1 a2'
          'a1 a1 a2'
          'a3 a3 a3';
      }

      .container > div {
        background-color: rgb(235, 227, 228);
      }

      .container > div:nth-of-type(1) {
        grid-area: a1;
      }

      .container > div:nth-of-type(2) {
        grid-area: a2;
      }

      .container > div:nth-of-type(3) {
        grid-area: a3;
      }
    </style>
```

grid-template:是 grid-template-columns+grid-template-rows+
grid-template-areas 合并的写法

```HTML
     <style>
      .container {
        display: grid;
        width: 300px;
        height: 300px;
        background-color: rgb(36, 73, 128);
        grid-template:
          'a1 a1 a2' 1fr
          'a1 a1 a2' 1fr
          'a3 a3 a3' 1fr
          / 1fr 1fr 1fr;
      }

      .container > div {
        background-color: rgb(235, 227, 228);
      }

      .container > div:nth-of-type(1) {
        grid-area: a1;
      }

      .container > div:nth-of-type(2) {
        grid-area: a2;
      }

      .container > div:nth-of-type(3) {
        grid-area: a3;
      }
    </style>
```

**网格间隙**
row-gap column-gap gap (弹性布局也能用)

```HTML
   <style>
      .container {
        display: grid;
        width: 300px;
        height: 300px;
        background-color: rgb(36, 73, 128);
        grid-template-areas:
          'a1 a1 a2'
          'a1 a1 a2'
          'a3 a3 a3';
        /* row-gap: 10px;
        column-gap: 5px; */
        gap: 10px 5px;
      }

      .container > div {
        background-color: rgb(235, 227, 228);
      }

      .container > div:nth-of-type(1) {
        grid-area: a1;
      }

      .container > div:nth-of-type(2) {
        grid-area: a2;
      }

      .container > div:nth-of-type(3) {
        grid-area: a3;
      }
    </style>
```

**对齐方式**
子项在网格中的对齐方式(子项小于 grid)
justify-items: start;
align-items: center;
place-items: start center;
**默认靠左靠上**
网格在容器中的对齐方式 (grid 小于容器)
justify-content: start;
align-content: center;
place-content: center start;
**默认靠左靠上**

```HTML
<style>
      .main {
        width: 300px;
        height: 300px;
        background: skyblue;
        display: grid;
        grid-template-columns: 100px 100px 100px;
        grid-template-rows: 100px 100px 100px;
        justify-items: start;
        align-items: center;
        /* place-items: start center; */
      }

      .main div {
        width: 50px;
        height: 50px;
        background: pink;
      }

      .main2 {
        width: 500px;
        height: 500px;
        background: skyblue;
        display: grid;
        grid-template-columns: 100px 100px 100px;
        grid-template-rows: 100px 100px 100px;
        justify-content: start;
        align-content: center;
        /* place-content: center start; */
      }

      .main2 div {
        width: 50px;
        height: 50px;
        background: pink;
      }
    </style>
```

**隐式显式网格可以做到自适应效果**
/_ 行产生隐式网格 _/
`默认 row`
grid-auto-flow: row;
/_ 隐式网格高度 _/
grid-auto-rows: 100px;

grid-auto-flow: column;
grid-auto-columns: 100px;

```HTML
    <style>
      .container {
        display: grid;
        width: 300px;
        height: 300px;
        background-color: rgb(36, 73, 128);
        grid-template-columns: 100px 100px 100px;
        grid-template-rows: 100px;
        /* 行产生隐式网格 */
        grid-auto-flow: row;
        /* 隐式网格高度 */
        grid-auto-rows: 150px;

        /* column 就是列产生隐式网格 */
        grid-auto-flow: column;
        /* 可以调节产生隐式网格的宽度 */
        grid-auto-columns: 100px;
      }

      .container > div {
        background-color: rgb(235, 227, 228);
      }
    </style>
```

**grid 子项属性**
**控制子项起始占的格子**
grid-column-start: 2;
grid-column-end: 3;
grid-column: 2 / 3

grid-row-start: 1;
grid-row-end: 2;
grid-row: 1 / 2

```HTML
 /* grid-column-start: 2;
        grid-column-end: 3; */
        /* 默认值：auto */
        /* grid-row-start: 1;
            grid-row-end: 2; */
        grid-column-start: 2;
```

**子项对齐方式**
justify-self
align-self
place-self

**repeat 与 minmax**
**repeat(个数,尺寸)**
grid-template-columns: 150px 100px 100px;
grid-template-columns: 150px repeat(2, 100px);

auto-fill 表示个数根据父容器自动决定
grid-template-columns: repeat(**auto-fill**, 100px);

**minmax: 最大最小范围**
两边固定 中间自适应大小
grid-template-columns: 100px **minmax(100px, 1fr)** 100px;

**根据分辨率自适应元素个数**
每行多少个?每个多大？

```HTML
<style>
      .main {
        background: skyblue;
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
        grid-template-rows: 100px;
        grid-auto-rows: 100px;
        grid-gap: 20px 20px;
      }

      .main div {
        background: pink;
        border: 1px black solid;
      }
    </style>
```

**比定位更方便的叠加布局**
实现思路
1:产生的隐式网格全部定义为相同的 grid-area
2。 调节子项 justify-self 与 align-self

```HTML
<style>
      .main {
        width: 530px;
        height: 300px;
        background: skyblue;
        display: grid;
      }

      .main img {
        grid-area: 1/1/1/1;
      }

       .main span {
        grid-area: 1/1/1/1;
        justify-self: end;
        align-self: end;
        margin: 5px;
      }

      .main p {
        grid-area: 1/1/1/1;
        align-self: center;
        margin: 0;
        padding: 0;
        background: rgba(0, 0, 0, 0.5);
        height: 30px;
        line-height: 30px;
        color: white;
      }
    </style>

  <body>
    <div class="main">
      <img src="./phone.png" alt="" />
      <!-- 隐式网格 -->
      <span>自制</span>
      <!-- 隐式网格 -->
      <p>手机热卖中.....</p>
    </div>
  </body>
```

**多种排列组合布局**
实现思路
1:划分最小网格
2:调节 grid-area (左上角是 1/1)

```HTML
<style>
      .main {
        width: 300px;
        height: 300px;
        background: skyblue;
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        grid-template-rows: repeat(3, 1fr);
        gap: 5px;
      }
      .main div {
        background: pink;
      }
      .main div:nth-of-type(1) {
        /* grid-area: 1/1/span 2/span 2; */
        grid-area: 2/1 / span 2 / span 2;
      }
    </style>
```

**栅格布局**
网页一般分成 12/24 个栅格
实现思路

1.划分列 grid-template-columns: repeat(12, 1fr);

2.隐式网格配置
grid-template-rows: 50px;
grid-auto-rows: 50px;

3.设置类名占的区域
grid-area: auto/auto/auto/span 1;

```HTML
<style>
        .row{
            background:skyblue;
            display: grid;
            grid-template-columns: repeat(12, 1fr);
            grid-template-rows: 50px;
            grid-auto-rows: 50px;
        }
        .row div{
            background:pink;
            border:1px black solid;
        }
        .row .col-1{
            grid-area: auto/auto/auto/span 1;
        }
        .row .col-2{
            grid-area: auto/auto/auto/span 2;
        }
        .row .col-3{
            grid-area: auto/auto/auto/span 3;
        }
        .row .col-4{
            grid-area: auto/auto/auto/span 4;
        }
        .row .col-5{
            grid-area: auto/auto/auto/span 5;
        }
        .row .col-6{
            grid-area: auto/auto/auto/span 6;
        }
        .row .col-7{
            grid-area: auto/auto/auto/span 7;
        }
        .row .col-8{
            grid-area: auto/auto/auto/span 8;
        }
        .row .col-9{
            grid-area: auto/auto/auto/span 9;
        }
        .row .col-10{
            grid-area: auto/auto/auto/span 10;
        }
        .row .col-11{
            grid-area: auto/auto/auto/span 11;
        }
        .row .col-12{
            grid-area: auto/auto/auto/span 12;
        }
    </style>

    <body>
      <div class="row">
          <div class="col-6">1</div>
          <div class="col-3">2</div>
          <div class="col-4">3</div>
          <div class="col-5">4</div>
      </div>
    </body>
```

**容器自适应行列布局**

**行自适应**

1. 设置列数
2. 设置隐式网格高度(自动)

```HTML
 <style>
      .main {
        display: grid;
        width: 300px;
        /* height: 300px; */
        background-color: pink;
        grid-template-columns: repeat(3, 1fr);
        grid-auto-rows: 100px;
        gap: 5px;
      }

      .main > div {
        background-color: rgb(44, 45, 104);
      }
    </style>
```

**百度热词风云榜**
划分网格加 grid-area
**小米商品导航**
列自适应

7. 移动端布局
   关于 px

   - **逻辑像素：CSS 中的像素，绝对单位，保证不同设备下元素的尺寸是相同的**。
   - **物理像素：设备屏幕实际拥有的像素点，相对单位，不同设备下物理像素大小不同**。

   1. **rem 布局**
      rem + 动态计算 font-size
      发展过程:flexible.js
      利用 vw 动态换算

   ```HTML
    <style>
    * {
      margin: 0;
      padding: 0;
    }
    html {
      /* font-size: 100vw; */ /* 在iphone6 -> 375px */
      font-size: 26.666667vw; /* 在iphone6 -> 100px */
    }
    body {
      font-size: 0.16rem; /* rem布局一定要在body重置font-size大小 */
    }
    .box {
      width: 1rem; /* 页面可视区分成了100vw和100vh */
      height: 1rem;
      background: pink;
    }
   </style>
   ```

   - px to rem 插件 alt + z 批量转换
   - 蓝湖/**pxCook** 量取 rem 注意 2 倍

     点击更多效果：

     ```HTML
     <style>
        .nav-sub__closed{
            height: 0.7rem;
            overflow: hidden;
        }
        .nav-sub__closed + .nav-sub-arrow{
            transform: rotate(0);
        }
     </style>

     ```

   2. **移动端 vw 布局及插件使用**
      px -to - vw 插件

8. 响应式布局
   媒体查询

   1. 媒体类型 screen / print
   2. 媒体特性 width / aspect-ratio/orientation(lanscape,protrait)
   3. 逻辑操作符 and/not/only/,
   4. link 标签方式

   ```HTML
   <link rel="stylesheet" href="./a.css" media="(orientation: portrait)" />
   <link rel="stylesheet" href="./b.css" media="(orientation: landscape)" />

   ```

   **媒体查询编写位置及顺序**

   - 样式表底部，对 css 进行优先级的覆盖
   - 移动端到 PC 端: **min-width 从小到大** 576 768 992 1200 1400

   **响应式栅格系统**
   栅格布局+断点设定
   外层加.row 包裹
   内层加 .col-sm-4/.col-md-4 分格数
   **响应式交互实现**
   例如元素显示隐藏
   使用`状态`

   - :checked 伪类

```HTML
 <style>
    ul {
      display: none;
    }

    input {
      display: none;
    }
    /* +找相邻ul */
    input:checked + ul {
      display: block;
    }

    @media (min-width: 700px) {
      ul {
        display: block;
      }
      span {
        display: none;
      }
    }
  </style>
   <body>
    <label for="menu">
      <span> 菜单按钮 </span>
    </label>
    <input id="menu" type="checkbox" />
    <ul>
      <li>首页</li>
      <li>教程</li>
      <li>论坛</li>
      <li>文章</li>
    </ul>
  </body>
```

- ghost 博客记录

1. 点击显示隐藏怎么做:
   :checked 伪类
2. nav 竖变横怎么做：
   大于 768px display 切换为 `flex !important`
3. 版芯如何响应式(文章最大宽度)
   使用

   ```HTML
   <style>
    :root {
       --container: 100%;
     }


    .版芯 {
      padding: 0 15px;
      width: var(--container);
      margin: 0 auto;
      box-sizing: border-box;
    }
   </style>
   ```

   然后在@media 里改变--container

9.antd 综合布局
BEM
BEM 中的块、元素和修饰符需要全部小写，名称中的单词用连字符（-）分隔，元素由双下划线（\_\_）分隔，修饰符由双连字符（--）分隔。注意，块和元素都既不能是 HTML 元素名或 ID，也不依赖其它块或元素。

```css
.setting-menu {
}
.setting-menu--open {
}
.setting-menu__head {
}
.setting-menu__head--fixed {
}
.setting-menu__content {
}
```

        笔记

        1. 根据状态固定左侧栏

        ```CSS
        .g-ant-sider__wrap--fixed{
          position: fixed;
          /* left:0;
          top:0;
          width: inherit; */
          overflow: hidden;
        }
        ```

        2. 根据状态折叠左侧栏

        ```CSS
        .g-ant__sider--closed{
            width: 48px;
            overflow: hidden;
        }

        ```

        3. 指定元素上添加自定义滚动条

        ```CSS
        .g-ant-sider__main{
            flex-grow: 1;
            overflow: hidden auto;
        }

        .g-ant-sider__main::-webkit-scrollbar {
            width : 6px;
            height: 6px;
        }
        .g-ant-sider__main::-webkit-scrollbar-thumb {
            background   : #51606d;
            border-radius: 3px;
        }
        .g-ant-sider__main::-webkit-scrollbar-track {
            background   : #263849;
            border-radius: 3px;
        }
        ```

        4. 主体网格布局
        5. 主题色切换
          ```JS
            document.documentElement.style.setProperty('--theme', color);
          ```
        6. 右侧 panel
          关闭时 right:-300px
          打开时:right:0

10.瀑布流
grid 不推荐(网格必须要规矩)
**横向瀑布流**
推荐 flex

```HTML
   <style>
      .main {
        display: flex;
        flex-wrap: wrap;
        gap: 10px;
      }
      .main .item {
        flex-grow: 1;
        height: 200px;
        /* flex-basis: 200px; */
      }
      .main .item img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        /* display: block; */
      }
    </style>

     <script>
      var items = document.querySelectorAll('.item')
      for (var i = 0; i < items.length; i++) {
        items[i].style.flexBasis = Math.random() * 200 + 200 + 'px'
      }
    </script>
```

**竖向瀑布流(Lofter)**
不推荐 flex (flex-wrap 必须要设定高度)
**多列方式** 需要 JS 配合实现
**Multi-Columns 多栏布局** 纯 CSS

```CSS
        /* 多兰个数 */
        column-count: 4;
        column-width: 300px;
        column-gap: 20px;
        /* 多兰分割线 */
        column-rule: 1px gray dashed;
```

多栏布局缺陷：**无法实现加载更多**(新的会到上面去)

11.视差

```CSS
/* CSS */
background-attachment:fixed
```

监听滚动条插件:skrollr
控制 opacity top font-size transform 等属性
