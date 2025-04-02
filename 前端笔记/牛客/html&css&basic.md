1. 怎样判断一个引用类型变量是 Proxy 的实例
   **改写实例的 Symbol.toStringTag 属性**

   ```JS
   // 没有任何基构造函数具有指定数量的类型参数
   class Foo extends Proxy {}  // 错误的写法！！！
   ```

2. 子网掩码
   请来一位先生主持公道，先生这么来操作：
   192.168.1.0-192.168.1.63 分给大儿子
   192.168.1.64-192.168.1.127 分给二儿子
   192.168.1.128 -192.168.1.191 分给三儿子
   192.168.1.192-192.168.1.255 分给四儿子
   那如何来描述四个儿子的子网网段呢？
   192.168.1.0/26
   192.168.1.64/26
   192.168.1.128/26
   192.168.1.192/26
   26 个 1，为了简化子网掩码的表示，用/26 代替 255.255.255.192。
   子网掩码的左边是网络位，用二进制数字“1”表示，1 的数目等于网络位的长度；右边是主机位，用二进制数字“0”表示，0 的数目等于主机位的长度。
   255.255.255.0 ‐＞ 11111111.11111111.11111111.00000000

   我们来看大儿子的网段 192.168.1.0/26 是如何得到的？
   以大儿子的房间为例：
   192.168.1.0-192.168.1.63
   二进制表示：
   192.168.1.0
   11000000.10101000.00000001.00000000
   11000000.10101000.00000001.00111111
   `首尾地址完全相同`的是多少？
   11000000.10101000.00000001.00
   数数一共多少位？26 位！那么用这个“192.168.1.0/26”就可以表示大儿子所有房间。

   ip 地址分为两部分，前面部分叫“网络位”，类似于学号中的班级编号，只要两个 ip 的网络位相同，那么他们就是同网段。

   ip 地址后面的部分叫“主机位”，代表这个主机的个体编号，类似学号中老师每次点名喊的数字。

   IP 地址与子网掩码相与得到网络号：

3. NodeJS 的 Buffer 是堆外内存
4. 所谓的标准模式是指，浏览器按 W3C 标准解析执行代码；怪异模式则是使用浏览器自己的方式解析执行代码，因为不同浏览器解析执行的方式不一样，所以我们称之为怪异模式。
   浏览器解析时到底使用标准模式还是怪异模式，与你网页中的 DTD（文档类型定义） 声明直接相关， DTD 声明定义了标准文档的类型（标准模式解析）文档类型，会使浏览器使用相应的方式加载网页并显示，忽略 DTD 声明 ,将使网页进入怪异模式。
   Doctype，知道这是干什么的吗？
   声明位于文档中的最前面的位置，处于标签之前。此标签可告知浏览器文档使用哪种 HTML 或 XHTML 规范。（重点：**告诉浏览器按照何种规范解析页面**）
5. 对 HTML5 认识?（是什么,为什么）
   HTML5 指的是包括 HTML 、 CSS 和 JavaScript 在内的一套技术组合。
   HTML5 是 HTML 最新版本， 2014 年 10 月由万维网联盟（ W3C ）完成标准制定。目标是替换 1999 年所制定的 HTML 4.01 和 XHTML 1.0 标准，以期能在互联网应用迅速发展的时候，使网络标准达到匹配当代的网络需求。
