1. 某主机的 IP 地址为 212.212.77.55，子网掩码为 255.255.252.0。若该主机向其所在子网发送`广播`分组，则目的地址可以是？
   由子网掩码可知前 22 位为`子网号`、后 10 位为`主机号`
   `主机号全部置为 1 即广播地址`
   212.212.（10001101）.55 => 212.212. 10001111.255
   在使用 TCP/IP 协议的网络中，主机标识段 host ID 为全 1 的 IP 地址为广播地址
2. 标签模板其实不是模板，而是函数调用的一种特殊形式。“标签”指的就是函数，紧跟在后面的模板字符串就是它的参数。

```JS
let a = 5; let b = 10; tag`Hello ${ a + b } world ${ a * b }`; // 等同于 tag(['Hello ', ' world ', ''], 15, 50);
```

3. 关于将 Promise.all 和 Promise.race 传入空数组的两段代码的输出结果说法正确的是：

```JS
Promise.all([]).then((res) => {
    console.log('all');
});
Promise.race([]).then((res) => {
    console.log('race');
});
```

all 会被输出，而 race 不会被输出
Promise.all([ ])中，数组为空数组，则立即决议为成功执行 resolve( )；
Promise.race([ ])中数组为空数组，就不会执行，永远挂起

4. 内存泄漏是 javascript 代码中必须尽量避免的，以下几段代码可能会引起内存泄漏的有（）

意外的全局变量

```JS
function getName() {
    name = 'javascript'
}
getName()
```

脱离 DOM 的引用
此时，仍旧存在一个全局的 #button 的引用 elements 字典。
button 元素仍旧在内存中，不能被 GC 回收。
`应该直接 elements.button=null`

```JS
const elements = {
    button: document.getElementById('button')
};
function removeButton() {
    document.body.removeChild(elements.button);
}
removeButton()
```

1. 身在乙方的小明去甲方对一网站做渗透测试，甲方客户告诉小明，目标站点由 wordpress 搭建，请问下面操作最为合适的是
   使用 wpscan 对网站进行扫描
2. 发生冲突的情况下，max-height 会覆盖掉 height, min-height 又会覆盖掉 max-height
3. 关于网络请求延迟增大的问题，以下哪些描述是正确的()
   使用 ping 来测试 TCP 端口是不是可以正常连接 (x) `ping 是网络层的，tcp 是传输层的`
   使用 tcpdump 抓包分析网络请求 (√)
   使用 strace 观察进程的网络系统调用 (√)
   使用 Wireshark 分析网络报文的收发情况 (√)
4. JSON.stringify 方法，返回值 res 分别是

   ```JS
   const fn = function(){}
   const res = JSON.stringify(fn)
   const num = 123
   const res = JSON.stringify(num)
   const res = JSON.stringify(NaN)
   const b = true
   const res = JSON.stringify(b)

   undefined、'123'、'null'、'true'
   ```

   1. 转换值如果有 toJSON() 方法，该方法定义什么值将被序列化。
   2. 函数、undefined 被单独转换时，会返回 undefined，如 JSON.stringify(function(){}) or JSON.stringify(undefined).
   3. 对包含循环引用的对象（对象之间相互引用，形成无限循环）执行此方\*\*\*抛出错误。
   4. NaN 和 Infinity 格式的数值及 null 都会被当做 null。
   5. undefined、任意的函数以及 symbol 值，在序列化过程中会被忽略（出现在非数组对象的属性值中时）或者被转换成 null（出现在数组中时）。

5. SASS/SCSS 依赖 Ruby 或 Node.js 编译环境，因此需要编译后才能在浏览器使用；而 less 本身由 JavaScript 实现，可以在浏览器中完成编译并可直接使用
6. 若一个哈夫曼树有 N 个叶子节点，则其节点总数为 2N-1

