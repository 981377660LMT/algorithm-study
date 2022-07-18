#### 1.介绍一下标准的 CSS 的盒子模型？低版本 IE 的盒子模型有什么不同的？

回答：

盒模型都是由四个部分组成的，分别是 margin、border、padding 和 content。

标准盒模型和 IE 盒模型的区别在于设置 width 和 height 时，所对应的范围不同。
标准盒模型(content-box)的 width 和 height 属性的范围只包含了 content，
而 IE 盒模型(border-box)的 width 和 height 属性的范围包含了 border、padding 和 content。

一般来说，我们可以通过修改元素的 box-sizing 属性来改变元素的盒模型。

详细的资料可以参考：
[《CSS 盒模型详解》](https://juejin.im/post/59ef72f5f265da4320026f76)

#### 2.CSS 选择符有哪些？

（1）id 选择器（#myid）
（2）类选择器（.myclassname）
（3）标签选择器（div,h1,p）
（4）后代选择器（h1 p）
（5）相邻后代选择器（子）选择器（ul>li）
（6）兄弟选择器（li~a）
（7）相邻兄弟选择器（li+a）
（8）属性选择器（a[rel="external"]）
（9）伪类选择器（a:hover,li:nth-child）
（10）伪元素选择器（::before、::after）
（11）通配符选择器（\*）

#### 3.::before 和:after 中双冒号和单冒号有什么区别？解释一下这 2 个伪元素的作用。

相关知识点：

单冒号（:）用于 CSS3 伪类，双冒号（::）用于 CSS3 伪元素。（伪元素由双冒号和伪元素名称组成）
双冒号是在当前规范中引入的，用于区分伪类和伪元素。`不过浏览器需要同时支持旧的已经存在的伪元素写法，比如:first-line、:first-letter、:before、:after 等， 而新的在 CSS3 中引入的伪元素则不允许再支持旧的单冒号的写法。`

回答：

在 css3 中使用单冒号来表示伪类，用双冒号来表示伪元素。但是为了兼容已有的伪元素的写法，在一些浏览器中也可以使用单冒号来表示伪元素。

伪类一般匹配的是元素的一些特殊状态，如 hover、link 等，而伪元素一般匹配的特殊的位置，比如 after、before 等。

#### 4.伪类与伪元素的区别

css 引入伪类和伪元素概念是为了`格式化文档树以外的信息`。也就是说，伪类和伪元素是用来修饰不在文档树中的部分，比如，一句话中的第一个字母，或者是列表中的第一个元素。

伪类用于当已有的元素处于某个状态时，为其添加对应的样式，这个状态是根据用户行为而动态变化的。比如说，当用户悬停在指定的元素时，我们可以通过:hover 来描述这个元素的状态。

伪元素用于创建一些不在文档树中的元素，并为其添加样式。它们允许我们为元素的某些部分设置样式。比如说，我们可以通过::before 来在一个元素前增加一些文本，并为这些文本添加样式。虽然用户可以看到这些文本，但是这些文本实际上`不在文档树中。`

详细资料可以参考：
[《总结伪类与伪元素》](http://www.alloyteam.com/2016/05/summary-of-pseudo-classes-and-pseudo-elements/)

#### 6.CSS 优先级算法如何计算？

相关知识点：

CSS 的优先级是根据样式声明的特殊性值来判断的。

选择器的特殊性值分为四个等级，如下：

（1）标签内选择符 x,0,0,0
（2）ID 选择符 0,x,0,0
（3）class 选择符/属性选择符/伪类选择符 0,0,x,0
（4）元素和伪元素选择符 0,0,0,x
(5) 通配符
(6) 继承样式

计算方法：

（1）每个等级的初始值为 0
（2）每个等级的叠加为选择器出现的次数相加
（3）`不可进位`，比如 0,99,99,99
（4）依次表示为：0,0,0,0
（5）每个等级计数之间没关联
（6）等级判断从左向右，如果某一位数值相同，则判断下一位数值
（7）`如果两个优先级相同，则最后出现的优先级高，!important 也适用`
（8）通配符选择器的特殊性值为：0,0,0,0
（9）继承样式优先级最低，通配符样式优先级高于继承样式
（10）!important（权重），它没有特殊性值，但它的优先级是最高的，为了方便记忆，可以认为它的特殊性值为 1,0,0,0,0。

计算实例：

（1）#demo a{color: orange;}/_特殊性值：0,1,0,1_/
（2）div#demo a{color: red;}/_特殊性值：0,1,0,2_/

注意：
（1）样式应用时，css 会先查看规则的权重（!important），加了权重的优先级最高，当权重相同的时候，会比较规则的特殊性。

（2）特殊性值越大的声明优先级越高。

（3）相同特殊性值的声明，根据样式引入的顺序，后声明的规则优先级高（距离元素出现最近的）

(4) 部分浏览器由于字节溢出问题出现的进位表现不做考虑

回答：

判断优先级时，首先我们会判断一条属性声明是否有权重，也就是是否在声明后面加上了`!important`。一条声明如果加上了权重，那么它的优先级就是最高的，前提是它之后不再出现相同权重的声明。如果权重相同，我们则需要去比较匹配规则的特殊性。

一条匹配规则一般由多个选择器组成，一条规则的特殊性由组成它的选择器的特殊性累加而成。选择器的特殊性可以分为四个等级，
第一个等级是`行内样式`，为 1000，第二个等级是 `id 选择器`，为 0100，第三个等级是`类选择器、伪类选择器和属性选择器`，为 0010，
第四个等级是`元素选择器和伪元素`选择器，为 0001。规则中每出现一个选择器，就将它的特殊性进行叠加，这个叠加只限于对应的等级的叠加，不会产生进位。选择器特殊性值的比较是从左向右排序的，也就是说以 1 开头的特殊性值比所有以 0 开头的特殊性值要大。
比如说特殊性值为 1000 的的规则优先级就要比特殊性值为 0999 的规则高。如果两个规则的特殊性值相等的时候，那么就会根据它们`引入的顺序`，后出现的规则的优先级最高。

对于组合声明的特殊性值计算可以参考：
[《CSS 优先级计算及应用》](https://www.jianshu.com/p/1c4e639ff7d5)
[《CSS 优先级计算规则》](http://www.cnblogs.com/wangmeijian/p/4207433.html)
[《有趣：256 个 class 选择器可以干掉 1 个 id 选择器》](https://www.zhangxinxu.com/wordpress/2012/08/256-class-selector-beat-id-selector/)

#### 7.关于伪类 LVHA 的解释?

a 标签有四种状态：链接访问前、链接访问后、鼠标滑过、激活，分别对应四种伪类:link、:visited、:hover、:active；

当链接未访问过时：
（1）当鼠标滑过 a 链接时，满足:link 和:hover 两种状态，要改变 a 标签的颜色，就必须将`:hover 伪类在:link 伪类后面声明`；
（2）当鼠标点击激活 a 链接时，同时满足:link、:hover、:active 三种状态，要显示 a 标签激活时的样式（:active），`必须将:active 声明放到:link 和:hover 之后`。因此得出 LVHA 这个顺序。

当链接访问过时，情况基本同上，只不过需要将:link 换成:visited。
这个顺序能不能变？可以，但也`只有:link 和:visited 可以交换位置`，因为一个链接要么访问过要么没访问过，不可能同时满足，
也就不存在覆盖的问题。

#### 8.CSS3 新增伪类有那些？

（1）elem:nth-child(n)选中父元素下的第 n 个子元素，并且这个子元素的标签名为 elem，n 可以接受具体的数
值，也可以接受函数。

（2）elem:nth-last-child(n)作用同上，不过是从后开始查找。

（3）elem:last-child 选中最后一个子元素。

（4）elem:only-child 如果 elem 是父元素下唯一的子元素，则选中之。

（5）elem:nth-of-type(n)选中父元素下第 `n 个 elem 类型元素`，n 可以接受具体的数值，也可以接受函数。

（6）elem:first-of-type 选中父元素下第一个 elem 类型元素。

（7）elem:last-of-type 选中父元素下最后一个 elem 类型元素。

（8）elem:only-of-type 如果父元素下的子元素只有一个 elem 类型元素，则选中该元素。

（9）elem:empty 选中不包含子元素和内容的 elem 类型元素。

（10）elem:target 选择当前活动的 elem 元素。

（11）:not(elem)选择非 elem 元素的每个元素。

（12）:enabled 控制表单控件的禁用状态。

（13）:disabled 控制表单控件的禁用状态。

(14):checked 单选框或复选框被选中。

详细的资料可以参考：
[《CSS3 新特性总结(伪类)》](https://www.cnblogs.com/SKLthegoodman/p/css3.html)
[《浅谈 CSS 伪类和伪元素及 CSS3 新增伪类》](https://blog.csdn.net/zhouziyu2011/article/details/58605705)

#### 9.如何居中 div？

-水平居中：给 div 设置一个宽度，然后添加 margin:0 auto 属性

```CSS
div {
  width: 200px;
  margin: 0 auto;
}
```

-水平垂直居中二

```CSS
/*未知容器的宽高，利用`transform`属性*/
div {
  position: absolute; /*相对定位或绝对定位均可*/
  width: 500px;
  height: 300px;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background-color: pink; /*方便看效果*/
}
```

-水平垂直居中三

```CSS
/*利用flex布局实际使用时应考虑兼容性*/
.container {
  display: flex;
  align-items: center; /*垂直居中*/
  justify-content: center; /*水平居中*/
}

.containerdiv {
  width: 100px;
  height: 100px;
  background-color: pink; /*方便看效果*/
}
```

回答：

一般常见的几种居中的方法有 3 种：

对于宽高固定的元素

（1）我们可以利用 margin:0 auto 来实现元素的水平居中。

（2）利用绝对定位，先将元素的左上角通过 top:50%和 left:50%定位到页面的中心，然后再通过 translate 来调整元素的中心点到页面的中心。

（3）使用 flex 布局，通过 align-items:center 和 justify-content:center 设置容器的垂直和水平方向上为居中对
齐，然后它的子元素也可以实现垂直和水平的居中。

对于宽高不定的元素，上面的后面两种方法，可以实现元素的垂直和水平的居中。

#### 10.display 有哪些值？说明他们的作用。

block 块类型。默认宽度为`父元素宽度`，可设置宽高，换行显示。
none 元素不显示，并从文档流中移除。
inline 行内元素类型。默认宽度为内容宽度，不可设置宽高，同行显示。
inline-block 默认宽度为内容宽度，可以设置宽高，同行显示。
list-item 像块类型元素一样显示，并添加样式列表标记。
`inherit 规定应该从父元素继承 display 属性的值。`

详细资料可以参考：
[《CSS display 属性》](http://www.w3school.com.cn/css/pr_class_display.asp)

#### 11.position 的值 relative 和 absolute 定位原点是？

相关知识点：

**absolute**
生成绝对定位的元素，相对于值不为 static 的第一个父元素的 padding box 进行定位，也可以理解为离自己这一级元素最近的一级 position 设置为 absolute 或者 relative 的父元素的 padding box 的左上角为原点的。
**如果没有父元素设置绝对或相对定位，那么元素相对于根元素也就是 html 元素定位**

**fixed**（老 IE 不支持）
生成绝对定位的元素，相对于浏览器窗口进行定位。

**relative**
生成相对定位的元素，相对于其`元素本身`所在正常位置进行定位。

**static**
默认值。没有定位，元素出现在正常的流中（忽略 top,bottom,left,right,z-index 声明）。

**inherit**
规定从父元素继承 position 属性的值。

回答：

relative 定位的元素，是相对于元素本身的正常位置来进行定位的。

absolute 定位的元素，是`相对于它的第一个 position 值不为 static 的祖先元素的 padding box 来进行定位的`。这句话我们可以这样来理解，我们首先需要找到绝对定位元素的一个 position 的值不为 static 的祖先元素，然后相对于这个祖先元
素的 padding box 来定位，也就是说在计算定位距离的时候，`padding 的值也要算进去。`

#### 12.CSS3 有哪些新特性？（根据项目回答）

`新增各种 CSS 选择器` （`:not(.foo)：所有 class 不是“foo”的节点`）
`圆角` （border-radius:8px）
多列布局 （multi-column layout）
阴影和反射 （Shadow\Reflect）
文字特效 （text-shadow）
文字渲染 （Text-decoration）
线性渐变 （gradient）
`旋转 （transform）`
缩放，定位，倾斜，`动画`，多背景
例如：transform:\scale(0.85,0.90)\translate(0px,-30px)\skew(-9deg,0deg)\Animation:

#### 13.请解释一下 CSS3 的 Flex box（弹性盒布局模型），以及适用场景？

相关知识点：

Flex 是 FlexibleBox 的缩写，意为"弹性布局"，用来为盒状模型提供最大的灵活性。

任何一个容器都可以指定为 Flex 布局。行内元素也可以使用 Flex 布局。注意，设为 Flex 布局以后，`子元素的 float、clear 和 vertical-align 属性将失效`。

采用 Flex 布局的元素，称为 Flex 容器（`flex container`），简称"容器"。它的所有子元素自动成为容器成员，称为 Flex 项目（`flex item`），简称"项目"。

容器默认存在两根轴：水平的主轴（main axis）和垂直的交叉轴（cross axis），项目默认沿主轴排列。

以下 6 个属性设置在容器上。
flex-direction 属性决定主轴的方向（即项目的排列方向）。
flex-wrap 属性定义，如果一条轴线排不下，如何换行。
flex-flow 属性是 flex-direction 属性和 flex-wrap 属性的简写形式，默认值为 row nowrap。
justify-content 属性定义了`项目在主轴上的对齐方式`。
align-items 属性定义项目在`交叉轴上如何对齐`。
align-content 属性定义了多根轴线的对齐方式。如果项目只有一根轴线，该属性不起作用。

以下 6 个属性设置在项目上。
order 属性定义项目的排列顺序。数值越小，排列越靠前，默认为 0。
`flex-grow` 属性定义项目的放大比例，`默认为 0`，即如果存在剩余空间，也不放大。
`flex-shrink` 属性定义了项目的缩小比例，`默认为 1`，即如果空间不足，该项目将缩小。
`flex-basis` 属性定义了在分配多余空间之前，项目占据的主轴空间。浏览器根据这个属性，计算主轴是否有多余空间。它的默认值为 `auto`，即项目的本来大小。

btw,flex 属性是 flex-grow，flex-shrink 和 flex-basis 的简写，默认值为 0 1 auto。

`align-self` 属性允许单个项目有与其他项目不一样的对齐方式，可覆盖 align-items 属性。默认值为 auto，表示继承父元素的 align-items 属性，如果没有父元素，则等同于 stretch。

回答：
flex 布局是 CSS3 新增的一种布局方式，我们可以通过将一个元素的 display 属性值设置为 flex 从而使它成为一个 flex 容器，它的所有子元素都会成为它的项目。
一个容器默认有两条轴，一个是水平的主轴，一个是与主轴垂直的交叉轴。我们可以使用 flex-direction 来指定主轴的方向。我们可以使用 justify-content 来指定元素在主轴上的排列方式，使用 align-items 来指定元素在交叉轴上的排列方式。还可以使用 flex-wrap 来规定当一行排列不下时的换行方式。
对于容器中的项目，我们可以使用 order 属性来指定项目的排列顺序，还可以使用 flex-grow 来指定当排列空间有剩余的时候，项目的放大比例。还可以使用 flex-shrink 来指定当排列空间不足时，项目的缩小比例。

详细资料可以参考：
[《Flex 布局教程：语法篇》](http://www.ruanyifeng.com/blog/2015/07/flex-grammar.html)
[《Flex 布局教程：实例篇》](http://www.ruanyifeng.com/blog/2015/07/flex-examples.html)

#### 14.用纯 CSS 创建一个三角形的原理是什么？

```css
采用的是相邻边框连接处的均分原理。
  将元素的宽高设为0，只设置
  border，把任意三条边隐藏掉（颜色设为transparent），剩下的就是一个三角形。
  #demo {
  width: 0;
  height: 0;
  border-width: 20px;
  border-style: solid;
  border-color: transparent transparent red transparent;
}
```

#### 18.li 与 li 之间有看不见的空白间隔是什么原因引起的？有什么解决办法？

浏览器会把 `inline 元素间的空白字符（空格、换行、Tab 等）渲染成一个空格`。而为了美观。我们通常是一个<li>放在一行，这导致<li>`换行后产生换行字符，它变成一个空格`，占用了一个字符的宽度。

解决办法：

（1）为<li>设置 float:left。不足：有些容器是不能设置浮动，如左右切换的焦点图等。

（2）将所有<li>写在同一行。不足：代码不美观。

（3）`将<ul>内的字符尺寸直接设为 0，即 font-size:0`。不足：<ul>中的其他字符尺寸也被设为 0，`需要额外重新设定其他 字符尺寸`，且在 Safari 浏览器依然会出现空白间隔。

（4）消除<ul>的字符间隔 letter-spacing:-8px，不足：这也设置了<li>内的字符间隔，因此需要将<li>内的字符
间隔设为默认 letter-spacing:normal。

详细资料可以参考：
[《li 与 li 之间有看不见的空白间隔是什么原因引起的？》](https://blog.csdn.net/sjinsa/article/details/70919546)

#### 22.width:auto 和 width:100%的区别

设置了 padding 和 margin 后:

width:100%会使元素 `box 的宽度等于父元素的 content box 的宽度`。
**会发生内容溢出父节点的情况**

width:auto 会使`元素撑满整个父元素`，margin、border、padding、content 区域会自动分配水平空间。浏览器会自己选择一个合适的宽度值,不用担心当元素自身有 margin、padding 、border 时，宽度会超过父节点。

#### 27.对 BFC 规范（块级格式化上下文：block formatting context）的理解？

相关知识点：

块格式化上下文（Block Formatting Context，BFC）是 Web 页面的可视化 CSS 渲染的一部分，是布局过程中生成块级盒子的区域，也是浮动元素与其他元素的交互限定区域。

通俗来讲

•`BFC 是一个独立的布局环境，可以理解为一个容器，在这个容器中按照一定规则进行物品摆放，并且不会影响其它环境中的物品`。
•如果一个元素符合触发 BFC 的条件，则 BFC 中的元素布局不受外部影响。

创建 BFC

（1）根元素或包含根元素的元素
（2）浮动元素 float ＝ left|right 或 inherit（≠none）
（3）`绝对定位元素` position ＝ absolute 或 fixed
（4）display ＝ inline-block|flex|inline-flex|table-cell 或 table-caption
（5）`overflow ＝ hidden|auto 或 scroll(≠visible)`

回答：

BFC 指的是块级格式化上下文，一个元素形成了 BFC 之后，那么它内部元素产生的布局不会影响到外部元素，外部元素的布局也不会影响到 BFC 中的内部元素。一个 BFC 就像是一个隔离区域，和其他区域互不影响。

一般来说根元素是一个 BFC 区域，浮动和绝对定位的元素也会形成 BFC，display 属性的值为 inline-block、flex 这些属性时也会创建 BFC。还有就是元素的 overflow 的值不为 visible 时都会创建 BFC。

详细资料可以参考：
[《深入理解 BFC 和 MarginCollapse》](https://www.w3cplus.com/css/understanding-bfc-and-margin-collapse.html)
[《前端面试题-BFC（块格式化上下文）》](https://segmentfault.com/a/1190000013647777)

#### 29.请解释一下为什么需要清除浮动？清除浮动的方式

浮动元素可以左右移动，直到遇到另一个浮动元素或者遇到它外边缘的包含框。浮动框不属于文档流中的普通流，当元素浮动之后，
不会影响块级元素的布局，只会影响内联元素布局。此时文档流中的普通流就会表现得该浮动框不存在一样的布局模式。当包含框
的高度小于浮动框的时候，此时就会出现“高度塌陷”。

`清除浮动是为了清除使用浮动元素产生的影响。浮动的元素，高度会塌陷，而高度的塌陷使我们页面后面的布局不能正常显示。`

清除浮动的方式

（1）使用 clear 属性清除浮动。参考 28。

（2）使用 BFC 块级格式化上下文来清除浮动。参考 26。

因为 BFC 元素不会影响外部元素的特点，所以 BFC 元素也可以用来清除浮动的影响，因为如果不清除，子元素浮动则父元素高度塌陷，必然会影响后面元素布局和定位，这显然有违 BFC 元素的子元素不会影响外部元素的设定。

#### 30.使用 clear 属性清除浮动的原理？

使用 clear 属性清除浮动，其语法如下：

clear:none|left|right|both

如果单看字面意思，clear:left 应该是“清除左浮动”，clear:right 应该是“清除右浮动”的意思，实际上，这种解释是有问题的，因为浮动一直还在，并没有清除。

官方对 clear 属性的解释是：“`元素盒子的边不能和前面的浮动元素相邻`。”，我们对元素设置 clear 属性是为了`避免浮动元素对该元素的影响`，而不是清除掉浮动。

还需要注意的一点是 clear 属性指的是元素盒子的边不能和前面的浮动元素相邻，注意这里“前面的”3 个字，也就是 clear 属性对“后面的”浮动元素是不闻不问的。考虑到 float 属性要么是 left，要么是 right，不可能同时存在，同时由于 clear 属性对“后面的”浮动元素不闻不问，因此，当 clear:left 有效的时候，clear:right 必定无效，也就是此时 clear:left 等同于设置 clear:both；同样地，clear:right 如果有效也是等同于设置 clear:both。由此可见，clear:left 和 clear:right 这两个声明就没有任何使用的价值，至少在 CSS 世界中是如此，`直接使用 clear:both 吧`。

一般使用伪元素的方式清除浮动

```CSS
.clear::after{
  content:'';
  display:block;//也可以是'block'，或者是'list-item'
  clear:both;
}
```

`clear 属性只有块级元素才有效的，而::after 等伪元素默认都是内联水平`，这就是借助伪元素清除浮动影响时需要设置 display 属性值的原因。

#### 33.使用 CSS 预处理器吗？喜欢哪个？

SASS（SASS、LESS 没有本质区别，只因为团队前端都是用的 SASS）

#### 35.浏览器是怎样解析 CSS 选择器的？

**样式系统从关键选择器开始匹配，然后左移查找规则选择器的祖先元素**。只要选择器的子树一直在工作，样式系统就会持续左移，直到和规则匹配，或者是因为不匹配而放弃该规则。

试想一下，`如果采用从左至右的方式读取 CSS 规则，那么大多数规则读到最后（最右）才会发现是不匹配的`，这样做会费时耗能，最后有很多都是无用的；而如果采取从右向左的方式，那么只要发现`最右边选择器不匹配，就可以直接舍弃了，避免了许多无效匹配`。

哈希表维护底层结点,倒过来搜???
理由:

- CSS 选择器 对应一个 word ,相当于在 Trie 中搜索结点
- 从父亲到子孙有很多条路,但是从子孙到父亲只有一条向上的路

详细资料可以参考：
[《探究 CSS 解析原理》](https://juejin.im/entry/5a123c55f265da432240cc90)

#### 36.在网页中应该使用奇数还是偶数的字体？为什么呢？

偶数
（1）`偶数字号相对更容易和 web 设计的其他部分构成比例关系`。比如：当我用了 14px 的正文字号，我可能会在一些地方用 14
×0.5=7px 的 margin，在另一些地方用 14×1.5=21px 的标题字号。
（2）浏览器缘故，`低版本的浏览器 ie6 会把奇数字体强制转化为偶数`，即 13px 渲染为 14px。
（3）系统差别，早期的 Windows 里，中易宋体点阵只有 12 和 14、15、16px，唯独缺少 13px。

详细资料可以参考：
[《谈谈网页中使用奇数字体和偶数字体》](https://blog.csdn.net/jian_xi/article/details/79346477)
[《现在网页设计中的为什么少有人用 11px、13px、15px 等奇数的字体？》](https://www.zhihu.com/question/20440679)

#### 37.margin 和 padding 分别适合什么场景使用？

`margin 是用来隔开元素与元素的间距`；padding 是用来隔开元素与内容的间隔。
margin 用于布局分开元素使元素与元素互不相干。
`padding 用于元素与内容之间的间隔`，让内容（文字）与（包裹）元素之间有一段距离。

何时应当使用 margin：
•需要在 border 外侧添加空白时。
•空白处不需要背景（色）时。
•上下相连的两个盒子之间的空白，需要相互抵消时。如 15px+20px 的 margin，将得到 20px 的空白。

何时应当时用 padding：
•需要在 border 内测添加空白时。
•空白处需要背景（色）时。
•上下相连的两个盒子之间的空白，希望等于两者之和时。如 15px+20px 的 padding，将得到 35px 的空白。

#### 38.抽离样式模块怎么写，说出思路，有无实践经验？[阿里航旅的面试题]

我的理解是把常用的 css 样式单独做成 css 文件……通用的和业务相关的分离出来，通用的做成样式模块儿共享，业务相关的，放进业务相关的库里面做成对应功能的模块儿。

详细资料可以参考：
[《CSS 规范-分类方法》](http://nec.netease.com/standard/css-sort.html)

#### 51.设备像素、css 像素、设备独立像素、dpr、ppi 之间的区别？

`设备像素指的是物理像素`，一般手机的分辨率指的就是设备像素，一个设备的设备像素是不可变的。

`css 像素和设备独立像素是等价的，不管在何种分辨率的设备上，css 像素的大小应该是一致的`，css 像素是一个相对单位，它是相对于设备像素的，一个 css 像素的大小取决于页面缩放程度和 dpr 的大小。

`dpr 指的是设备像素和设备独立像素的比值 Device Pixels Rate`，一般的 pc 屏幕，dpr=1。在 iphone4 时，苹果推出了 retina 屏幕，它的 dpr 为 2。屏幕的缩放会改变 dpr 的值。

ppi 指的是每英寸的物理像素的密度 `Pixels Per Inch`，ppi 越大，屏幕的分辨率越大。

详细资料可以参考：
[《什么是物理像素、虚拟像素、逻辑像素、设备像素，什么又是 PPI,DPI,DPR 和 DIP》](https://www.cnblogs.com/libin-1/p/7148377.html)
[《前端工程师需要明白的「像素」》](https://www.jianshu.com/p/af6dad66e49a)
[《CSS 像素、物理像素、逻辑像素、设备像素比、PPI、Viewport》](https://github.com/jawil/blog/issues/21)
[《前端开发中像素的概念》](https://github.com/wujunchuan/wujunchuan.github.io/issues/15)

#### 58.png、jpg、gif 这些图片格式解释一下，分别什么时候用。有没有了解过 webp？

相关知识点：

（1）BMP，是无损的、既支持索引色也支持直接色的、点阵图。这种图片格式几乎没有对数据进行压缩，所以 BMP 格式的图片通常具有较大的文件大小。

（2）GIF 是无损的、采用索引色的、点阵图。采用 LZW 压缩算法进行编码。文件小，是 GIF 格式的优点，同时，GIF 格式还具有支持动画以及透明的优点。但，GIF 格式仅支持 8bit 的索引色，所以 GIF 格式适用于对色彩要求不高同时需要文件体积较小的场景。

（3）JPEG 是有损的、采用直接色的、点阵图。JPEG 的图片的优点，是采用了直接色，`得益于更丰富的色彩，JPEG 非常适合用来存储照片`，与 GIF 相比，JPEG 不适合用来存储企业 Logo、线框类的图。因为有损压缩会导致图片模糊，而直接色的选用，又会导致图片文件较 GIF 更大。

（4）PNG-8 是无损的、使用索引色的、点阵图。PNG 是一种比较新的图片格式，PNG-8 是非常好的 GIF 格式替代者，在可能的
情况下，应该尽可能的使用 PNG-8 而不是 GIF，因为在相同的图片效果下，PNG-8 具有更小的文件体积。除此之外，PNG-8
还支持透明度的调节，而 GIF 并不支持。现在，除非需要动画的支持，否则我们没有理由使用 GIF 而不是 PNG-8。

（5）PNG-24 是无损的、使用直接色的、点阵图。PNG-24 的优点在于，它压缩了图片的数据，使得同样效果的图片，PNG-24 格
式的文件大小要比 BMP 小得多。当然，PNG24 的图片还是要比 JPEG、GIF、PNG-8 大得多。

（6）SVG 是无损的、矢量图。SVG 是矢量图。这意味着 SVG 图片由直线和曲线以及绘制它们的方法组成。当你放大一个 SVG 图
片的时候，你看到的还是线和曲线，而不会出现像素点。这意味着 SVG 图片在放大时，不会失真，`所以它非常适合用来绘制企业 Logo、Icon 等。`

（7）`WebP 是谷歌开发的一种新图片格式`，WebP 是同时支持有损和无损压缩的、使用直接色的、点阵图。从名字就可以看出来它是
为 Web 而生的，什么叫为 Web 而生呢？`就是说相同质量的图片，WebP 具有更小的文件体积`。现在网站上充满了大量的图片，如果能够降低每一个图片的文件大小，那么将大大减少浏览器和服务器之间的数据传输量，进而降低访问延迟，提升访问体验。

•在无损压缩的情况下，相同质量的 WebP 图片，文件大小要比 PNG 小 26%；
•在有损压缩的情况下，具有相同图片精度的 WebP 图片，文件大小要比 JPEG 小 25%~34%；
•WebP 图片格式支持图片透明度，一个无损压缩的 WebP 图片，如果要支持透明度只需要 22%的格外文件大小。

但是目前只有 Chrome 浏览器和 Opera 浏览器支持 WebP 格式，兼容性不太好。

回答：

我了解到的一共有七种常见的图片的格式。

（1）第一种是 BMP 格式，它是无损压缩的，支持索引色和直接色的点阵图。由于它基本上没有进行压缩，因此它的文件体积一般比较大。

（2）第二种是 GIF 格式，它是无损压缩的使用索引色的点阵图。由于使用了 LZW 压缩方法，因此文件的体积很小。并且 GIF 还支持动画和透明度。但因为它使用的是索引色，所以它适用于一些对颜色要求不高且需要文件体积小的场景。

（3）第三种是 JPEG 格式，它是`有损压缩`的使用直接色的点阵图。由于使用了直接色，色彩较为丰富，一般适用于来存储照片。但
由于使用的是直接色，可能文件的体积相对于 GIF 格式来说更大。

（4）第四种是 PNG-8 格式，它是`无损压缩`的使用索引色的点阵图。它是 GIF 的一种很好的替代格式，它也支持透明度的调整，并且文件的体积相对于 GIF 格式更小。一般来说如果不是需要动画的情况，我们都可以使用 PNG-8 格式代替 GIF 格式。

（5）第五种是 PNG-24 格式，它是无损压缩的使用直接色的点阵图。PNG-24 的优点是它使用了压缩算法，所以它的体积比 BMP
格式的文件要小得多，但是相对于其他的几种格式，还是要大一些。

（6）第六种格式是 svg 格式，它是矢量图，它记录的图片的绘制方式，因此对矢量图进行放大和缩小不会产生锯齿和失真。它一般
适合于用来制作一些网站 logo 或者图标之类的图片。

（7）第七种格式是 webp 格式，它是支持有损和无损两种压缩方式的使用直接色的点阵图。使用 webp 格式的最大的优点是，在相
同质量的文件下，它拥有更小的文件体积。因此它非常适合于网络图片的传输，因为图片体积的减少，意味着请求时间的减小，
这样会提高用户的体验。这是谷歌开发的一种新的图片格式，目前在兼容性上还不是太好。

详细资料可以参考：
[《图片格式那么多，哪种更适合你？》](https://www.cnblogs.com/xinzhao/p/5130410.html)

#### 61.style 标签写在 body 后与 body 前有什么区别？

页面加载自上而下当然是先加载样式。写在 body 标签后由于浏览器以逐行方式对 HTML 文档进行解析，当解析到写在尾部的样式表（外联或写在 style 标签）`会导致浏览器停止之前的渲染，等待加载且解析样式表完成之后重新渲染`，在 windows 的 IE 下可
能会出现 FOUC 现象（`即样式失效导致的页面闪烁问题`）

#### 62.什么是 CSS 预处理器/后处理器？

CSS 预处理器定义了一种新的语言，其基本思想是，用一种专门的编程语言，为 CSS 增加了一些编程的特性，将 CSS 作为目标生成文件，然后开发者就只要使用这种语言进行编码工作。通俗的说，CSS 预处理器用一种专门的编程语言，进行 Web 页面样式设计，然后再编译成正常的 CSS 文件。

`预处理器例如：LESS、Sass、Stylus`，用来`预编译` Sass 或 less csssprite，增强了 css 代码的复用性，还有层级、mixin、变量、循环、函数等，具有很方便的 UI 组件模块化开发能力，极大的提高工作效率。

CSS 后处理器是对 CSS 进行处理，并最终生成 CSS 的预处理器，它属于广义上的 CSS 预处理器。我们很久以前就在用 CSS 后处理器了，最典型的例子是 CSS 压缩工具（如 clean-css），只不过以前没单独拿出来说过。还有最近比较火的 Autoprefixer，以 CanIUse 上的浏览器支持数据为基础，自动处理兼容性问题。

后处理器例如：PostCSS，通常被视为在完成的样式表中根据 CSS 规范处理 CSS，让其更有效；目前最常做的是给 `CSS 属性添加浏览器私有前缀，实现跨浏览器兼容性的问题。`

详细资料可以参考：
[《CSS 预处理器和后处理器》](https://blog.csdn.net/yushuangyushuang/article/details/79209752)

#### 66.画一条 0.5px 的线

采用 meta viewport 的方式

采用 border-image 的方式

`采用 transform:scale()的方式`

详细资料可以参考：
[《怎么画一条 0.5px 的边（更新）》](https://juejin.im/post/5ab65f40f265da2384408a95)

#### 67.transition 和 animation 的区别

`transition 关注的是 CSS property 的变化`，property 值和时间的关系是一个`三次贝塞尔曲线。`

`animation 作用于元素本身而不是样式属性`，可以使用关键帧的概念，应该说可以实现更自由的动画效果。

详细资料可以参考：
[《CSSanimation 与 CSStransition 有何区别？》](https://www.zhihu.com/question/19749045)
[《CSS3Transition 和 Animation 区别及比较》](https://blog.csdn.net/cddcj/article/details/53582334)
[《CSS 动画简介》](http://www.ruanyifeng.com/blog/2014/02/css_transition_and_animation.html)
[《CSS 动画：animation、transition、transform、translate》](https://juejin.im/post/5b137e6e51882513ac201dfb)

#### 69.为什么 height:100%会无效？

对于普通文档流中的元素，`百分比高度值要想起作用，其父级必须有一个可以生效的高度值。`

原因是如果包含块的高度没有显式指定（即高度由内容决定），并且该元素不是绝对定位，则计算值为 auto，因为解释成了 auto，所以无法参与计算。

使用绝对定位的元素会有计算值，即使祖先元素的 height 计算为 auto 也是如此。

#### 98.常见的元素隐藏方式？

-（1）使用 display:none;隐藏元素，渲染树不会包含该渲染对象，因此该元素不会在页面中占据位置，也不会响应绑定的监听事件。

-（2）使用 visibility:hidden;隐藏元素。`元素在页面中仍占据空间，但是不会响应绑定的监听事件。`

-（3）使用 opacity:0;将元素的透明度设置为 0，以此来实现元素的隐藏。元素在页面中仍然占据空间，并且能够响应元素绑定的监听事件。

-（4）`通过使用绝对定位将元素移除可视区域内，以此来实现元素的隐藏。`

-（5）通过 `z-index 负值`，来使其他元素遮盖住该元素，以此来实现隐藏。

-（6）通过` clip/clip-path` 元素裁剪的方法来实现元素的隐藏，这种方法下，元素仍在页面中占据位置，但是不会响应绑定的监听事件。

-（7）通过 `transform:scale(0,0)`来将元素缩放为 0，以此来实现元素的隐藏。这种方法下，元素仍在页面中占据位置，但是不会响应绑定的监听事件。

详细资料可以参考：
[《CSS 隐藏元素的八种方法》](https://juejin.im/post/584b645a128fe10058a0d625#heading-2)

#### 99.css 实现上下固定中间自适应布局？

```css
利用绝对定位实现 body {
  padding: 0;
  margin: 0;
}

.header {
  position: absolute;
  top: 0;
  width: 100%;
  height: 100px;
  background: red;
}

.container {
  position: absolute;
  top: 100px;
  bottom: 100px;
  width: 100%;
  background: green;
}

.footer {
  position: absolute;
  bottom: 0;
  height: 100px;
  width: 100%;
  background: red;
}

利用flex布局实现 html,
body {
  height: 100%;
}

body {
  display: flex;
  padding: 0;
  margin: 0;
  flex-direction: column;
}

.header {
  height: 100px;
  background: red;
}

.container {
  flex-grow: 1;
  background: green;
}

.footer {
  height: 100px;
  background: red;
}
```

详细资料可以参考：
[《css 实现上下固定中间自适应布局》](https://www.jianshu.com/p/30bc9751e3e8)

#### 100.css 两栏布局的实现？

相关资料：

```css
/*两栏布局一般指的是页面中一共两栏，左边固定，右边自适应的布局，一共有四种实现的方式。*/
/*以左边宽度固定为200px为例*/

/*（1）利用浮动，将左边元素宽度设置为200px，并且设置向左浮动。将右边元素的margin-left设置为200px，宽度设置为auto（默认为auto，撑满整个父元素）。*/
.outer {
  height: 100px;
}

.left {
  float: left;

  height: 100px;
  width: 200px;

  background: tomato;
}

.right {
  margin-left: 200px;

  width: auto;
  height: 100px;

  background: gold;
}

/*（2）第二种是利用flex布局，将左边元素的放大和缩小比例设置为0，基础大小设置为200px。将右边的元素的放大比例设置为1，缩小比例设置为1，基础大小设置为auto。*/
.outer {
  display: flex;

  height: 100px;
}

.left {
  flex-shrink: 0;
  flex-grow: 0;
  flex-basis: 200px;

  background: tomato;
}

.right {
  flex: auto;
  /*11auto*/

  background: gold;
}

/*（3）第三种是利用绝对定位布局的方式，将父级元素设置相对定位。左边元素设置为absolute定位，并且宽度设置为
200px。将右边元素的margin-left的值设置为200px。*/
.outer {
  position: relative;

  height: 100px;
}

.left {
  position: absolute;

  width: 200px;
  height: 100px;

  background: tomato;
}

.right {
  margin-left: 200px;
  height: 100px;

  background: gold;
}

/*（4）第四种还是利用绝对定位的方式，将父级元素设置为相对定位。左边元素宽度设置为200px，右边元素设置为绝对定位，左边定位为200px，其余方向定位为0。*/
.outer {
  position: relative;

  height: 100px;
}

.left {
  width: 200px;
  height: 100px;

  background: tomato;
}

.right {
  position: absolute;

  top: 0;
  right: 0;
  bottom: 0;
  left: 200px;

  background: gold;
}
```

[《两栏布局 demo 展示》](http://cavszhouyou.top/Demo-Display/TwoColumnLayout/index.html)

回答：

两栏布局一般指的是页面中一共两栏，左边固定，右边自适应的布局，一共有四种实现的方式。

以左边宽度固定为 200px 为例

-（1）利用浮动，`将左边元素宽度设置为 200px，并且设置向左浮动。将右边元素的 margin-left 设置为 200px，宽度设置为 auto`（默认为 auto，撑满整个父元素）。

-（2）第二种是利用 flex 布局，`将左边元素的放大和缩小比例设置为 0，基础大小设置为 200px。将右边的元素的放大比例设置为 1，缩小比例设置为 1，基础大小设置为 auto。`

-（3）第三种是利用绝对定位布局的方式，将父级元素设置相对定位。`左边元素设置为 absolute 定位，并且宽度设置为 200px。将右边元素的 margin-left 的值设置为 200px。`

-（4）第四种还是利用绝对定位的方式，将父级元素设置为相对定位。左边元素宽度设置为 200px，右边元素设置为绝对定位，左边定位为 200px，其余方向定位为 0。

#### 101.css 三栏布局的实现？

相关资料：

```css
/*三栏布局一般指的是页面中一共有三栏，左右两栏宽度固定，中间自适应的布局，一共有五种实现方式。

这里以左边宽度固定为100px，右边宽度固定为200px为例。*/

/*（1）利用绝对定位的方式，左右两栏设置为绝对定位，中间设置对应方向大小的margin的值。*/
.outer {
  position: relative;

  height: 100px;
}

.left {
  position: absolute;

  width: 100px;
  height: 100px;
  background: tomato;
}

.right {
  position: absolute;
  top: 0;
  right: 0;

  width: 200px;
  height: 100px;
  background: gold;
}

.center {
  margin-left: 100px;
  margin-right: 200px;
  height: 100px;
  background: lightgreen;
}

/*（2）利用flex布局的方式，左右两栏的放大和缩小比例都设置为0，基础大小设置为固定的大小，中间一栏设置为auto*/
.outer {
  display: flex;
  height: 100px;
}

.left {
  flex: 00100px;
  background: tomato;
}

.right {
  flex: 00200px;
  background: gold;
}

.center {
  flex: auto;
  background: lightgreen;
}

/*（3）利用浮动的方式，左右两栏设置固定大小，并设置对应方向的浮动。中间一栏设置左右两个方向的margin值，注意这种方式，中间一栏必须放到最后。*/
.outer {
  height: 100px;
}

.left {
  float: left;
  width: 100px;
  height: 100px;
  background: tomato;
}

.right {
  float: right;
  width: 200px;
  height: 100px;
  background: gold;
}

.center {
  height: 100px;
  margin-left: 100px;
  margin-right: 200px;
  background: lightgreen;
}

/*（4）圣杯布局，利用浮动和负边距来实现。父级元素设置左右的 padding，三列均设置向左浮动，中间一列放在最前面，宽度设置为父级元素的宽度，因此后面两列都被挤到了下一行，通过设置 margin 负值将其移动到上一行，再利用相对定位，定位到两边。*/
.outer {
  height: 100px;
  padding-left: 100px;
  padding-right: 200px;
}

.left {
  position: relative;
  left: -100px;

  float: left;
  margin-left: -100%;

  width: 100px;
  height: 100px;
  background: tomato;
}

.right {
  position: relative;
  left: 200px;

  float: right;
  margin-left: -200px;

  width: 200px;
  height: 100px;
  background: gold;
}

.center {
  float: left;

  width: 100%;
  height: 100px;
  background: lightgreen;
}

/*（5）双飞翼布局，双飞翼布局相对于圣杯布局来说，左右位置的保留是通过中间列的 margin 值来实现的，而不是通过父元
素的 padding 来实现的。本质上来说，也是通过浮动和外边距负值来实现的。*/

.outer {
  height: 100px;
}

.left {
  float: left;
  margin-left: -100%;

  width: 100px;
  height: 100px;
  background: tomato;
}

.right {
  float: left;
  margin-left: -200px;

  width: 200px;
  height: 100px;
  background: gold;
}

.wrapper {
  float: left;

  width: 100%;
  height: 100px;
  background: lightgreen;
}

.center {
  margin-left: 100px;
  margin-right: 200px;
  height: 100px;
}
```

[《三栏布局 demo 展示》](http://cavszhouyou.top/Demo-Display/ThreeColumnLayout/index.html)

回答：

三栏布局一般指的是页面中一共有三栏，左右两栏宽度固定，中间自适应的布局，一共有五种实现方式。

这里以左边宽度固定为 100px，右边宽度固定为 200px 为例。

（1）利用绝对定位的方式，左右两栏设置为绝对定位，中间设置对应方向大小的 margin 的值。

（2）利用 flex 布局的方式，左右两栏的放大和缩小比例都设置为 0，基础大小设置为固定的大小，中间一栏设置为 auto。

（3）利用浮动的方式，左右两栏设置固定大小，并设置对应方向的浮动。中间一栏设置左右两个方向的 margin 值，注意这种方式，中间一栏必须放到最后。

（4）圣杯布局，利用浮动和负边距来实现。父级元素设置左右的 padding，三列均设置向左浮动，中间一列放在最前面，宽度设置为父级元素的宽度，因此后面两列都被挤到了下一行，通过设置 margin 负值将其移动到上一行，再利用相对定位，定位到两边。圣杯布局中间列的宽度不能小于两边任意列的宽度，而双飞翼布局则不存在这个问题。

（5）双飞翼布局，双飞翼布局相对于圣杯布局来说，左右位置的保留是通过中间列的 margin 值来实现的，而不是通过父元素的 padding 来实现的。本质上来说，也是通过浮动和外边距负值来实现的。

#### 102.实现一个宽高自适应的正方形

```css
/*1.第一种方式是利用vw来实现*/
.square {
  width: 10%;
  height: 10vw;
  background: tomato;
}

/*2.第二种方式是利用元素的margin/padding百分比是相对父元素width的性质来实现*/
.square {
  width: 20%;
  height: 0;
  padding-top: 20%;
  background: orange;
}

/*3.第三种方式是利用子元素的margin-top的值来实现的*/
.square {
  width: 30%;
  overflow: hidden;
  background: yellow;
}

.square::after {
  content: '';
  display: block;
  margin-top: 100%;
}
```

[《自适应正方形 demo 展示》](http://cavszhouyou.top/Demo-Display/AdaptiveSquare/index.html)