6. cookies，sessionStorage 和 localStorage 的区别?
7. link 和@import 的区别?
   下面介绍一下 `<link>` 和 `@import` 在引入 CSS 文件时的主要区别：

   1. **使用方式**

      - `<link>`：是在 HTML 文件的 `<head>` 部分通过 HTML 标签引入外部 CSS 文件，例如：

        ```html
        <!-- filepath: /path/to/index.html -->
        <head>
          <link rel="stylesheet" href="styles.css" />
        </head>
        ```

      - `@import`：是在 CSS 文件或 `<style>` 标签中使用的 CSS 语法，用于引入其他 CSS 文件，例如：

        ```css
        /* filepath: /path/to/styles.css */
        @import url('other-styles.css');
        ```

   2. **加载顺序和性能**

      - `<link>`：浏览器在解析 HTML 时会立即开始下载 `<link>` 指定的 CSS 文件，并行加载多个 CSS 文件，性能通常更好。
      - `@import`：会延迟加载 CSS 文件，且多个 `@import` 必须串行加载，会影响页面的渲染速度，尤其在 IE 等老版本浏览器中性能较差。

   3. **浏览器兼容性**

      - `<link>`：所有现代浏览器都支持，兼容性非常好。
      - `@import`：也可以正常使用，但某些早期版本的浏览器可能存在兼容性问题，且与 `<link>` 比较时效率较低。

   4. **维护和组织代码**
      - `<link>`：结构清晰，所有外部样式在 HTML 中明确声明，便于查找和维护。
      - `@import`：可以将样式拆分并合并到一个 CSS 中，但过多嵌套会导致管理难度加大。

   **记忆小技巧：**  
   把 `<link>` 想象成“直接连接”外部样式，速度快；而 `@import` 则是“间接导入”，需要经过额外的步骤，因此加载上可能略慢。

   总结来说，在需要提高页面加载性能和兼容性的场景下，推荐使用 `<link>` 标签引入外部样式。

8. 说说你对 SVG 理解?
   SVG 可缩放矢量图形（ Scalable Vector Graphics ）是基于可扩展标记语言（ XML ），用于描述二维矢量图形的一种图形格式
   SVG 严格遵从 XML 语法，并用文本格式的描述性语言来描述图像内容，因此是一种和图像分辨率无关的矢量图形格式
9. `data-`属性的作用是什么？
   `data-`为 H5 新增的为前端开发者提供自定义的属性，这些属性集可以通过对象的 `dataset` 属性获取，不支持该属性的浏览器可以通过 `getAttribute` 方法获取
10. 介绍一下你对浏览器内核的理解？
    **渲染引擎**(layout engineer 或 Rendering Engine), 浏览器的内核的不同对于网页的语法解释会有不同，所以渲染的效果也不相同; 和 **JS 引擎**,解析和执行 javascript 来实现网页的动态效果。
    最开始渲染引擎和 JS 引擎并没有区分的很明确，后来 JS 引擎越来越独立，内核就倾向于只指渲染引擎。
11. iframe 有那些优缺点
    优点：

    1. 解决加载缓慢的第三方内容如图标和广告等的`加载问题`
    2. `并行加载脚本`
       缺点:
    3. iframe 会阻塞主页面的 Onload 事件
    4. 搜索引擎的检索程序无法解读这种页面，不利于 SEO;

12. 简述一下 src 与 href 的区别？

    `src` 和 `href` 都是在 HTML 中用于指定资源位置的属性，但它们有不同的用途和行为：

    - **src（source）**

      - 用于嵌入外部资源，并在当前位置“注入”该资源内容，比如 `<img>`、`<script>`、`<iframe>` 等标签。
      - 浏览器会下载并替换当前元素的内容。如果资源加载失败，通常会导致该元素无法正常显示或执行。

    - **href（hypertext reference）**
      - 用于指定链接地址，通常用于 `<a>` 标签表示超链接，也用于 `<link>` 标签引入样式表、图标等。
      - 浏览器会根据需要进行导航或下载资源，不会直接在当前位置嵌入内容。加载失败时不会影响当前文档的主体显示。

    **总结记忆要点：**

    - `src` 是 “source”，`直接嵌入资源`。
    - `href` 是 “hyperlink reference”，用来`指向链接或资源地址`。

    ```HTML
    <script src ='js.js'></script>

    <link href='common.css' rel='stylesheet'/>
    ```

13. HTTP 协议决定了服务器与客户端之间的连接方式，无法直接实现消息推送（ F5 已坏） , 一些变相的解决办法：
    双向通信与消息推送：客户端定时向服务器发送 Ajax 请求，服务器接到请求后马上返回响应信息并关闭连接；适于小型应用
    长轮询：客户端向服务器发送 Ajax 请求，服务器接到请求后 hold 住连接，直到有新消息才返回响应信息并关闭连接，客户端处理完响应信息后再向服务器发送新的请求。WebQQ、 Hi 网页版、 Facebook IM 。
    长连接：在页面里嵌入一个隐蔵 iframe，将这个隐蔵 iframe 的 src 属性设为对一个长连接的请求或是采用 xhr 请求，服务器端就能源源不断地往客户端输入数据。Gmail 聊天
    Websocket:能够实现真正意义上的推送功能