7. 分页存储管理将进程的逻辑地址空间分成若干个页，并为各页加以编号，从 0 开始，若某一计算机主存按字节编址，逻辑地址和物理地址都是 32 位，`页表项`大小为 4 字节，若使用一级页表的分页存储管理方式，逻辑地址结构为页号（20 位），`页内偏移量`（12 位），则页的大小是（4KB ）？页表最大占用（4MB ）？

地址长度为 32 位，其中 0~11 位为页内地址（即页内偏移量），2^12 即`每页大小为 4KB；`
同样地，12~31 位为页号，地址空间最多允许有 2^20 = 1M 页，又页表项 4 字节， 所以页表最大占用` 1M * 4 = 4MB`

12. PostCSS
    假定在项目使用 less 完成 CSS 的预编译，PostCSS 进行 CSS 代码转换的时间`晚于` less
13. 以下哪种方法可以用来清理僵尸进程()
    `向僵尸进程的父进程发送 SIGKILL 信号`

14. 类

```JS
function Person() { }
var person = new Person();
```

- 每一个原型都有一个 constructor 属性指向关联的构造函数。
- person.constructor === Person (去`person.__proto__`上找 constructor)
- Object.getPrototypeOf(person) === Person.prototype

15. resolve 同步会立即执行
<!-- 2 3 5 4 1 -->

```JS
setTimeout(function(){
    console.log(1);
}, 0)
new Promise(function(resolve){
    console.log(2);
    resolve();
    console.log(3);
}).then(function(){
    console.log(4);
})
console.log(5);

```

16. flex 分配空间问题

- 当多个子盒子的总宽度 flex-basis<=父盒子时，使用的是 flex-grow 属性，按比例分配剩余空间；
- 当多个子盒子的总宽度 flex-basis>父盒子时，使用的是 flex-shrink 属性

```CSS
  .box{
        display:flex;
        width:200px;
        height:50px;
  }

  .left{
      flex:3 2 50px
  }

  .right{
      flex: 2 1 200px
  }

  计算规则：
      (1)计算超出父盒子的宽度：200+50-200 = 50;
      (2).left需要减少：(50*2)/(50*2+200*1) * 50 = 50/3
          .right需要减少：(200*1)/(50*2+200*1) * 50 = 100/3
      (3)最终left宽度：50-50/3 = 100/3
          right宽度：200-100/3 = 500/3
      (4)left和right比例：  1:5
```

17. 小牛开发文件上传功能时，遇到了一些安全问题，那么对于文件上传漏洞，有效防御手段有哪些？
    服务器端限制文件扩展名
    将上传的文件存储在静态文件服务器中
18. 将一个整数序列整理为降序，两趟处理后序列变为{36, 31, 29, 14, 18, 19, 32}则采用的排序算法可能是`____插入排序____`。

19. 如果存储结构由数组变为链表，那么下列哪些算法的时间复杂度量级会升高
    希尔排序
    堆排序
    希尔排序、堆排序使用数组存储的话，方便获取指定位置的数据。这两个排序都需取`指定位置`的数据，而使用链表增加了获取指定位置的时间。
20. FIFO 为先进先出的顺序来完成页面的访问，而如果在采用先进先出页面淘汰算法的系统中，一进程在`内存占 3 块`（开始为空），页面访问序列为 1、2、3、4、1、2、5、1、2、3、4、5、6。运行时会产生（ ）次缺页中断？
    先来先服务利用队列来进行页面读取。队列大小为 3，刚开始队列为空：
    访问 1，`队列中没 1，缺页一次，读入页面 1`
    访问 2，队列中没 2，缺页两次，读入页面 2，队列为 1,2
    访问 3，队列中没 3，缺页三次，读入页面 3，队列为 1,2,3
    访问 4，没 4，缺页 4 次，读入页面 4，队列为 2,3,4
    访问 1，没 1，缺页 5 次，读入页面 1，队列为 3,4,1
    访问 2，没 2，缺页 6 次，读入页面 2，队列为 4,1,2
    访问 5，没 5，缺页 7 次，读入页面 5，队列为 1,2,5
    访问 1，不存在缺页，队列中为 1,2,5
    访问 2，不缺页，队列中为 1,2,5
    访问 3，缺页 8 次……
    总的次数为 10 次

