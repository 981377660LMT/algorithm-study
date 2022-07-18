#### 8. 行内元素与块级元素的区别？

`HTML4` 中，元素被分成两大类：inline （内联元素）与 block （块级元素）。

（1） 格式上，默认情况下，行内元素不会以新行开始，而块级元素会新起一行。
（2） 内容上，默认情况下，`行内元素只能包含文本和其他行内元素`。而块级元素可以包含行内元素和其他块级元素。
（3） 行内元素与块级元素属性的不同，主要是盒模型属性上：
行内元素设置 width 无效，height 无效（可以设置 line-height），`设置 margin 和 padding 的上下不会对其他元素产生影响。`

#### 12. 页面导入样式时，使用 link 和 @import 有什么区别？

（1）从属关系区别。 @import 是 `CSS 提供的语法规则`，只有导入样式表的作用；link 是 HTML 提供的标签，不仅可以加
载 CSS 文件，还可以定义 `RSS、rel 连接属性、引入网站图标`等。

（2）`加载顺序`区别。加载页面时，link 标签引入的 CSS 被同时加载；@import 引入的 CSS 将在页面加载完毕后被加载。

（3）兼容性区别。@import 是 CSS2.1 才有的语法，故只可在 IE5+ 才能识别；link 标签作为 HTML 元素，不存在兼容
性问题。

（4）DOM 可控性区别。可以通过 JS 操作 DOM ，插入 link 标签来改变样式；由于 DOM 方法是基于文档的，无法使用 @import 的方式插入样式。

#### 18. 渲染过程中遇到 JS 文件怎么处理？（浏览器解析过程）

JavaScript 的加载、解析与执行会阻塞文档的解析，也就是说，`在构建 DOM 时，HTML 解析器若遇到了 JavaScript，那么它会暂停文档的解析，将控制权移交给 JavaScript 引擎，等 JavaScript 引擎运行完毕，浏览器再从中断的地方恢复继续解析文档。`

也就是说，如果你想首屏渲染的越快，就越不应该在首屏就加载 JS 文件，这也是都建议将 script 标签`放在 body 标签底部的原因`。
当然在当下，`并不是说 script 标签必须放在底部，因为你可以给 script 标签添加 defer 或者 async 属性。`

#### 19. async 和 defer 的作用是什么？有什么区别？（浏览器解析过程）

（1）脚本没有 defer 或 async，浏览器会立即加载并执行指定的脚本，也就是说不等待后续载入的文档元素，读到就加载并执行。

（2）defer 属性表示延迟执行引入的 JavaScript，`即这段 JavaScript 加载时 HTML 并未停止解析`，这两个过程是并行的。
`当整个 document 解析完毕后再执行脚本文件`，在 DOMContentLoaded 事件触发之前完成。
`多个脚本按顺序执行。`

（3）async 属性表示异步执行引入的 JavaScript，与 defer 的区别在于，如果已经加载好，就会开始执行，
也就是说`它的执行仍然会阻塞文档的解析，只是它的加载过程不会阻塞`。
`多个脚本的执行顺序无法保证。`