14. 表单的基本组成部分有哪些，表单的主要用途是什么？
    组成：表单标签、表单域、表单按钮
    主要用途：表单在网页中主要负责数据采集的功能，和向服务器传送数据。
15. HTML5 标准提供了哪些新的 API？
    Data Transfer API
    History API
    Media API
    Text Track API
    User Interaction
16. HTML5 应用程序缓存和浏览器缓存有什么区别
    应用程序缓存是 HTML5 的重要特性之一，提供了离线使用的功能让应用程序可以获取本地的网站内容，例如 HTML 、 CSS 、图片以及 JavaScript 。这个特性可以提高网站性能，它的实现借助于 manifest 文件，如下：

```HTML

<!doctype html>
<html manifest=”example.appcache”>
</html>

与传统浏览器缓存相比，它不强制用户访问的网站内容被缓存。
```

20. 兄弟选择器 ~ 相邻兄弟选择器+
    子选择器 >
    https://segmentfault.com/a/1190000016563980
21. 请列举几种清除浮动的方法(至少两种)?

- 结尾处加**空 div 标签** clear:both
  此方法是以前主要使用的一种清除浮动方法
- 父级 div 定义 **伪类:after**
  推荐使用，建议定义公共类，以减少 CSS 代码。
  ```CSS
  #parent:after{
    content:" ";
    height:0; // 该方法中必须为需要清除浮动元素的伪对象中设置 0，否则该元素会比实际高出若干像素；
    visibility:hidden;
    display:block;
    clear:both;
  }
  ```
- **BFC**:父级 div 设置`overflow`为`hidden`或者 auto
  只推荐没有使用 position 或对 overflow:hidden 理解比较深的朋友使用。

22. display:inline
    • display:inline
    a、inline 元素不会独占一行，多个相邻的行内元素会排列在同一行里，直到一行排列不下，才会新换一行，其宽度随元素的内容而变化。
    b、inline 元素设置 width,height 属性无效。
    c、inline 元素的 margin 和 padding 属性，水平方向的 padding-left, padding-right, margin-left, margin-right 都产生边距效果；但竖直方向的 padding-top, padding-bottom, margin-top, margin-bottom 不会产生边距效果。

23. em 与 rem 的重要区别： 它们计算的规则一个是依赖父元素另一个是依赖根元素计算
24. box-sizing 属性的的用法
    box-sizing:content-box： **padding 和 border 不被包含在定义的 width 和 height 之内**
    box-sizing:border-box： padding 和 border 被包含在定义的 width 和 height 之内
25. 浏览器标准模式和怪异模式之间的区别是什么？
    标准模式是指，浏览器按 W3C 标准解析执行代码；怪异模式则是使用浏览器自己的方式解析执行代码
    浏览器解析时到底使用标准模式还是怪异模式，与你网页中的 DTD 声明直接相关
    设置行内元素的高宽： 在 Standards 模式下，给<span>等行内元素设置 wdith 和 height 都不会生效，而在 quirks 模式下，则会生效。
26. 说说你对边距折叠的理解?
    a、参加折叠的 margin 都是正值：取其中 margin 较大的值为最终 margin 值。
    b、参与折叠的 margin 都是负值：取的是其中`绝对值较大`的，然后，从 0 位置，负向位移。
    c、参与折叠的 margin 中有正值，有负值：`先取出负 margin 中绝对值中最大的，然后，和正 margin 值中最大的 margin 相加。`

27. 如何产生 BFC
    a、float 的值不为 none
    b、overflow 的值不为 visible
    c、display 的值为 table-cell, table-caption, inline-block 中的任何一个
    d、position 的值不为 relative 和 static
28. 如何解决多个元素重叠问题？
    z-index 属性
29. 说一下 http 和 https
30. WebSocket 的实现和应用
    WebSocket 是基于 Http 协议的，或者说借用了 Http 协议来完成一部分握手，在握手阶段与 Http 是相同的。我们来看一个 websocket 握手协议的实现，基本是 2 个属性，upgrade，connection。
31. head 请求和 options 请求
    head：类似于 get 请求，只不过返回的响应中没有具体的内容，**用户获取报头**
    options：允许客户端查看服务器的性能，比如说**服务器支持的请求方式**等等。