21. require 题

```JS
有a.js和b.js两个文件，请选择b文件中代码的输出
// a.js
let a = 1
let b = {}
setTimeout(() => {
a = 2
b.b = 2
}, 100)
module.exports = { a, b }

// b.js
const a = require('./a')
console.log(a.a)
console.log(a.b)
setTimeout(() => {
console.log(a.a)
console.log(a.b)
}, 500)
```

1 {} 1 {b:2}

`commonjs 导出的是值的拷贝` ，a 所以 a 一直是 1;
b 是浅拷贝，拷贝的是对象的引用,所以 a.js 的 b 改变时,b.js 的 a.b 也改变

22. vue 的 computed 缓存属性

```JS
请选择下面代码输出1的次数
var vm = new Vue({
el: '#example',
data: {
message: 'Hello'
},
computed: {
test: function () {
console.log(1)
return this.message
}
},
created: function (){
this.message = 'World'
for (var i = 0; i < 5; i++) {
console.log(this.test)
}
}
})
```

1 次
因为 vue 的 computed 具有缓存功能。`message 只更新了一次，所以 test 只触发一次`，执行一次 console.log(1)。

23. 以下这种写法不规范，但是不会报错，其在浏览器中的表现形式是

<p>1<p>2</p></p>
```HTML
<p>1</p><p>2</p><p></p>
```

24. 元素的 border 是由三角形组合而成

```CSS
div {
    width: 0;
    height: 0;
    border: 40px solid;
    border-color: orange blue red green;
}
```

![](image/note/1646028724276.png)

25. 标签中使用多个 class，不看这些 class 添加的顺序，而是看 style 中定义的顺序。
    数字 `1` 和 `2` 被浏览器渲染出来的颜色分别是是？

```HTML
<html>
  <head>
    <style>
      .classA { color: blue; }
      .classB { color: red; }
    </style>
  </head>
  <body>
    <p class='classB classA'>1</p>
    <p class='classA classB'>2</p>
  </body>
</html>
```

26. 当网站对<script>标签进行过滤时，可以通过哪种方式进行绕过且有效攻击

```HTML
<img src="" onerror=alert(1)>
```

27. linux 文件权限
    Linux 的每个文件一般都有三个权限 r--读，w--写，x--执行，其分别对应的数值为 4，2，1。
    修改/home 下 test 目录以及目录下所有文件，可以支持所有人可读可写的，以下能实现的有?
    `chmod 666 /home/test -R`
    `chmod 777 /home/test -R`

28. 下列选项中，可能导致当前 linux`进程阻塞`?

- 进程申请临界资源
- 进程 从磁盘读数据

29. transition 与 animation 的区别
    - transition 着重属性的变化，而 animation 重点是在创建帧，让不同帧在不同时间点发生不同变化；
30. 关于 HTML 的描述，不推荐的是
    可以使用 center 标签来设置元素居中；

    - 纯渲染层面，建议是 CSS

31. 在面向对象技术中，多态性是指（）
    针对一消息，不同对象可以以适合自身的方式加以响应
32. 操作系统的基本特征：并发、共享、虚拟、异步

33. 关于 node.js 中的模块化规范，以下说法正确的有哪些

- require 加载模块是一个`同步`的过程
- require 函数可以在代码的任意位置执行
- exports 或 module.exports 其中一个一旦重新赋值，exports 将失效

34. Web Worker 为 Web 内容在后台线程中运行脚本提供了一种简单的方法。请列举出 Web worker 的常用 API 并列举至少 1 个 Web Worker 的常见用途。

```
常用API：
  1. new Worker(url)，用于创建一个worker实例，url指向一个js文件，浏览器会创建一个单独的线程来执行这个文件
  2. worker.prototype.postMessage()，用于从worker向主线程传递消息，第一个参数是被传递的消息，可以传递对象/基础类型的数据
  3. onmessage / addEventListener('message', callback)，可以用在主线程或worker上，用于监听message事件，接收对方传递来的消息，消息被放置在事件对象的data属性中

常见用途：
  1. 处理密集型数***算
  2. 大数据排序
  3. 数据处理，如压缩、音频处理等
  4. 用于执行网络操作，如Ajax、WebSocket
```