详细资料可以参考：
[《defer 和 async 的区别》](https://segmentfault.com/q/1010000000640869)

#### 24. 什么是重绘和回流？（浏览器绘制过程）

重绘: 当渲染树中的一些元素需要更新属性，而这些属性只是影响元素的外观、风格，而不会影响布局的操作，比如 background-color，我们将这样的操作称为重绘。

回流：当渲染树中的一部分（或全部）因为元素的规模尺寸、布局、隐藏等改变而需要重新构建的操作，会影响到布局的操作，这的操作我们称为回流。

常见引起回流属性和方法：

任何会改变元素几何信息（元素的位置和尺寸大小）的操作，都会触发回流。

（1）添加或者删除可见的 DOM 元素；
（2）元素尺寸改变——边距、填充、边框、宽度和高度
（3）内容变化，比如用户在 input 框中输入文字
（4）浏览器窗口尺寸改变——resize 事件发生时
（5）计算 offsetWidth 和 offsetHeight 属性
（6）设置 style 属性的值
（7）当你修改网页的默认字体时。

回流必定会发生重绘，重绘不一定会引发回流。回流所需的成本比重绘高的多，改变父节点里的子节点很可能会导致父节点的一系列回流。

常见引起重绘属性和方法：

![常见引起回流属性和方法](https://cavszhouyou-1254093697.cos.ap-chongqing.myqcloud.com/note-14.png)

常见引起回流属性和方法：

![常见引起重绘属性和方法](https://cavszhouyou-1254093697.cos.ap-chongqing.myqcloud.com/note-13.png)

详细资料可以参考：
[《浏览器的回流与重绘》](https://juejin.im/post/5a9923e9518825558251c96a)

#### 25. 如何减少回流？（浏览器绘制过程）

（1）使用 transform 替代 top
（2）把 DOM 离线后修改。如：使用 documentFragment 对象在内存里操作 DOM
（3）`不要一条一条地修改 DOM 的样式`。与其这样，还不如预先定义好 css 的 class，然后修改 DOM 的 className。

#### 26. 为什么操作 DOM 慢？（浏览器绘制过程）

一些 DOM 的操作或者属性访问可能会引起页面的`回流和重绘`，从而引起性能上的消耗。

#### 27. DOMContentLoaded 事件和 Load 事件的区别？

当初始的 `HTML 文档被完全加载和解析完成之后，DOMContentLoaded 事件被触发`，而无需等待样式表、图像和子框架的加载完成。

Load 事件是当`所有资源加载完成后触发的`。

详细资料可以参考：
[《DOMContentLoaded 事件 和 Load 事件的区别？》](https://www.jianshu.com/p/ca8dae435a2c)

#### 28. HTML5 有哪些新特性、移除了那些元素？

新增的有：

绘画 canvas;
用于媒介回放的 video 和 audio 元素;
本地离线`存储` localStorage 长期存储数据，浏览器关闭后数据不丢失;
sessionStorage 的数据在浏览器关闭后自动删除;
`语意化更好的内容元素`，比如 article、footer、header、nav、section;
表单控件，calendar、date、time、email、url、search;
`新的技术` webworker, websocket;
新的文档属性 document.visibilityState

移除的元素有：

纯表现的元素：basefont，big，center，font, s，strike，tt，u;
对可用性产生负面影响的元素：frame，frameset，noframes；

#### 32. 前端需要注意哪些 SEO ？

（1）合理的 `title、description`、keywords：搜索对着三项的权重逐个减小，title 值强调重点即可，重要关键词出现不要超过 2 次，而且要靠前，不同页面 title 要有所不同；description 把页面内容高度概括，长度合适，不可过分堆砌关键词，不同页面 description 有所不同；keywords 列举出重要关键词即可。

（2）`语义化的 HTML 代码`，符合 W3C 规范：语义化代码让搜索引擎容易理解网页。

（3）重要内容 HTML 代码放在最前：搜索引擎抓取 HTML 顺序是从上到下，有的搜索引擎对抓取长度有限制，保证重要内容肯定被抓取。

（5）少用 iframe：搜索引擎不会抓取 iframe 中的内容

（6）`非装饰性图片必须加 alt`

（7）提高网站速度：网站速度是搜索引擎排序的一个重要指标

#### 38. Label 的作用是什么？是怎么用的？

label 标签来定义表单控制间的关系，当用户选择该标签时，浏览器会自动将焦点转到和标签相关的表单控件上。

<label for="Name">Number:</label>
<input type=“text“ name="Name" id="Name"/>

#### 44. 实现不使用 border 画出 1 px 高的线，在不同浏览器的标准模式与怪异模式下都能保持一致的效果。

html

<div style="height:1px;overflow:hidden;background:red"></div>

#### 46. `<img>` 的 title 和 alt 有什么区别？

title 通常当鼠标滑动到元素上的时候显示

alt 是 <img> 的特有属性，是图片内容的等价描述，用于图片无法加载时显示、读屏器阅读图片。可提图片高可访问性，除了纯装
饰图片外都必须设置有意义的值，搜索引擎会重点分析。

#### 47. Canvas 和 SVG 有什么区别？

Canvas 是一种通过 JavaScript 来绘制 2D 图形的方法。`Canvas 是逐像素来进行渲染的，因此当我们对 Canvas 进行缩放时，会出现锯齿或者失真的情况。`

SVG 是一种`使用 XML 描述 2D 图形的语言`。SVG 基于 XML，这意味着 SVG DOM 中的每个元素都是可用的。我们可以为某个元素附加 JavaScript 事件监听函数。并且 `SVG 保存的是图形的绘制方法，因此当 SVG 图形缩放时并不会失真。`

详细资料可以参考：
[《SVG 与 HTML5 的 canvas 各有什么优点，哪个更有前途？》](https://www.zhihu.com/question/19690014)

#### 48. 网页验证码是干嘛的，是为了解决什么安全问题？

（1）区分用户是计算机还是人的公共全自动程序。可以防止恶意破解密码、刷票、论坛灌水
（2）有效防止黑客对某一个特定注册用户用特定程序暴力破解方式进行不断的登陆尝试 `CSRF`

#### 49. 渐进增强和优雅降级的定义

渐进增强：`针对低版本浏览器进行构建页面，保证最基本的功能`，然后再针对高级浏览器进行效果、交互等改进和追加功能达到更好的用户体验。

优雅降级：一开始就根据高版本浏览器构建完整的功能，然后再`针对低版本浏览器进行兼容`。

#### 50. attribute 和 property 的区别是什么？

attribute 是 dom 元素在文档中作为 `html 标签拥有的属性`；
property 就是 dom 元素在 `js 中作为对象拥有的属性`。
对于 html 的标准属性来说，attribute 和 property 是同步的，是会自动更新的，
但是对于自定义的属性来说，他们是不同步的。
(props,attrs)

#### 52. IE 各版本和 Chrome 可以并行下载多少个资源？

（1） IE6 2 个并发
（2） iE7 升级之后的 6 个并发，之后版本也是 6 个
（3） `Firefox，chrome 也是 6 个`

#### 55. 浏览器架构

- 用户界面
  - 主进程
  - 内核
    - 渲染引擎
    - JS 引擎
      - 执行栈
    - 事件触发线程
      - 消息队列
        - 微任务
        - 宏任务
    - 网络异步线程
    - 定时器线程

#### 56. 常用的 meta 标签

 <meta> 元素可提供有关页面的元信息（meta-information），比如针对搜索引擎和更新频度的描述和关键词。
 <meta> 标签位于文档的头部，不包含任何内容。<meta> 标签的属性定义了与文档相关联的名称/值对。

 <!DOCTYPE html>  H5标准声明，使用 HTML5 doctype，不区分大小写
 <head lang="en"> 标准的 lang 属性写法
 <meta charset="utf-8">    声明文档使用的字符编码
 <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>   优先使用 IE 最新版本和 Chrome
 <meta name="description" content="不超过150个字符"/>       页面描述
 <meta name="keywords" content=""/>      页面关键词者
 <meta name="author" content="name, email@gmail.com"/>    网页作
 <meta name="robots" content="index,follow"/>      搜索引擎抓取
 <meta name="viewport" content="initial-scale=1, maximum-scale=3, minimum-scale=1, user-scalable=no"> 为移动设备添加 viewport
 <meta http-equiv="X-UA-Compatible" content="IE=edge">     避免IE使用兼容模式

`设置页面不缓存`

 <meta http-equiv="pragma" content="no-cache">
 <meta http-equiv="cache-control" content="no-cache">
 <meta http-equiv="expires" content="0">

详细资料可以参考：
[《Meta 标签用法大全》](http://www.cnblogs.com/qiumohanyu/p/5431859.html)

#### 57. css reset 和 normalize.css 有什么区别？

css reset 是最早的一种解决浏览器间样式不兼容问题的方案，它的基本思想是将浏览器的所有样式都`重置掉`，从而达到所有浏览器样式保持一致的效果。但是使用这种方法，可能会带来一些性能上的问题，并且对于一些元素的不必要的样式的重置，其实反而会造成画蛇添足的效果。

后面出现一种更好的解决浏览器间样式不兼容的方法，就是 normalize.css ，它的思想是尽量的`保留浏览器自带的样式(求同存异)`，通过在原有的样式的基础上进行调整，来保持各个浏览器间的样式表现一致。相对与 css reset，normalize.css 的方法保留了有价值的默认值，并且修复了一些浏览器的 bug，而且使用 normalize.css 不会造成元素复杂的继承链。

#### 60. head 标签中必不少的是？

<title>

#### 62. 在 HTML5 中，哪个方法用于获得用户的当前位置？

```js
window.navigator.geolocation.getCurrentPosition(function (position) {
  var pos = {
    lat: position.coords.latitude,
    lng: position.coords.longitude,
  }
  console.log(pos)
})
```

#### 66. 前端性能优化？

前端性能优化主要是为了提高页面的加载速度，优化用户的访问体验。我认为可以从这些方面来进行优化。

第一个方面是`页面`的内容方面

（1）通过文件合并、css 雪碧图、使用 base64 等方式来减少 `HTTP 请求数`，避免过多的请求造成等待的情况。
（2）通过 DNS 缓存等机制来减少 `DNS 的查询次数`。
（3）通过设置缓存策略，对常用不变的资源进行`缓存`。
（4）使用`延迟加载`的方式，来减少页面首屏加载时需要请求的资源。延迟加载的资源当用户需要访问时，再去请求加载
（5）通过用户行为，对某些资源使用`预加载`的方式，来提高用户需要访问资源时的响应速度。

第二个方面是`服务器`方面

（1）使用` CDN` 服务，来提高用户对于资源请求时的响应速度。
（2）服务器端启用 Gzip、Deflate 等方式对于传输的资源进行压缩，减小文件的体积。
（3）`尽可能减小 cookie 的大小`，并且通过将静态资源分配到其他域名下，来避免对静态资源请求时携带不必要的 cookie

第三个方面是 CSS 和 JavaScript 方面

（1）把样式表放在页面的 head 标签中，减少页面的首次渲染的时间。
（2）`避免使用 @import` 标签。
（3）尽量把 js 脚本放在页面底部或者使用 `defer 或 async` 属性，避免脚本的加载和执行阻塞页面的渲染。
（4）通过对 JavaScript 和 CSS 的文件进行`压缩`，来减小文件的体积。