32. 一个图片 url 访问后直接下载怎样实现？
    - 设置 response header 中的 disposition-type 为 attachment
    - h5 新增的 download 属性
    ```HTML
    <a href="/images/logo.png" download="/images/logo.png">
        <img border="0" src="/images/logo.png" alt="runoob.com" >
    </a>
    ```
33. 几个很实用的 BOM 属性对象方法?
    Bom 是浏览器对象。有哪些常用的 Bom 属性呢？
34. drag API
    dragstart：事件主体是被拖放元素，在开始拖放被拖放元素时触发，。
    drag：事件主体是被拖放元素，在正在拖放被拖放元素时触发。
    dragend：事件主体是被拖放元素，在整个拖放操作结束时触发
    dragenter：事件主体是目标元素，在被拖放元素进入某元素时触发。
    dragover：事件主体是目标元素，在被拖放在某元素内移动时触发。
    dragleave：事件主体是目标元素，在被拖放元素移出目标元素是触发。
    drop：事件主体是目标元素，在目标元素完全接受被拖放元素时触发。
35. 说一下 http2.0

    HTTP2.0 是基于 1999 年发布的 HTTP1.0 之后的首次更新。
    新特性

    1. 提升访问速度：相比 HTTP1.0，请求资源所需时间更少，访问速度更快。
    2. 允许多路复用：多路复用允许同时通过单一的 HTTP/2 连接发送多重请求-响应信息。改善了在 HTTP1.1 中，浏览器客户端在同一时间，针对同一域名下的请求有一定数量限制（连接数量），超过限制会被阻塞。
    3. 二进制分帧：HTTP2.0 会将所有的传输信息分割为更小的信息或者帧，并对他们进行二进制编码。
    4. 首部压缩：使报头更紧凑，更快速传输，有利于移动网络环境。
    5. 服务器端推送（server push）：还没有收到浏览器的请求，服务器就把各种资源推送给浏览器。 比如，浏览器只请求了 index.html，但是服务器把 index.html、style.css、example.png 全部发送给浏览器。

36. 400:Bad Request 请求报文中存在语法错误
    401:Unauthorized 例如 token 不对
    403:Forbidden 没有该权限
37. 一句话概括 RESTFUL
    就是**用 URL 定位资源，用 HTTP 描述操作**
38. fetch 发送 2 次请求的原因
    fetch 方法对于 HTTP GET 方式不会发送两次请求。
    fetch 发送 post 请求(复杂请求)需要发送两次的原因是，第一次是向服务器发送一个 options 请求询问服务器是否`支持修改请求头`(此时返回状态码 No Content 204)，如果可以才在第二次中发送真正的请求

39. 说一下 web worker
    Web Worker 的作用：为单线程的 JS 创造多线程环境。 43. iframe 是什么？有什么缺点？
    定义：iframe 元素会创建包含另一个文档的内联框架
    提示：可以将提示文字放在<iframe></iframe>之间，来提示某些不支持 iframe 的浏览器
    缺点：

        会阻塞主页面的 onload 事件

        搜索引擎无法解读这种页面，不利于 SEO

        iframe 和主页面共享连接池，而浏览器对相同区域有限制所以会影响性能。

40. 讲讲 viewport 和移动端布局

    ```HTML
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1">

    ```

    content 里面的这些属性必须知道
    width: 控制 viewport 的大小，可以是一个指定的值，或者是特殊值，如 device-width 是为设备指定的宽度。

    height： 和 width 相对应，指定高度

    initial-scale： 初始缩放比例，也是当前页面第一次 load 的时候缩放的比例。

    maximum-scale： 允许用户缩放的最大比例。

    minimum-scale： 允许用户缩放的最小比例。

    user-scalable： 用户是否可以手动缩放。

41. click 在 ios 上有 300ms 延迟，原因及如何解决？
    原因：ios 上有双击缩放的功能，点击一次屏幕浏览器无法判断用户是想要进行单击操作还是双击操作，所以要等待 300ms 是否有第二次点击
    解决方法：
    (1)粗暴型，禁用缩放
    <meta name="viewport" content="width=device-width, user-scalable=no">

    (2)利用 FastClick，其原理是：

    检测到 touchend 事件后，`立刻出发模拟 click 事件，并且把浏览器 300 毫秒之后真正出发的事件给阻断掉`