35. 200 与 202
    200 OK 代表请求没问题，响应实体的主体部分`包含了所请求的资源`
    202 Accepted 代表请求已接受，但服务器还`未对其执行任何动作`
36. p 标签是块级元素，常理来说，块级元素是可以嵌套块级元素和行内元素的，但是 p 标签是个特殊，它里面不能嵌套块级元素
p 标签内`遇到下一个块级元素的标签会自动结束，`
例如<p>p 开始<div>xxxxx</div>p 结束</p>，会被解析为
<p>p开始</p><div>xxxxx</div><p>p结束</p>，相当于写了两个p段落，这不是我们想要的结果

37. 下面语句的执行结果是？

```JS
'a.b.c'.replace(/(.)\.(.)\.(.)/, '$2.$1.$0')
```

b.a.$0
$0 是没有特殊意义的

38. Number 与 parseInt
    Number(undefined) NaN
    Number(null) 0
    parseInt(undefined) NaN
    parseInt(null) NaN
39. Math.round 正数时 0.5 向上取，负数时 0.5 向下取
    Math.round(-0.5) 四舍五入。-0.5 四舍五入为 0
40. html5 中新增了 manifest 标签，它有什么作用？
    `应用缓存资源清单`
    manifest 文件是一个简单的文本文件，列举出了浏览器用于离线访问而缓存的资源。
41. 未声明的变量是全局作用域

```JS
(function(){
var a = b = 100;
})();
console.log(typeof a);
console.log(typeof b);
```

undefined，number

42. 代码输出

```JS
var res=null instanceof Object;
console.log(res);
代码输出是？
//JavaScript
false
```

42. 下面哪些是 javascript 中 document 的方法？
    getElementById
    getElementsByTagName
    getElementsByName
43. css 加载
    css 加载不会阻塞 DOM 树的解析
    css 加载会阻塞 DOM 树的渲染
    css 加载会阻塞后面 js 语句的执行
44. 简述 document.write 和 innerHTML 的区别。
    document.write 只能重绘整个页面,
    innerHTML 可以重绘页面的一部分。
45. JavaScript 异步编程的四种方法
    回调函数
    事件监听（eventEmitter)
    发布/订阅
    promise 对象
46. 用户从手机的浏览器访问www.baidu.com，看到的可能跟桌面PC电脑，是不太一样的网页效果，会更适合移动设备使用。请简要分析一下，实现这种网页区分显示的原因及技术原理。
    - 手机的网速问题、屏幕大小、内存、CPU 等。通过不同设备的特征，实现不同的网页展现或输出效果。根据` useragent`、屏幕大小信息、IP、网速、css media Query 等原理，实现前端或后端的特征识别和行为改变。
    - www.baidu.com解析的域名时首先调到www.a.shifen.com, 百度可能在这个服务器上会分析访问浏览器的 `user-agent` 来分析来自桌面设备还是一地哦那个设备, 以此来确定要提供的网页.
47. Flappy Bird 是风靡一时的手机游戏，玩家要操作一只小鸟穿过无穷无尽的由钢管组成的障碍。如果要你在 HTML 前端开发这个游戏，为了保证游戏的流畅运行，并长时间运行也不会崩溃，请列举`开发要注意的性能问题和解决的方法。`