42. addEventListener 参数
    addEventListener(event, function, useCapture)
    其中，event 指定事件名；function 指定要事件触发时执行的函数；useCapture 指定事件是否在捕获或冒泡阶段执行。
43. cookie session 区别
    参考回答：

    1.  cookie 数据存放在客户的浏览器上，session 数据放在服务器上。
    2.  cookie 不是很安全，别人可以分析存放在本地的 COOKIE 并进行 COOKIE 欺骗.考虑到安全应当使用 session。
    3.  session 会在一定时间内保存在服务器上。当访问增多，会比较占用你服务器的性能.考虑到减轻服务器性能方面，应当使用 COOKIE。
    4.  单个 cookie 保存的数据不能超过 4K，很多浏览器都限制一个站点最多保存 20 个 cookie。

44. http 常用请求头
45. 常见状态码
    https://blog.csdn.net/qq_44647809/article/details/115276258
46. 什么是 FOUC（无样式内容闪烁）？你如何来避免 FOUC？
    FOUC - Flash Of Unstyled Content 文档样式闪烁
    引用 CSS 文件的@import 就是造成这个问题的罪魁祸首。I
    IE 会先加载整个 HTML 文档的 DOM，然后再去导入外部的 CSS 文件，因此，在页面 DOM 加载完成到 CSS 导入完成中间会有一段时间页面上的内容是没有样式的，这段时间的长短跟网速，电脑速度都有关系。
    解决方法简单的出奇，只要在<head>之间加入一个<link>或者<script>元素就可以了。
47. js 延迟加载的方式有哪些？
    1. defer 和 async
    2. 动态创建 DOM 方式（创建 script，插入到 DOM 中，加载完毕后 callBack）
    3. 按需异步载入 js
48. document.write 和 innerHTML 的区别
    document.write 只能重绘整个页面

        innerHTML 可以重绘页面的一部分

49. 哪些操作会造成内存泄漏？
    1. setTimeout 的第一个参数使用字符串而非函数的话，会引发内存泄漏。
    2. 闭包
    3. 控制台日志
    4. 循环（在两个对象彼此引用且彼此保留时，就会产生一个循环）

Stacking Context（堆叠上下文）是浏览器中用于管理元素在 z 轴上的堆叠顺序的一种概念。简单来说，它决定了哪些元素显示在其他元素之上或之下。下面详细说明：

---

### 1. 什么是 Stacking Context

- **定义：**  
  Stacking Context 是一个独立的渲染层，在该层内，子元素的堆叠顺序只受内部规则控制，与外部的堆叠环境隔离。

- **作用：**  
  它确保了一组元素按照特定规则进行层叠显示，而不会被外部层级影响，从而减少了布局冲突。

---

### 2. 形成 Stacking Context 的条件

一个元素会形成独立的 Stacking Context，主要有以下几种情况：

- 根元素（HTML）总是一个 Stacking Context。
- 具有 `position` 属性（相对、绝对或固定定位）且 `z-index` 属性不为 `auto` 的元素。
- 使用了 CSS 属性 `opacity` 且值小于 1 的元素。
- 应用了 CSS 变换，如 `transform`, `filter` 或 `perspective` 的元素。
- 使用 `mix-blend-mode` 或 `isolation: isolate` 的元素。
- 其他部分属性（如 `clip-path`、`mask` 等）在特定条件下也会触发新的 Stacking Context。

---

### 3. Stacking Context 内的堆叠顺序

在一个 Stacking Context 内，元素的堆叠顺序按以下规则排列（从背景到前景）：

1. 背景和边框
2. 负 z-index 的子元素
3. 块级格式化上下文内普通内容
4. z-index 为 0 的子元素
5. 正 z-index 的子元素

注意：不同 Stacking Context 之间是互相独立的，外部元素不会直接影响内部元素的层级关系。

---

### 4. 记忆技巧

- **“上下文”概念：** 就像一个独立的小世界，子元素在这个小世界内有自己的“排名”，不受外界影响。
- **触发条件记忆法：** “定位+z-index、非 1 的透明度、变换效果”可以迅速提示是否会形成新的 Stacking Context。

---

理解 Stacking Context 有助于解决页面中出现的重叠问题与层级混乱，掌握其规则可以让你更精确地控制元素的显示顺序。