- 长时间运行会`崩溃的原因就是‘内存泄漏’`。我们在日常的 JS 程序中并不太在意内存泄漏问题，因为 JS 解释器会垃圾回收机制，大部分无效内存会被回收，另一方方面 JS 运行在客户端，即使出现内存泄漏也不是太大的问题，简单的刷新页面即可。但是，如果出现要预防内存泄漏的场景还是要注意一些问题。
- 背景的卷轴效果优化。背景不能是无限长的图片拼接，必须有回收已移出的场景的方法。
- 将复杂运算从主 UI 线程中解耦。比如场景中小鸟的运动轨迹、碰撞算法等，需要在空闲时间片运算，不能和 UI 动画同时进行。
- 将比较大的运算分解成不同的时间片，防止阻塞主 UI 线程。最好使用 webworker。
- 注意内存泄漏和回收。使用对象池管理内存，提高内存检测和垃圾回收。
- 进行预处理。将一些常用的过程进行预处理，
- 控制好帧率。将 1 秒分解成多个时间片，在固定间隔时间片进行 UI 动画，其他时间片用在后台运算。
- 通过 GPU 加速和 CSS transition 将小鸟飞行动画和背景动画分离
  需要对减少 `dom 的 reflow`，将动画元素如小鸟设置绝对定位
- 使用 canvas 绘图

48. 在 html 中，下列哪个标签可以创建一个下拉菜单？
    <select><option>baidu</option></select>
49. 设散列表长度为 m，散列函数为 H（key）=key%p，为了减少发生冲突的可能性，p 应取（　`小于m的最大素数`　）
50. 下面程序输出的结果为:
    1、此处的“+”为字符串连接符 2、`020 是八进制`，对应的十进制数字是 16.

```JS
function add(a){
    return a + '010';
}
console.log(add(020));

```

51. enctype 的默认值是 application/x-www-form-urlencoded
    application/x-www-form-urlencoded
    在发送前编码所有字符（默认）

    multipart/form-data
    不对字符编码。
    在使用文件上传的表单时，必须使用该值

    text/plain
    空格转换为“+”加号，但不对特殊字符编码

52. CSS 如何使用服务端的字体？
    @font-face
    @font-face { font-family : name ; src : url( url ) ; sRules }
    说明：
    name : 　字体名称
    url : 　使用绝对或相对地址指定 OpenType 字体
    sRules : 　样式表定义

```CSS

@font-face { font-family: dreamy; font-weight: bold; src: url(http://www.example.com/font.eot); }
```

53. 一个具有 4 个节点的二叉树可能有几种形态?
    公式 C[n,2n] / (n+1)
54. dom 的操作，常用的有哪些，如何创建、添加、移除、移动、复制、查找节点？

```JS
创建：
      createDocumentFragment()    //创建一个DOM片段
      createElement()   //创建一个具体的元素
      createTextNode()   //创建一个文本节点
添加：
    appendChild()
移出：
    removeChild()
替换：
      replaceChild()
插入：
      insertBefore()
复制：
      cloneNode(true)
查找：
      getElementsByTagName()    //通过标签名称
      getElementsByClassName()    //通过标签名称
      getElementsByName()    //通过元素的Name属性的值
      getElementById()    //通过元素Id，唯一性
```

55. 轮播图
    采用复制第一帧和最后一帧的方式来保证在首尾两张图片的无缝滚动切换
56. text-shadow 属性中的四个值 (length、length、length、color) 分别是什么意义：
    阴影离开文字的横方向距离，阴影离开文字的纵方向距离，阴影的模糊半径，阴影的颜色
57. 移动端前端开发与 PC 端比有哪些不同？
    开放题；考察前端开发移动端与 PC 端的熟悉程度；
    回答得分占比：
    （1）浏览器兼容、手机端兼容及一些常见的兼容问题； 30%
    （2）Js 事件处理，click 事件、Touch 事件处理； 30%
    （3）布局方面; 30%
    （4）开发调试; 10%
    （5）其它一些开发技巧。
58. 树组件

    1. 写出页面的 HTML 结构(20%)
    2. 写出页面基本的 css 样式(20%)
    3. 通过递归遍历数据生成树形，写出具体点(20%)
    4. 点击父节点全选/全取消，点击叶节点能选中或取消，写出具体点(20%)
    5. 点击箭头图标展开子项(10%)
    6. 通过`委托绑定事件`(10%)
