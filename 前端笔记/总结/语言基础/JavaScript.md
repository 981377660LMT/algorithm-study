## JavaScript 面试知识点总结

本部分主要是笔者在复习 JavaScript 相关知识和一些相关面试题时所做的笔记，如果出现错误，希望大家指出！

### 目录

- [JavaScript 面试知识点总结](#javascript-面试知识点总结)
  - [目录](#目录)
    - [1. 介绍 js 的基本数据类型。](#1-介绍-js-的基本数据类型)
    - [2. JavaScript 有几种类型的值？你能画一下他们的内存图吗？](#2-javascript-有几种类型的值你能画一下他们的内存图吗)
    - [3. 什么是堆？什么是栈？它们之间有什么区别和联系？](#3-什么是堆什么是栈它们之间有什么区别和联系)
    - [4. 内部属性 [[Class]] 是什么？](#4-内部属性-class-是什么)
    - [5. 介绍 js 有哪些内置对象？](#5-介绍-js-有哪些内置对象)
    - [6. undefined 与 undeclared 的区别？](#6-undefined-与-undeclared-的区别)
    - [7. null 和 undefined 的区别？](#7-null-和-undefined-的区别)
    - [8. 如何获取安全的 undefined 值？](#8-如何获取安全的-undefined-值)
    - [10. JavaScript 原型，原型链？ 有什么特点？](#10-javascript-原型原型链-有什么特点)
    - [11. js 获取原型的方法？](#11-js-获取原型的方法)
    - [12. 在 js 中不同进制数字的表示方式](#12-在-js-中不同进制数字的表示方式)
    - [13. js 中整数的安全范围是多少？](#13-js-中整数的安全范围是多少)
    - [14. typeof NaN 的结果是什么？](#14-typeof-nan-的结果是什么)
    - [15. isNaN 和 Number.isNaN 函数的区别？](#15-isnan-和-numberisnan-函数的区别)
    - [16. Array 构造函数只有一个参数值时的表现？](#16-array-构造函数只有一个参数值时的表现)
    - [17. 其他值到字符串的转换规则？](#17-其他值到字符串的转换规则)
    - [18. 其他值到数字值的转换规则？](#18-其他值到数字值的转换规则)
    - [20. {} 和 [] 的 valueOf 和 toString 的结果是什么？](#20--和--的-valueof-和-tostring-的结果是什么)
    - [21. 什么是假值对象？](#21-什么是假值对象)
    - [22. ~ 操作符的作用？](#22--操作符的作用)
    - [23. 解析字符串中的数字和将字符串强制类型转换为数字的返回结果都是数字，它们之间的区别是什么？](#23-解析字符串中的数字和将字符串强制类型转换为数字的返回结果都是数字它们之间的区别是什么)
    - [24. `+` 操作符什么时候用于字符串的拼接？](#24--操作符什么时候用于字符串的拼接)
    - [25. 什么情况下会发生布尔值的隐式强制类型转换？](#25-什么情况下会发生布尔值的隐式强制类型转换)
    - [27. Symbol 值的强制类型转换？](#27-symbol-值的强制类型转换)
    - [28. == 操作符的强制类型转换规则？](#28--操作符的强制类型转换规则)
    - [29. 如何将字符串转化为数字，例如 '12.3b'?](#29-如何将字符串转化为数字例如-123b)
    - [30. 如何将浮点数点左边的数每三位添加一个逗号，如 12000000.11 转化为『12,000,000.11』?](#30-如何将浮点数点左边的数每三位添加一个逗号如-1200000011-转化为1200000011)
    - [37. Javascript 的作用域链？](#37-javascript-的作用域链)
    - [38. 谈谈 This 对象的理解。](#38-谈谈-this-对象的理解)
    - [39. eval 是做什么的？](#39-eval-是做什么的)
    - [40. 什么是 DOM 和 BOM？](#40-什么是-dom-和-bom)
    - [42. 事件是什么？IE 与火狐的事件机制有什么区别？ 如何阻止冒泡？](#42-事件是什么ie-与火狐的事件机制有什么区别-如何阻止冒泡)
    - [44. 事件委托是什么？](#44-事件委托是什么)
    - [45. ["1", "2", "3"].map(parseInt) 答案是多少？](#45-1-2-3mapparseint-答案是多少)
    - [46. 什么是闭包，为什么要用它？](#46-什么是闭包为什么要用它)
    - [47. javascript 代码中的 "use strict"; 是什么意思 ? 使用它区别是什么？](#47-javascript-代码中的-use-strict-是什么意思--使用它区别是什么)
    - [48. 如何判断一个对象是否属于某个类？](#48-如何判断一个对象是否属于某个类)
    - [51. Javascript 中，有一个函数，执行时对象查找时，永远不会去查找原型，这个函数是？](#51-javascript-中有一个函数执行时对象查找时永远不会去查找原型这个函数是)
    - [52. 对于 JSON 的了解？](#52-对于-json-的了解)
    - [54. js 延迟加载的方式有哪些？](#54-js-延迟加载的方式有哪些)
    - [55. Ajax 是什么? 如何创建一个 Ajax？](#55-ajax-是什么-如何创建一个-ajax)
    - [57. Ajax 解决浏览器缓存问题？](#57-ajax-解决浏览器缓存问题)
    - [58. 同步和异步的区别？](#58-同步和异步的区别)
    - [59. 什么是浏览器的同源政策？](#59-什么是浏览器的同源政策)
    - [60. 如何解决跨域问题？](#60-如何解决跨域问题)
    - [62. 简单谈一下 cookie ？](#62-简单谈一下-cookie-)
    - [64. js 的几种模块规范？](#64-js-的几种模块规范)
    - [66. ES6 模块与 CommonJS 模块、AMD、CMD 的差异。](#66-es6-模块与-commonjs-模块amdcmd-的差异)
    - [67. requireJS 的核心原理是什么？（如何动态加载的？如何避免多次加载的？如何 缓存的？）](#67-requirejs-的核心原理是什么如何动态加载的如何避免多次加载的如何-缓存的)
    - [70. documen.write 和 innerHTML 的区别？](#70-documenwrite-和-innerhtml-的区别)
    - [74. JavaScript 类数组对象的定义？](#74-javascript-类数组对象的定义)
    - [77. [,,,] 的长度？](#77--的长度)
    - [78. JavaScript 中的作用域与变量声明提升？](#78-javascript-中的作用域与变量声明提升)
    - [80. 简单介绍一下 V8 引擎的垃圾回收机制](#80-简单介绍一下-v8-引擎的垃圾回收机制)
    - [81. 哪些操作会造成内存泄漏？](#81-哪些操作会造成内存泄漏)
    - [83. 如何判断当前脚本运行在浏览器还是 node 环境中？（阿里）](#83-如何判断当前脚本运行在浏览器还是-node-环境中阿里)
    - [86. 什么是“前端路由”？什么时候适合使用“前端路由”？“前端路由”有哪些优点和缺点？](#86-什么是前端路由什么时候适合使用前端路由前端路由有哪些优点和缺点)
    - [88. 检测浏览器版本版本有哪些方式？](#88-检测浏览器版本版本有哪些方式)
    - [89. 什么是 Polyfill ？](#89-什么是-polyfill-)
    - [101. toPrecision 和 toFixed 和 Math.round 的区别？](#101-toprecision-和-tofixed-和-mathround-的区别)
    - [102. 什么是 XSS 攻击？如何防范 XSS 攻击？](#102-什么是-xss-攻击如何防范-xss-攻击)
    - [103. 什么是 CSP？](#103-什么是-csp)
    - [105. 什么是 Samesite Cookie 属性？](#105-什么是-samesite-cookie-属性)
    - [110. Object.defineProperty 介绍？](#110-objectdefineproperty-介绍)
    - [111. 使用 Object.defineProperty() 来进行数据劫持有什么缺点？](#111-使用-objectdefineproperty-来进行数据劫持有什么缺点)
    - [112. 什么是 Virtual DOM？为什么 Virtual DOM 比原生 DOM 快？](#112-什么是-virtual-dom为什么-virtual-dom-比原生-dom-快)
    - [113. 如何比较两个 DOM 树的差异？](#113-如何比较两个-dom-树的差异)
    - [116. offsetWidth/offsetHeight,clientWidth/clientHeight 与 scrollWidth/scrollHeight 的区别？](#116-offsetwidthoffsetheightclientwidthclientheight-与-scrollwidthscrollheight-的区别)
    - [117. 谈一谈你理解的函数式编程？](#117-谈一谈你理解的函数式编程)
    - [118. 异步编程的实现方式？](#118-异步编程的实现方式)
    - [120. get 请求传参长度的误区](#120-get-请求传参长度的误区)
    - [122. get 和 post 请求在缓存方面的区别](#122-get-和-post-请求在缓存方面的区别)
    - [123. 图片的懒加载和预加载](#123-图片的懒加载和预加载)
    - [124. mouseover 和 mouseenter 的区别？](#124-mouseover-和-mouseenter-的区别)
    - [125. js 拖拽功能的实现](#125-js-拖拽功能的实现)
    - [127. let 和 const 的注意点？](#127-let-和-const-的注意点)
    - [129. 什么是尾调用，使用尾调用有什么好处？](#129-什么是尾调用使用尾调用有什么好处)
    - [130. Symbol 类型的注意点？](#130-symbol-类型的注意点)
    - [131. Set 和 WeakSet 结构？](#131-set-和-weakset-结构)
    - [134. Reflect 对象创建目的？](#134-reflect-对象创建目的)
    - [135. require 模块引入的查找方式？](#135-require-模块引入的查找方式)
    - [138. 如何检测浏览器所支持的最小字体大小？](#138-如何检测浏览器所支持的最小字体大小)
    - [140. 单例模式模式是什么？](#140-单例模式模式是什么)
    - [141. 策略模式是什么？ 查表](#141-策略模式是什么-查表)
    - [142. 代理模式是什么？](#142-代理模式是什么)
    - [143. 中介者模式是什么？](#143-中介者模式是什么)
    - [144. 适配器模式是什么？](#144-适配器模式是什么)
    - [145. 观察者模式和发布订阅模式有什么不同？](#145-观察者模式和发布订阅模式有什么不同)
    - [147. Vue 的各个生命阶段是什么？](#147-vue-的各个生命阶段是什么)
    - [150. vue-router 中的导航钩子函数](#150-vue-router-中的导航钩子函数)
    - [151. $route 和 $router 的区别？](#151-route-和-router-的区别)
    - [152. vue 常用的修饰符？](#152-vue-常用的修饰符)
    - [157. 开发中常用的几种 Content-Type ？](#157-开发中常用的几种-content-type-)
    - [165. 如何确定页面的可用性时间，什么是 Performance API？](#165-如何确定页面的可用性时间什么是-performance-api)
    - [167. js 语句末尾分号是否可以省略？](#167-js-语句末尾分号是否可以省略)
    - [168. Object.assign()](#168-objectassign)
    - [170. js for 循环注意点](#170-js-for-循环注意点)
    - [171. 一个列表，假设有 100000 个数据，这个该怎么办？](#171-一个列表假设有-100000-个数据这个该怎么办)
    - [172. js 中倒计时的纠偏实现？](#172-js-中倒计时的纠偏实现)

#### 1. 介绍 js 的基本数据类型。

js 一共有 7 种基本数据类型，分别是 Undefined、Null、Boolean、Number、String，还有在 ES6 中新增的 Symbol 和 ES10 中新增的 BigInt 类型。
Symbol 代表创建后独一无二且不可变的数据类型，它的出现我认为主要是为了`解决可能出现的全局变量冲突的问题`/`属性的身份标识符例如NONE`/`类的标识符`/`魔法方法`。
BigInt 是一种数字类型的数据，它可以表示任意精度格式的整数，使用 BigInt 可以安全地存储和操作大整数，即使这个数已经超出了 Number 能够表示的安全整数范围。

#### 2. JavaScript 有几种类型的值？你能画一下他们的内存图吗？

涉及知识点：

- 栈：原始数据类型（Undefined、Null、Boolean、Number、String）
- 堆：引用数据类型（对象、数组和函数）

两种类型的区别是：存储位置不同。
原始数据类型直接存储在栈（stack）中的简单数据段，占据空间小、大小固定，属于被频繁使用数据，所以放入栈中存储。

引用数据类型存储在堆（heap）中的对象，占据空间大、大小不固定。如果存储在栈中，将会影响程序运行的性能；`引用数据类型在栈中存储了指针`，该指针指向堆中该实体的起始地址。当解释器寻找引用值时，会首先检索其在栈中的地址，取得地址后从堆中获得实体。

回答：

js 可以分为两种类型的值，一种是基本数据类型，一种是复杂数据类型。

基本数据类型....（参考 1）

复杂数据类型指的是 Object 类型，所有其他的如 Array、Date 等数据类型都可以理解为 Object 类型的子类。

两种类型间的主要区别是它们的存储位置不同，基本数据类型的值直接保存在栈中，而复杂数据类型的值保存在堆中，通过使用在栈中保存对应的指针来获取堆中的值。

详细资料可以参考：
[《JavaScript 有几种类型的值？》](https://blog.csdn.net/lxcao/article/details/52749421)
[《JavaScript 有几种类型的值？能否画一下它们的内存图；》](https://blog.csdn.net/jiangjuanjaun/article/details/80327342)

#### 3. 什么是堆？什么是栈？它们之间有什么区别和联系？

在操作系统中，内存被分为栈区和堆区。

`栈区内存由编译器自动分配释放`，存放函数的参数值，局部变量的值等。其操作方式类似于数据结构中的栈。

堆区内存一般由程序员分配释放，若程序员不释放，程序结束时`可能由垃圾回收机制回收`。闭包存在这里。

详细资料可以参考：
[《什么是堆？什么是栈？他们之间有什么区别和联系？》](https://www.zhihu.com/question/19729973)

#### 4. 内部属性 [[Class]] 是什么？

所有 typeof 返回值为 "object" 的对象（如数组）都包含一个内部属性 [[Class]]（我们可以把它看作一个内部的分类，而非
传统的面向对象意义上的类）。这个属性无法直接访问，一般通过 Object.prototype.toString(..) 来查看。例如：

```JS

Object.prototype.toString.call( [1,2,3] );
// "[object Array]"

Object.prototype.toString.call( /regex-literal/i );
// "[object RegExp]"

// 我们自己创建的类就不会有这份特殊待遇，因为 toString() 找不到 toStringTag 属性时只好返回默认的 Object 标签
// 默认情况类的[[Class]]返回[object Object]
class Class1 {}
Object.prototype.toString.call(new Class1()); // "[object Object]"
// 需要定制[[Class]]
class Class2 {
get [Symbol.toStringTag]() {
return "Class2";
}
}
Object.prototype.toString.call(new Class2()); // "[object Class2]"

```

#### 5. 介绍 js 有哪些内置对象？

涉及知识点：

全局的对象（ global objects ）或称标准内置对象，不要和 "全局对象（global object）" 混淆。这里说的全局的对象是说在
全局作用域里的对象。全局作用域中的其他对象可以由用户的脚本创建或由宿主程序提供。

标准内置对象的分类

（1）值属性，这些全局属性返回一个简单值，这些值没有自己的属性和方法。

例如 Infinity、NaN、undefined、null 字面量

（2）函数属性，全局函数可以直接调用，不需要在调用时指定所属对象，执行结束后会将结果直接返回给调用者。

例如 eval()、parseFloat()、parseInt() 等

（3）基本对象，基本对象是定义或使用其他对象的基础。基本对象包括一般对象、函数对象和错误对象。

例如 Object、Function、Boolean、Symbol、Error 等

（4）数字和日期对象，用来表示数字、日期和执行数学计算的对象。

例如 Number、Math、Date

（5）字符串，用来表示和操作字符串的对象。

例如 String、RegExp

（6）可索引的集合对象，这些对象表示按照索引值来排序的数据集合，包括数组和类型数组，以及类数组结构的对象。例如 Array

（7）使用键的集合对象，这些集合对象在存储数据时会使用到键，支持按照插入顺序来迭代元素。

例如 Map、Set、WeakMap、WeakSet

（8）矢量集合，SIMD 矢量集合中的数据会被组织为一个数据序列。

例如 SIMD 等

（9）结构化数据，这些对象用来表示和操作结构化的缓冲区数据，或使用 JSON 编码的数据。

例如 JSON 等

（10）控制抽象对象

例如 Promise、Generator 等

（11）反射

例如 Reflect、Proxy

（12）国际化，为了支持多语言处理而加入 ECMAScript 的对象。

例如 Intl、Intl.Collator 等

（13）WebAssembly

（14）其他

例如 arguments

回答：

js 中的内置对象主要指的是在程序执行前存在全局作用域里的由 js 定义的一些全局值属性、函数和用来实例化其他对象的构造函数对象。一般我们经常用到的如全局变量值 NaN、undefined，全局函数如 parseInt()、parseFloat() 用来实例化对象的构造函数如 Date、Object 等，还有提供数学计算的单体内置对象如 Math 对象。

详细资料可以参考：
[《标准内置对象的分类》](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects)
[《JS 所有内置对象属性和方法汇总》](https://segmentfault.com/a/1190000011467723#articleHeader24)

#### 6. undefined 与 undeclared 的区别？

已在作用域中声明但还没有赋值的变量，是 undefined 的。相反，还没有在作用域中声明过的变量，是 undeclared 的。

对于 undeclared 变量的引用，浏览器会报引用错误，如 ReferenceError: b is not defined 。`但是我们可以使用 typeof 的安全防范机制来避免报错，因为对于 undeclared（或者 not defined ）变量，typeof 会返回 "undefined"。`

#### 7. null 和 undefined 的区别？

首先 Undefined 和 Null 都是基本数据类型，这两个基本数据类型分别都只有一个值，就是 undefined 和 null。

undefined 代表的含义是`未定义`，null 代表的含义是`空对象`。一般变量声明了但`还没有定义`的时候会返回 undefined，null 主要用于赋值给一些可能会返回对象的变量，作为`初始化`。

undefined 在 js 中`不是一个保留字`，这意味着我们可以`使用 undefined 来作为一个变量名`，这样的做法是非常危险的，它会影响我们对 undefined 值的判断。但是我们可以`通过一些方法获得安全的 undefined 值，比如说 void 0。`

当我们对两种类型使用 typeof 进行判断的时候，Null 类型化会返回 “object”，这是一个历史遗留的问题。当我们使用双等对两种类型的值进行比较时会返回 true，使用三个等号时会返回 false。

详细资料可以参考：
[《JavaScript 深入理解之 undefined 与 null》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8Bundefined%E4%B8%8Enull.html)

#### 8. 如何获取安全的 undefined 值？

因为 undefined 是一个标识符，所以可以被当作变量来使用和赋值，但是这样会影响 undefined 的正常判断。

表达式 void \_\_\_ 没有返回值，因此返回结果是 undefined。void 并不改变表达式的结果，只是让表达式不返回值。

按惯例我们用 void 0 来获得 undefined。

#### 10. JavaScript 原型，原型链？ 有什么特点？

在 js 中我们是使用构造函数来新建一个对象的，`每一个构造函数的内部都有一个 prototype 属性值，这个属性值是一个对象`，这个对象包含了可以由该构造函数的所有实例共享的属性和方法。当我们使用构造函数新建一个对象后，在这个对象的内部将包含一个`指针，这个指针指向构造函数的 prototype 属性对应的值`，在 ES5 中这个指针被称为对象的原型。一般来说我们是不应该能够获取到这个值的，但是现在浏览器中都实现了 `__proto__` 属性来让我们访问这个属性，但是我们最好不要使用这个属性，因为它不是规范中规定的。ES5 中新增了一个 Object.getPrototypeOf() 方法，我们可以通过这个方法来获取对象的原型。

**当我们访问一个对象的属性时，如果这个对象内部不存在这个属性，那么它就会去它的原型对象里找这个属性**，这个原型对象又会有自己的原型，于是就这样一直找下去，也就是原型链的概念。`原型链的尽头一般来说都是 Object.prototype` `所以这就是我们新建的对象为什么能够使用 toString() 等方法的原因。`

特点：

JavaScript 对象是通过引用来传递的，我们创建的每个新对象实体中并没有一份属于自己的原型副本。当我们修改原型时，与之相关的对象也会继承这一改变。

详细资料可以参考：
[《JavaScript 深入理解之原型与原型链》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8B%E5%8E%9F%E5%9E%8B%E4%B8%8E%E5%8E%9F%E5%9E%8B%E9%93%BE.html)

#### 11. js 获取原型的方法？

- p.\_\_proto\_\_
- p.constructor.prototype
- `Object.getPrototypeOf(p)`

#### 12. 在 js 中不同进制数字的表示方式

- 以 0X、0x 开头的表示为十六进制。

- 以 0、0O、0o 开头的表示为八进制。

- 以 0B、0b 开头的表示为二进制格式。

#### 13. js 中整数的安全范围是多少？

安全整数指的是，在这个范围内的整数`转化为二进制存储的时候不会出现精度丢失`，能够被“安全”呈现的最大整数是 2^53 - 1，
即 9007199254740991，在 ES6 中被定义为 Number.MAX_SAFE_INTEGER。最小整数是-9007199254740991，在 ES6 中被定义为 Number.MIN_SAFE_INTEGER。

如果某次计算的结果得到了一个超过 JavaScript 数值范围的值，那么这个值会被自动转换为特殊的 Infinity 值。如果某次计算返回了正或负的 Infinity 值，那么该值将无法参与下一次的计算。判断一个数是不是有穷的，可以使用 Number.isFinite 函数来判断。

#### 14. typeof NaN 的结果是什么？

NaN 意指“不是一个数字”（not a number），NaN 是一个“警戒值”（sentinel value，有特殊用途的常规值），用于指出
数字类型中的错误情况，即“执行数学运算没有成功，这是失败后返回的结果”。

**typeof NaN; // "number"**

NaN 是一个特殊值，它和自身不相等，是唯一一个非自反（自反，reflexive，即 x === x 不成立）的值。而 NaN != NaN
为 true。

#### 15. isNaN 和 Number.isNaN 函数的区别？

函数 isNaN 接收参数后，会尝试将这个`参数转换为数值`，任何不能被转换为数值的的值都会返回 true，因此非数字值传入也会返回 true ，会影响 NaN 的判断。

**函数 Number.isNaN 会首先判断传入参数是否为数字**，如果是数字再继续判断是否为 NaN

#### 16. Array 构造函数只有一个参数值时的表现？

Array 构造函数只带一个数字参数的时候，该参数会被作为数组的预设长度（length），而非只充当数组中的一个元素。这样
创建出来的只是一个空数组，只不过它的 length 属性被设置成了指定的值。
构造函数 Array(..) 不要求必须带 new 关键字。不带时，它会被自动补上。

#### 17. 其他值到字符串的转换规则？

规范的 9.8 节中定义了抽象操作 toString ，它负责处理非字符串到字符串的强制类型转换。
（1）Null 和 Undefined 类型 ，null 转换为 "null"，undefined 转换为 "undefined"，
（2）Boolean 类型，true 转换为 "true"，false 转换为 "false"。
（3）Number 类型的值直接转换，不过那些极小和极大的数字会使用指数形式。
（4）Symbol 类型的值直接转换，但是只允许显式强制类型转换，使用隐式强制类型转换会产生错误。
（5）对普通对象来说，除非自行定义 toString() 方法，否则会调用 toString()（Object.prototype.toString()）

来返回内部属性 [[Class]] 的值，如"[object Object]"。如果对象有自己的 toString() 方法，字符串化时就会
调用该方法并使用其返回值。

#### 18. 其他值到数字值的转换规则？

有时我们需要将非数字值当作数字来使用，比如数学运算。为此 ES5 规范在 9.3 节定义了抽象操作 toNumber。

（1）`Undefined 类型的值转换为 NaN`。
（2）`Null 类型的值转换为 0`。
（3）Boolean 类型的值，true 转换为 1，false 转换为 0。
（4）String 类型的值转换如同使用 Number() 函数进行转换，`如果包含非数字值则转换为 NaN，空字符串为 0。`
（5）`Symbol 类型的值不能转换为数字，会报错。`
（6）对象（包括数组）会首先被转换为相应的基本类型值，如果返回的是非数字的基本类型值，则再遵循以上规则将其强制转换为数字。

为了将值转换为相应的基本类型值，抽象操作 ToPrimitive 会首先（通过内部操作 DefaultValue）检查该值是否有 valueOf() 方法。如果有并且返回基本类型值，就使用该值进行强制类型转换。如果没有就使用 toString() 的返回值（如果存在）来进行强制类型转换。

`如果 valueOf() 和 toString() 均不返回基本类型值，会产生 TypeError 错误。`

#### 20. {} 和 [] 的 valueOf 和 toString 的结果是什么？

{} 的 valueOf 结果为 {} ，`toString 的结果为 "[object Object]"`

[] 的 valueOf 结果为 [] ，`toString 的结果为 ""`

#### 21. 什么是假值对象？

浏览器在某些特定情况下，在常规 JavaScript 语法基础上自己创建了一些外来值，这些就是“假值对象”。
假值对象看起来和普通对象并无二致（都有属性，等等），但将它们强制类型转换为布尔值时结果为 false
**最常见的例子是 document.all**，它是一个类数组对象，包含了页面上的所有元素，由 DOM（而不是 JavaScript 引擎）提供给 JavaScript 程序使用。

#### 22. ~ 操作符的作用？

~ 返回 补码，并且 ~ 会将数字`转换为 32 位整数`，因此我们可以使用 ~ 来进行取整操作。
~x 大致等同于 -(x+1)。

#### 23. 解析字符串中的数字和将字符串强制类型转换为数字的返回结果都是数字，它们之间的区别是什么？

解析允许字符串（如 parseInt() ）中含有非数字字符，解析按从左到右的顺序，如果遇到非`数字字符就停止`。
而转换（如 Number ()）`不允许出现非数字字符，否则会失败并返回 NaN`。

#### 24. `+` 操作符什么时候用于字符串的拼接？

简单来说就是，如果 + 的其中一个操作数是字符串（或者通过以上步骤最终得到字符串），则执行字符串拼接，否则执行数字加法。

那么对于除了加法的运算符来说，只要其中一方是数字，那么另一方就会被转为数字。

#### 25. 什么情况下会发生布尔值的隐式强制类型转换？

（1） if (..) 语句中的条件判断表达式。
（2） for ( .. ; .. ; .. ) 语句中的条件判断表达式（第二个）。
（3） while (..) 和 do..while(..) 循环中的条件判断表达式。
（4） ? : 中的条件判断表达式。
（5） 逻辑运算符 ||（逻辑或）和 &&（逻辑与）左边的操作数（作为条件判断表达式）。

#### 27. Symbol 值的强制类型转换？

ES6 `允许从 Symbol 到字符串的显式强制类型转换，然而隐式强制类型转换会产生错误。`
Symbol 值不能够被强制类型转换为数字（显式和隐式都会产生错误），但可以被强制类型转换为布尔值（显式和隐式结果都是 true ）。

#### 28. == 操作符的强制类型转换规则？

（1）字符串和数字之间的相等比较，将`字符串转换为数字`之后再进行比较。
（2）其他类型和布尔类型之间的相等比较，先将`布尔值转换为数字`后，再应用其他规则进行比较。
（3）null 和 undefined 之间的相等比较，结果为真。其他值和它们进行比较都返回假值。
（4）对象和非对象之间的相等比较，对象先调用 `ToPrimitive` 抽象操作后，再进行比较。
（5）如果一个操作值为 NaN ，则相等比较返回 false（ NaN 本身也不等于 NaN ）。
（6）如果两个操作值都是对象，则比较它们是不是指向同一个对象。如果两个操作数都指向同一个对象，则相等操作符返回 true，否则，返回 false。

详细资料可以参考：
[《JavaScript 字符串间的比较》](https://www.jeffjade.com/2015/08/28/2015-09-02-js-string-compare/)

#### 29. 如何将字符串转化为数字，例如 '12.3b'?

（1）使用 Number() 方法，前提是所包含的字符串不包含不合法字符。
（2）使用 parseInt() 方法，parseInt() 函数可解析一个字符串，并返回一个整数。还可以设置要解析的数字的基数。当基数的值为 0，或没有设置该参数时，parseInt() 会根据 string 来判断数字的基数。
（3）使用 parseFloat() 方法，该函数解析一个字符串参数并返回一个浮点数。
（4）使用 + 操作符的隐式转换，前提是所包含的字符串不包含不合法字符。

详细资料可以参考：
[《详解 JS 中 Number()、parseInt() 和 parseFloat() 的区别》](https://blog.csdn.net/m0_38099607/article/details/72638678)

#### 30. 如何将浮点数点左边的数每三位添加一个逗号，如 12000000.11 转化为『12,000,000.11』?

```JS
js
// 方法一
function format(number) {
return number && number.replace(/(?!^)(?=(\d{3})+\.)/g, ",");
}
// 方法二
function format1(number) {
return Intl.NumberFormat().format(number)
}
// 方法三
function format2(number) {
return number.toLocaleString('en')
}
```

#### 37. Javascript 的作用域链？

作用域链的作用是保证对执行环境有权访问的所有变量和函数的有序访问，`通过作用域链，我们可以访问到外层环境的变量和函数。`

`作用域链的本质上是一个指向变量对象的指针列表`。变量对象是一个包含了执行环境中所有变量和函数的对象。
`作用域链的前端始终都是当前执行上下文的变量对象。全局执行上下文的变量对象（也就是全局对象）始终是作用域链的最后一个对象。`

当我们查找一个变量时，`如果当前执行环境中没有找到，我们可以沿着作用域链向后查找。`

作用域链的创建过程跟执行上下文的建立有关....

详细资料可以参考：
[《JavaScript 深入理解之作用域链》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8B%E4%BD%9C%E7%94%A8%E5%9F%9F%E9%93%BE.html)

<!-- TODO -->

#### 38. 谈谈 This 对象的理解。

this 是执行上下文中的一个属性，`它指向最后一次调用这个方法的对象`。在实际开发中，this 的指向可以通过四种调用模式来判断。

- 1.第一种是函数调用模式，当一个函数不是一个对象的属性时，直接作为函数来调用时，this 指向`全局`对象。

- 2.第二种是方法调用模式，如果一个函数作为一个对象的方法来调用时，this 指向这个`对象`。

- 3.第三种是构造器调用模式，如果一个函数用 new 调用时，函数执行前会新创建一个对象，`this 指向这个新创建的对象`。

- 4.第四种是 apply 、 call 和 bind 调用模式，这三个方法都可以显示的指定调用函数的 this 指向。其中 apply 方法接收两个参数：一个是 this 绑定的对象，一个是参数数组。call 方法接收的参数，第一个是 this 绑定的对象，后面的其余参数是传入函数执行的参数。也就是说，在使用 call() 方法时，传递给函数的参数必须逐个列举出来。bind 方法通过传入一个对象，返回一个 this 绑定了传入对象的新函数。这个函数的 this 指向除了使用 new 时会被改变，其他情况下都不会改变。

`这四种方式，使用构造器调用模式的优先级最高，然后是 apply 、 call 和 bind 调用模式，然后是方法调用模式，然后是函数调用模式。`

[《JavaScript 深入理解之 this 详解》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8Bthis%E8%AF%A6%E8%A7%A3.html)

#### 39. eval 是做什么的？

**用过 解析字符串成对象的时候**

它的功能是把`对应的字符串解析成 JS 代码并运行`。

应该避免使用 eval，不安全，非常耗性能（2 次，`一次解析成 js 语句，一次执行`）。

详细资料可以参考：
[《eval()》](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/eval)

#### 40. 什么是 DOM 和 BOM？

`DOM 指的是文档对象模型`，它指的是把`文档当做一个对象来对待`，这个对象主要定义了`处理网页内容的方法和接口`。

`BOM 指的是浏览器对象模型`，它指的是把`浏览器当做一个对象来对待`，这个对象主要定义了与浏览器进行交互的法和接口。BOM 的核心是 window，而 window 对象具有双重角色，它既是通过 js 访问浏览器窗口的一个接口，又是一个 Global（全局）对象。这意味着在网页中定义的任何对象，变量和函数，都作为全局对象的一个属性或者方法存在。window 对象含有 location 对象、navigator 对象、screen 对象等子对象，并且 DOM 的最根本的对象 document 对象也是 BOM 的 window 对象的子对象。

详细资料可以参考：
[《DOM, DOCUMENT, BOM, WINDOW 有什么区别?》](https://www.zhihu.com/question/33453164)
[《Window 对象》](http://www.w3school.com.cn/jsref/dom_obj_window.asp)
[《DOM 与 BOM 分别是什么，有何关联？》](https://www.zhihu.com/question/20724662)
[《JavaScript 学习总结（三）BOM 和 DOM 详解》](https://segmentfault.com/a/1190000000654274#articleHeader21)

#### 42. 事件是什么？IE 与火狐的事件机制有什么区别？ 如何阻止冒泡？

- 1.事件是用户操作网页时发生的交互动作，比如 click/move， 事件除了用户触发的动作外，还可以是文档加载，窗口滚动和大小调整。事件被封装成一个 event 对象，包含了该事件发生时的所有相关信息（ event 的属性）以及可以对事件进行的操作（ event 的方法）。

- 2.事件处理机制：`IE 支持事件冒泡、Firefox 同时支持两种事件模型，也就是：事件冒泡和事件捕获`。

- 3.event.stopPropagation() 或者 ie 下的方法 event.cancelBubble = true;

详细资料可以参考：
[《Javascript 事件模型系列（一）事件及事件的三种模型》](https://www.cnblogs.com/lvdabao/p/3265870.html)
[《Javascript 事件模型：事件捕获和事件冒泡》](https://blog.csdn.net/wuseyukui/article/details/13771493)

#### 44. 事件委托是什么？

事件委托本质上是利用了浏览器事件冒泡的机制。因为事件在冒泡过程中会上传到父节点，并且父节点可以通过事件对象获取到
目标节点，因此可以把`子节点的监听函数定义在父节点上`，由父节点的监听函数统一处理多个子元素的事件，这种方式称为事件代理。

使用事件代理我们可以不必要为每一个子元素都绑定一个监听事件，这样减少了内存上的消耗。并且使用事件代理我们还可以实现事件的`动态绑定`，比如说新增了一个子节点，我们并不需要单独地为它添加一个监听事件，它所发生的事件会交给父元素中的监听函数来处理。

详细资料可以参考：
[《JavaScript 事件委托详解》](https://zhuanlan.zhihu.com/p/26536815)

#### 45. ["1", "2", "3"].map(parseInt) 答案是多少？

parseInt() 函数能解析一个字符串，并返回一个整数，需要两个参数 (val, radix)，其中 radix 表示要解析的数字的基数。（该值介于 2 ~ 36 之间，并且字符串中的数字不能大于 radix 才能正确返回数字结果值）。

此处 map 传了 3 个参数 (element, index, array)，默认第三个参数被忽略掉，因此三次传入的参数分别为 "1-0", "2-1", "3-2"

因为字符串的值不能大于基数，因此后面两次调用均失败，返回 NaN ，`第一次基数为 0 ，按十进制解析返回 1。`

详细资料可以参考：
[《为什么 ["1", "2", "3"].map(parseInt) 返回 [1,NaN,NaN]？》](https://blog.csdn.net/justjavac/article/details/19473199)

#### 46. 什么是闭包，为什么要用它？

闭包是指`有权访问另一个函数作用域中变量的函数`，创建闭包的最常见的方式就是在一个函数内创建另一个函数，创建的函数可以访问到当前函数的局部变量。

闭包有两个常用的用途。

闭包的第一个用途是使我们在`函数外部能够访问到函数内部的变量`。通过使用闭包，我们可以通过在外部调用闭包函数，从而在外部访问到函数内部的变量，可以使用这种方法来创建私有变量。

函数的另一个用途是使已经运行`结束的函数上下文中的变量对象继续留在内存中`，因为闭包函数保留了这个变量对象的引用，所以这个变量对象不会被回收。

其实闭包的`本质就是作用域链的一个特殊的应用`，只要了解了作用域链的创建过程，就能够理解闭包的实现原理。

详细资料可以参考：
[《JavaScript 深入理解之闭包》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8B%E9%97%AD%E5%8C%85.html)

#### 47. javascript 代码中的 "use strict"; 是什么意思 ? 使用它区别是什么？

相关知识点：

use strict 是一种 ECMAscript5 添加的（严格）运行模式，这种模式使得 Javascript 在更严格的条件下运行。

设立"严格模式"的目的，主要有以下几个：

- 消除 Javascript 语法的一些不合理、不严谨之处，减少一些怪异行为;
- 消除代码运行的一些不安全之处，保证代码运行的安全；
- 提高编译器效率，增加运行速度；
- 为未来新版本的 Javascript 做好铺垫。

区别：

- 1.`禁止使用 with 语句`。
- 2.`禁止 this 关键字指向全局对象`。
- 3.`对象不能有重名的属性`。

回答：

use strict 指的是严格运行模式，在这种模式对 js 的使用添加了一些限制。比如说禁止 this 指向全局对象，还有禁止使
用 with 语句等。设立严格模式的目的，主要是为了消除代码使用中的一些不安全的使用方式，也是为了消除 js 语法本身的一
些不合理的地方，以此来减少一些运行时的怪异的行为。同时使用严格运行模式也能够提高编译的效率，从而提高代码的运行速度。
我认为严格模式代表了 js 一种更合理、更安全、更严谨的发展方向。

详细资料可以参考：
[《Javascript 严格模式详解》](http://www.ruanyifeng.com/blog/2013/01/javascript_strict_mode.html)

#### 48. 如何判断一个对象是否属于某个类？

第一种方式是使用 `instanceof` 运算符来判断构造函数的 prototype 属性是否出现在对象的原型链中的任何位置。

第二种方式可以通过对象的 constructor 属性来判断，对象的 constructor 属性指向该对象的构造函数，但是这种方式不是很安全，因为 constructor 属性可以被改写。

第三种方式，如果需要判断的是某个内置的引用类型的话，可以使用 `Object.prototype.toString()` 方法来打印对象的
[[Class]] 属性来进行判断。

详细资料可以参考：
[《js 判断一个对象是否属于某一类》](https://blog.csdn.net/haitunmin/article/details/78418522)

#### 51. Javascript 中，有一个函数，执行时对象查找时，永远不会去查找原型，这个函数是？

**hasOwnProperty**

所有继承了 Object 的对象都会继承到 hasOwnProperty 方法。这个方法可以用来检测一个对象是否含有特定的自身属性，和 in 运算符不同，该方法会忽略掉那些从原型链上继承到的属性。

详细资料可以参考：
[《Object.prototype.hasOwnProperty()》](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Object/hasOwnProperty)

#### 52. 对于 JSON 的了解？

回答：

JSON 是一种基于文本的轻量级的数据交换格式。它可以被任何的编程语言读取和作为数据格式来传递。

在项目开发中，我们使用 JSON 作为前后端数据交换的方式。在前端我们通过将一个符合 JSON 格式的数据结构序列化为 JSON 字符串，然后将它传递到后端，后端通过 JSON 格式的字符串解析后生成对应的数据结构，以此来实现前后端数据的一个传递。

因为 JSON 的语法是基于 js 的，因此很容易将 JSON 和 js 中的对象弄混，但是我们应该注意的是 JSON 和 js 中的对象不是一回事，JSON 中对象格式更加严格，比如说在 JSON 中属性值不能为函数，不能出现 NaN 这样的属性值等，因此大多数的 js 对象是不符合 JSON 对象的格式的。

在 js 中提供了两个函数来实现 js 数据结构和 JSON 格式的转换处理，一个是 JSON.stringify 函数，通过传入一个符合 JSON 格式的数据结构，将其转换为一个 JSON 字符串。如果传入的数据结构不符合 JSON 格式，那么在序列化的时候会对这些值进行对应的特殊处理，使其符合规范。在前端向后端发送数据时，我们可以调用这个函数将数据对象转化为 JSON 格式的字符串。

另一个函数 JSON.parse() 函数，这个函数用来将 JSON 格式的字符串转换为一个 js 数据结构，如果传入的字符串不是标准的 JSON 格式的字符串的话，将会抛出错误。当我们从后端接收到 JSON 格式的字符串时，我们可以通过这个方法来将其解析为一个 js 数据结构，以此来进行数据的访问。

详细资料可以参考：
[《深入了解 JavaScript 中的 JSON 》](https://my.oschina.net/u/3284240/blog/874368)

#### 54. js 延迟加载的方式有哪些？

相关知识点：

js 延迟加载，也就是等页面加载完成之后再加载 JavaScript 文件。 js 延迟加载有助于提高页面加载速度。

一般有以下几种方式：

- defer 属性
- async 属性
- 动态创建 DOM 方式
- 使用 setTimeout 延迟方法
- 让 JS 最后加载

回答：

js 的加载、解析和执行会阻塞页面的渲染过程，因此我们希望 js 脚本能够尽可能的延迟加载，提高页面的渲染速度。

我了解到的几种方式是：

第一种方式是我们一般采用的是将 js 脚本放在文档的底部，来使 js 脚本尽可能的在最后来加载执行。

第二种方式是给 js 脚本添加 defer 属性，这个属性会让脚本的加载与文档的解析同步解析，然后在文档解析完成后再执行这个脚本文件，这样的话就能使页面的渲染不被阻塞。多个设置了 defer 属性的脚本按规范来说最后是`顺序执行`的，但是在一些浏览器中可能不是这样。

第三种方式是给 js 脚本添加 async 属性，这个属性会使脚本异步加载，不会阻塞页面的解析过程，但是当脚`本加载完成后立即执行 js 脚本`，这个时候如果文档没有解析完成的话同样会阻塞。多个 async 属性的脚本的执行顺序是不可预测的，`一般不会按照代码的顺序依次执行`。

第四种方式是动态创建 DOM 标签的方式，我们可以对文档的加载事件进行监听，**当文档加载完成后再动态的创建 script 标签来引入 js 脚本。**

详细资料可以参考：
[《JS 延迟加载的几种方式》](https://blog.csdn.net/meijory/article/details/76389762)
[《HTML 5 `<script>` `async` 属性》](http://www.w3school.com.cn/html5/att_script_async.asp)

#### 55. Ajax 是什么? 如何创建一个 Ajax？

相关知识点：

2005 年 2 月，AJAX 这个词第一次正式提出，它是 Asynchronous JavaScript and XML 的缩写，指的是通过 JavaScript 的
异步通信，从服务器获取 XML 文档从中提取数据，再更新当前网页的对应部分，而不用刷新整个网页。

具体来说，AJAX 包括以下几个步骤。

- 1.创建 XMLHttpRequest 对象，也就是创建一个异步调用对象
- 2.创建一个新的 HTTP 请求，并指定该 HTTP 请求的方法、URL 及验证信息
- 3.设置响应 HTTP 请求状态变化的函数
- 4.发送 HTTP 请求
- 5.获取异步调用返回的数据
- 6.使用 JavaScript 和 DOM 实现局部刷新

一般实现：

js
const SERVER_URL = "/server";

let xhr = new XMLHttpRequest();

// 创建 Http 请求
xhr.open("GET", SERVER_URL, true);

// 设置状态监听函数
xhr.onreadystatechange = function() {
if (this.readyState !== 4) return;

// 当请求成功时
if (this.status === 200) {
handle(this.response);
} else {
console.error(this.statusText);
}
};

// 设置请求失败时的监听函数
xhr.onerror = function() {
console.error(this.statusText);
};

// 设置请求头信息
xhr.responseType = "json";
xhr.setRequestHeader("Accept", "application/json");

// 发送 Http 请求
xhr.send(null);

// promise 封装实现：

function getJSON(url) {
// 创建一个 promise 对象
let promise = new Promise(function(resolve, reject) {
let xhr = new XMLHttpRequest();

    // 新建一个 http 请求
    xhr.open("GET", url, true);

    // 设置状态的监听函数
    xhr.onreadystatechange = function() {
      if (this.readyState !== 4) return;

      // 当请求成功或失败时，改变 promise 的状态
      if (this.status === 200) {
        resolve(this.response);
      } else {
        reject(new Error(this.statusText));
      }
    };

    // 设置错误监听函数
    xhr.onerror = function() {
      reject(new Error(this.statusText));
    };

    // 设置响应的数据类型
    xhr.responseType = "json";

    // 设置请求头信息
    xhr.setRequestHeader("Accept", "application/json");

    // 发送 http 请求
    xhr.send(null);

});

return promise;
}

回答：

我对 ajax 的理解是，`它是一种异步通信的方法，通过直接由 js 脚本向服务器发起 http 通信`，然后根据服务器返回的数据，更新网页的相应部分，而不用刷新整个页面的一种方法。

创建一个 ajax 有这样几个步骤

首先是创建一个 XMLHttpRequest 对象。

然后在这个对象上使用 open 方法创建一个 http 请求，open 方法所需要的参数是请求的方法、请求的地址、是否异步和用户的认证信息。

在发起请求前，我们可以为这个对象添加一些信息和监听函数。比如说我们可以通过 setRequestHeader 方法来为请求添加头信息。我们还可以为这个对象添加一个状态监听函数。一个 XMLHttpRequest 对象一共有 5 个状态，当它的状态变化时会触发 onreadystatechange 事件，我们可以通过设置监听函数，来处理请求成功后的结果。当对象的 readyState 变为 4 的时候，代表服务器返回的数据接收完成，这个时候我们可以通过判断请求的状态，如果状态是 2xx 或者 304 的话则代表返回正常。这个时候我们就可以通过 response 中的数据来对页面进行更新了。

当对象的属性和监听函数设置完成后，最后我们调用 sent 方法来向服务器发起请求，可以传入参数作为发送的数据体。

详细资料可以参考：
[《XMLHttpRequest 对象》](https://wangdoc.com/javascript/bom/xmlhttprequest.html)
[《从 ajax 到 fetch、axios》](https://juejin.im/post/5acde23c5188255cb32e7e76)
[《Fetch 入门》](https://juejin.im/post/5c160937f265da61180199b2)
[《传统 Ajax 已死，Fetch 永生》](https://segmentfault.com/a/1190000003810652)

#### 57. Ajax 解决浏览器缓存问题？

- 1.在 ajax 发送请求前加上 anyAjaxObj.setRequestHeader("If-Modified-Since","0")。

- 2.在 ajax 发送请求前加上 anyAjaxObj.setRequestHeader("Cache-Control","no-cache")。

- 3.在 URL 后面加上一个随机数： "fresh=" + Math.random();。

- 4.在 URL 后面加上时间戳："nowtime=" + new Date().getTime();。

- 5.如果是使用 jQuery，直接这样就可以了\$.ajaxSetup({cache:false})。这样页面的所有 ajax 都会执行这条语句就是不需要保存缓存记录。

详细资料可以参考：
[《Ajax 中浏览器的缓存问题解决方法》](https://www.cnblogs.com/cwzqianduan/p/8632009.html)
[《浅谈浏览器缓存》](https://segmentfault.com/a/1190000012573337)

#### 58. 同步和异步的区别？

相关知识点：

同步，可以理解为在执行完一个函数或方法之后，一直等待系统返回值或消息，这时程序是处于阻塞的，只有接收到返回的值或消息后才往下执行其他的命令。

异步，执行完函数或方法后，不必阻塞性地等待返回值或消息，只需要向系统委托一个异步过程，那么当系统接收到返回值或消息时，系统会自动触发委托的异步过程，从而完成一个完整的流程。

回答：

同步指的是当一个进程在执行某个请求的时候，如果这个请求需要等待一段时间才能返回，那么这个进程会一直等待下去，直到消息返回为止再继续向下执行。

异步指的是当一个进程在执行某个请求的时候，如果这个请求需要等待一段时间才能返回，这个时候进程会继续往下执行，不会阻塞等待消息的返回，当消息返回时系统再通知进程进行处理。

详细资料可以参考：
[《同步和异步的区别》](https://blog.csdn.net/tennysonsky/article/details/45111623)

#### 59. 什么是浏览器的同源政策？

我对浏览器的同源政策的理解是，`一个域下的 js 脚本在未经允许的情况下，不能够访问另一个域的内容`。这里的同源的指的是两个域的协议、域名、端口号必须相同，否则则不属于同一个域。

`同源政策主要限制了三个方面`

第一个是当前域下的 js 脚本不能够访问其他域下的 cookie、localStorage 和 indexDB。
第二个是当前域下的 js 脚本不能够操作访问操作其他域下的 DOM。
第三个是当前域下 ajax 无法发送跨域请求。

同源政策的目的主要是为了保证用户的信息安全，`它只是对 js 脚本的一种限制，并不是对浏览器的限制`，对于`一般的 img、或者script 脚本请求都不会有跨域的限制`，这是因为这些操作都不会通过响应结果来进行可能出现安全问题的操作。

#### 60. 如何解决跨域问题？

相关知识点：

- 1. 通过 jsonp 跨域
- 2. document.domain + iframe 跨域
- 3. location.hash + iframe
- 4. window.name + iframe 跨域
- 5. postMessage 跨域
- 6. 跨域资源共享（CORS)
- 7. nginx 代理跨域
- 8. nodejs 中间件代理跨域
- 9. WebSocket 协议跨域

回答：

解决跨域的方法我们可以根据我们想要实现的目的来划分。

首先我们如果只是想要实现主域名下的不同子域名的跨域操作，我们可以使用设置 document.domain 来解决。

（1）将 document.domain 设置为主域名，来实现相同子域名的跨域操作，这个时候主域名下的 cookie 就能够被子域名所访问。同时如果文档中含有主域名相同，子域名不同的 iframe 的话，我们也可以对这个 iframe 进行操作。

如果是想要解决不同跨域窗口间的通信问题，比如说一个页面想要和页面的中的不同源的 iframe 进行通信的问题，我们可以使用 location.hash 或者 window.name 或者 postMessage 来解决。

（2）使用 location.hash 的方法，我们可以在主页面动态的修改 iframe 窗口的 hash 值，然后在 iframe 窗口里实现监听函数来实现这样一个单向的通信。因为在 iframe 是没有办法访问到不同源的父级窗口的，所以我们不能直接修改父级窗口的 hash 值来实现通信，我们可以在 iframe 中再加入一个 iframe ，这个 iframe 的内容是和父级页面同源的，所以我们可以 window.parent.parent 来修改最顶级页面的 src，以此来实现双向通信。

（3）使用 window.name 的方法，主要是基于同一个窗口中设置了 window.name 后不同源的页面也可以访问，所以不同源的子页面可以首先在 window.name 中写入数据，然后跳转到一个和父级同源的页面。这个时候级页面就可以访问同源的子页面中 window.name 中的数据了，这种方式的好处是可以传输的数据量大。

（4）使用 postMessage 来解决的方法，这是一个 h5 中新增的一个 api。通过它我们可以实现多窗口间的信息传递，通过获取到指定窗口的引用，然后调用 postMessage 来发送信息，在窗口中我们通过对 message 信息的监听来接收信息，以此来实现不同源间的信息交换。

如果是像解决 ajax 无法提交跨域请求的问题，我们可以使用 jsonp、cors、websocket 协议、服务器代理来解决问题。

（5）使用 jsonp 来实现跨域请求，它的主要原理是通过动态构建 script 标签来实现跨域请求，因为浏览器对 script 标签的引入没有跨域的访问限制 。通过在请求的 url 后指定一个回调函数，然后服务器在返回数据的时候，构建一个 json 数据的包装，这个包装就是回调函数，然后返回给前端，前端接收到数据后，因为请求的是脚本文件，所以会直接执行，这样我们先前定义好的回调函数就可以被调用，从而实现了跨域请求的处理。这种方式只能用于 get 请求。

（6）使用 CORS 的方式，CORS 是一个 W3C 标准，全称是"跨域资源共享"。CORS 需要浏览器和服务器同时支持。目前，所有浏览器都支持该功能，因此我们只需要在服务器端配置就行。浏览器将 CORS 请求分成两类：简单请求和非简单请求。对于简单请求，浏览器直接发出 CORS 请求。具体来说，就是会在头信息之中，增加一个 Origin 字段。Origin 字段用来说明本次请求来自哪个源。服务器根据这个值，决定是否同意这次请求。对于如果 Origin 指定的源，不在许可范围内，服务器会返回一个正常的 HTTP 回应。浏览器发现，这个回应的头信息没有包含 Access-Control-Allow-Origin 字段，就知道出错了，从而抛出一个错误，ajax 不会收到响应信息。如果成功的话会包含一些以 Access-Control- 开头的字段。

非简单请求，浏览器会先发出一次预检请求，来判断该域名是否在服务器的白名单中，如果收到肯定回复后才会发起请求。

（7）使用 websocket 协议，这个协议没有同源限制。

（8）使用服务器来代理跨域的访问请求，就是有跨域的请求操作时发送请求给后端，让后端代为请求，然后最后将获取的结果发返回。

详细资料可以参考：
[《前端常见跨域解决方案（全）》](https://segmentfault.com/a/1190000011145364)
[《浏览器同源政策及其规避方法》](http://www.ruanyifeng.com/blog/2016/04/same-origin-policy.html)
[《跨域，你需要知道的全在这里》](https://juejin.im/entry/59feae9df265da43094488f6)
[《为什么 form 表单提交没有跨域问题，但 ajax 提交有跨域问题？》](https://www.zhihu.com/question/31592553)

#### 62. 简单谈一下 cookie ？

我的理解是 cookie 是服务器提供的一种用于维护会话状态信息的数据，通过服务器发送到浏览器，浏览器保存在本地，当下一次有同源的请求时，将保存的 cookie 值添加到请求头部，发送给服务端。这可以用来实现记录用户登录状态等功能。cookie 一般可以存储 4k 大小的数据，并且只能够被同源的网页所共享访问。

服务器端可以使用 Set-Cookie 的响应头部来配置 cookie 信息。一条 cookie 包括了 9 个属性值 name、value、expires、domain、path、secure、HttpOnly、SameSite、Priority。其中 name 和 value 分别是 cookie 的名字和值。expires 指定了 cookie 失效的时间，domain 是域名、path 是路径，`domain 和 path 一起限制了 cookie 能够被哪些 url 访问`。secure 规定了 cookie 只能在确保安全的情况下传输，`HttpOnly 规定了这个 cookie 只能被服务器访问，不能使用 js 脚本访问`。SameSite 属性用来限制第三方 cookie，可以有效防止 CSRF 攻击，从而减少安全风险。Priority 是 chrome 的提案，定义了三种优先级，当 cookie 数量超出时低优先级的 cookie 会被优先清除。

cookie 跨域:
在发同域请求时，浏览器会将 cookie 自动加在 request header 中
**在发生 xhr 的跨域请求的时候，即使是同源下的 cookie，也不会被自动添加到请求头部，除非显示地规定。**
根本原因是 cookies 也是一种认证信息，在跨域请求中，client 端必须手动设置 **xhr.withCredentials=true**，且 server 端也必须允许 request 能携带认证信息（即 response header 中包含 **Access-Control-Allow-Credentials:true**），这样浏览器才会自动将 cookie 加在 request header 中。

详细资料可以参考：
[《HTTP cookies》 ](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Cookies)
[《聊一聊 cookie》 ](https://segmentfault.com/a/1190000004556040)

#### 64. js 的几种模块规范？

js 中现在比较成熟的有四种模块加载方案。

第一种是 CommonJS 方案，它通过 require 来引入模块，通过 module.exports 定义模块的输出接口。这种模块加载方案是`服务器端的解决方案，它是以同步的方式来引入模块的`，因为在服务端文件都存储在本地磁盘，所以读取非常快，所以以同步的方式加载没有问题。但如果是在浏览器端，由于模块的加载是使用网络请求，因此使用异步加载的方式更加合适。

第二种是 AMD 方案，这种方案采用异步加载的方式来加载模块，模块的加载不影响后面语句的执行，所有依赖这个模块的语句都定
义在一个回调函数里，等到加载完成后再执行回调函数。require.js 实现了 AMD 规范。

第三种是 CMD 方案，这种方案和 AMD 方案都是为了解决异步模块加载的问题，sea.js 实现了 CMD 规范。它和 require.js
的区别在于模块定义时对依赖的处理不同和对依赖模块的执行时机的处理不同。参考 60

第四种方案是 ES6 提出的方案，使用 import 和 export 的形式来导入导出模块。这种方案和上面三种方案都不同。参考 61。

#### 66. ES6 模块与 CommonJS 模块、AMD、CMD 的差异。

- 1.CommonJS 模块输出的是一个`值的拷贝`，ES6 模块输出的是`值的引用`。CommonJS 模块输出的是值的拷贝，也就是说，一旦输出一个值，模块内部的变化就影响不到这个值。ES6 模块的运行机制与 CommonJS 不一样。JS 引擎对脚本静态分析的时候，遇到模块加载命令 import，就会生成一个只读引用。等到脚本真正执行时，再根据这个只读引用，到被加载的那个模块里面去取值。

- 2.CommonJS 模块是运行时加载，ES6 模块是编译时输出接口。CommonJS 模块就是对象，即在输入时是先加载整个模块，生成一个对象，然后再从这个对象上面读取方法，这种加载称为“`运行时加载`”。而 ES6 模块不是对象，它的对外接口只是一种静态定义，在`代码静态解析`阶段就会生成。

#### 67. requireJS 的核心原理是什么？（如何动态加载的？如何避免多次加载的？如何 缓存的？）

require.js 的核心原理是通过动态创建 script 脚本来异步引入模块，然后对每个脚本的 load 事件进行监听，如果每个脚本都加载完成了，再调用回调函数。

详细资料可以参考：
[《requireJS 的用法和原理分析》](https://github.com/HRFE/blog/issues/10)
[《requireJS 的核心原理是什么？》](https://zhuanlan.zhihu.com/p/55039478)
[《从 RequireJs 源码剖析脚本加载原理》](https://www.cnblogs.com/dong-xu/p/7160919.html)
[《requireJS 原理分析》](https://www.jianshu.com/p/5a39535909e4)

#### 70. documen.write 和 innerHTML 的区别？

document.write 的内容会代替整个文档内容，会重写整个页面。

innerHTML 的内容只是替代指定元素的内容，只会重写页面中的部分内容。

详细资料可以参考：
[《简述 document.write 和 innerHTML 的区别。》](https://www.nowcoder.com/questionTerminal/2c5d8105b2694d85b06eff85e871cf50)

#### 74. JavaScript 类数组对象的定义？

一个拥有 length 属性和若干索引属性的对象就可以被称为类数组对象，类数组对象和数组类似，但是不能调用数组的方法。

常见的类数组对象有 arguments 和 DOM 方法的返回结果，还有一个函数也可以被看作是类数组对象，因为它含有 length
属性值，代表可接收的参数个数。

常见的类数组转换为数组的方法有这样几种：

（1）通过 call 调用数组的 slice 方法来实现转换

js
Array.prototype.slice.call(arrayLike);

（4）通过 Array.from 方法来实现转换

js
Array.from(arrayLike);

详细的资料可以参考：
[《JavaScript 深入之类数组对象与 arguments》](https://github.com/mqyqingfeng/Blog/issues/14)
[《javascript 类数组》](https://segmentfault.com/a/1190000000415572)
[《深入理解 JavaScript 类数组》](https://blog.lxxyx.cn/2016/05/07/%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3JavaScript%E7%B1%BB%E6%95%B0%E7%BB%84/)

#### 77. [,,,] 的长度？

`尾后逗号` （有时叫做“终止逗号”）在向 JavaScript 代码添加元素、参数、属性时十分有用。如果你想要添加新的属性，并且上一行已经使用了尾后逗号，你可以仅仅添加新的一行，而不需要修改上一行。这使得版本控制更加清晰，以及代码维护麻烦更少。

JavaScript 一开始就支持数组字面值中的尾后逗号，随后向对象字面值（ECMAScript 5）中添加了尾后逗号。最近（ECMAS
cript 2017），又将其添加到函数参数中。但是 JSON 不支持尾后逗号。

如果使用了多于一个尾后逗号，会产生间隙。 带有间隙的数组叫做稀疏数组（密致数组没有间隙）。稀疏数组的长度为逗号数
量。

详细资料可以参考：
[《尾后逗号》](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Trailing_commas)

#### 78. JavaScript 中的作用域与变量声明提升？

变量提升的表现是，无论我们在函数中何处位置声明的变量，好像都被提升到了函数的首部，我们可以在变量声明前访问到而不会报错。

造成变量声明提升的本质原因是 js 变量作用域 是**静态作用域** 引擎在代码执行前有一个`解析的过程，创建了执行上下文`，初始化了一些代码执行时需要用到的对象。当我们访问一个变量时，我们会到当前执行上下文中的作用域链中去查找，而`作用域链的首端指向的是当前执行上下文的变量对象，这个变量对象是执行上下文的一个属性，它包含了函数的形参、所有的函数和变量声明`，`这个对象的是在代码解析的时候创建的`。这就是会出现变量声明提升的根本原因。

详细资料可以参考：
[《JavaScript 深入理解之变量对象》](http://cavszhouyou.top/JavaScript%E6%B7%B1%E5%85%A5%E7%90%86%E8%A7%A3%E4%B9%8B%E5%8F%98%E9%87%8F%E5%AF%B9%E8%B1%A1.html)

#### 80. 简单介绍一下 V8 引擎的垃圾回收机制

v8 的垃圾回收机制基于`分代回收机制`，这个机制又基于世代假说，这个假说有两个特点，一是新生的对象容易早死，另一个是不死的对象会活得更久。基于这个假说，v8 引擎将内存分为了新生代和老生代。

新创建的对象或者只经历过一次的垃圾回收的对象被称为新生代。经历过多次垃圾回收的对象被称为老生代。

新生代被分为 From 和 To 两个空间，To 一般是闲置的。当 From 空间满了的时候会执行 `Scavenge` 算法进行垃圾回收。当我们执行垃圾回收算法的时候应用逻辑将会停止，等垃圾回收结束后再继续执行。这个算法分为三步：

（1）首先检查 From 空间的存活对象，如果对象存活则判断对象是否满足晋升到老生代的条件，如果满足条件则晋升到老生代。如果不满足条件则移动 To 空间。

（2）如果对象不存活，则释放对象的空间。

（3）最后将 From 空间和 To 空间角色进行交换。

新生代对象晋升到老生代有两个条件：

（1）第一个是判断是对象否已经`经过一次 Scavenge 回收`。若经历过，则将对象从 From 空间复制到老生代中；若没有经历，则复制到 To 空间。

（2）第二个是 To 空间的内存使用占比是否超过限制。当对象从 From 空间复制到 To 空间时，若 `To 空间使用超过 25%`，则对象直接晋升到老生代中。设置 25% 的原因主要是因为算法结束后，两个空间结束后会交换位置，如果 To 空间的内存太小，会影响后续的内存分配。

老生代采用了标记清除法和标记压缩法。标记清除法首先会对内存中存活的对象进行标记，标记结束后清除掉那些没有标记的对象。由于标记清除后会造成很多的内存碎片，不便于后面的内存分配。所以了解决内存碎片的问题引入了标记压缩法。

由于在进行垃圾回收的时候会暂停应用的逻辑，对于新生代方法由于内存小，每次停顿的时间不会太长，但对于老生代来说每次垃圾回收的时间长，停顿会造成很大的影响。 为了解决这个问题 V8 引入了`增量标记`(支持并发，类似于 golang 三色标记混合写屏障)的方法，将一次停顿进行的过程分为了多步，每次执行完一小步就让运行逻辑执行一会，就这样交替运行。

详细资料可以参考：
[《深入理解 V8 的垃圾回收原理》](https://www.jianshu.com/p/b8ed21e8a4fb)
[《JavaScript 中的垃圾回收》](https://zhuanlan.zhihu.com/p/23992332)

#### 81. 哪些操作会造成内存泄漏？

相关知识点：

- 1.意外的全局变量 (链式赋值 let a=b=c=0 ,b,c 都会成为全局变量)
- 2.被遗忘的计时器或回调函数
- 3.脱离 DOM 的引用
- 4.闭包

回答：

第一种情况是我们由于**使用未声明的变量，而意外的创建了一个全局变量**，而使这个变量一直留在内存中无法被回收。

第二种情况是我们设置了 **setInterval 定时器，而忘记取消它**，如果循环函数有对外部变量的引用的话，那么这个变量会被一直留在内存中，而无法被回收。

第三种情况是我们获取一个 DOM 元素的引用，而后面这个元素被删除，由于我们一直保留了对这个元素的引用，所以它也无法被回收。

第四种情况是不合理的使用闭包，从而导致某些变量一直被留在内存当中。

详细资料可以参考：
[《JavaScript 内存泄漏教程》](http://www.ruanyifeng.com/blog/2017/04/memory-leak.html)
[《4 类 JavaScript 内存泄漏及如何避免》](https://jinlong.github.io/2016/05/01/4-Types-of-Memory-Leaks-in-JavaScript-and-How-to-Get-Rid-Of-Them/)
[《杜绝 js 中四种内存泄漏类型的发生》](https://juejin.im/entry/5a64366c6fb9a01c9332c706)
[《javascript 典型内存泄漏及 chrome 的排查方法》](https://segmentfault.com/a/1190000008901861)

#### 83. 如何判断当前脚本运行在浏览器还是 node 环境中？（阿里）

`typeof window === 'undefined' ? 'node' : 'browser';`

通过判断当前环境的 window 对象类型是否为 undefined，如果是 undefined，则说明当前脚本运行在 node 环境，否则说明运行在 window 环境。

#### 86. 什么是“前端路由”？什么时候适合使用“前端路由”？“前端路由”有哪些优点和缺点？

（1）什么是前端路由？
前端路由就是把不同路由对应不同的内容或页面的任务交给前端来做，之前是通过服务端根据 url 的不同返回不同的页面实现的。

（2）什么时候使用前端路由？
在单页面应用，大部分页面结构不变，只改变部分内容的使用

（3）前端路由有什么优点和缺点？
优点：用户体验好，不需要每次都从服务器全部获取，快速展现给用户
缺点：单页面无法`记住之前滚动的位置`，无法在前进，后退的时候记住滚动的位置
前端路由一共有两种实现方式，一种是通过 hash 的方式，一种是通过使用 pushState 的方式。

详细资料可以参考：
[《什么是“前端路由”》](https://segmentfault.com/q/1010000005336260)
[《浅谈前端路由》 ](https://github.com/kaola-fed/blog/issues/137)
[《前端路由是什么东西？》](https://www.zhihu.com/question/53064386)

#### 88. 检测浏览器版本版本有哪些方式？

检测浏览器版本一共有两种方式：
一种是检测 `window.navigator.userAgent` 的值，但这种方式很不可靠，因为 `userAgent 可以被改写`，并且早期的浏览器如 ie，会通过伪装自己的 userAgent 的值为 Mozilla 来躲过服务器的检测。
第二种方式是功能检测，`根据每个浏览器独有的特性来进行判断`，如 ie 下独有的 ActiveXObject。

详细资料可以参考：
[《JavaScript 判断浏览器类型》](https://www.jianshu.com/p/d99f4ca385ac)

#### 89. 什么是 Polyfill ？

Polyfill 指的是用于`实现浏览器并不支持的原生 API 的代码`。

比如说 querySelectorAll 是很多现代浏览器都支持的原生 Web API，但是有些古老的浏览器并不支持，那么假设有人写了一段代码来实现这个功能使这些浏览器也支持了这个功能，那么这就可以成为一个 Polyfill。

一个 shim 是一个库，有自己的 API，而`不是单纯实现原生不支持的 API`。

详细资料可以参考：
[《Web 开发中的“黑话”》](https://segmentfault.com/a/1190000002593432)
[《Polyfill 为何物》](https://juejin.im/post/5a579bc7f265da3e38496ba1)

#### 101. toPrecision 和 toFixed 和 Math.round 的区别？

toPrecision 用于处理精度，精度是`从左至右第一个不为 0 的数开始数起。`
toFixed 是对小数点后指定位数取整，`从小数点开始数起。`
Math.round 是将一个数字四舍五入到一个整数。

#### 102. 什么是 XSS 攻击？如何防范 XSS 攻击？

XSS 攻击指的是跨站脚本攻击，是一种代码注入攻击。攻击者通过在网站注入恶意脚本，使之在用户的浏览器上运行，从而盗取用户的信息如 cookie 等。

XSS 的本质是因为网站没有对恶意代码进行过滤，与正常的代码混合在一起了，浏览器没有办法分辨哪些脚本是可信的，从而导致了恶意代码的执行。

XSS 一般分为存储型、反射型和 DOM 型。

存储型指的是恶意代码提交到了网站的`数据库中`，当用户请求数据的时候，服务器将其拼接为 HTML 后返回给了用户，从而导致了恶意代码的执行。

反射型指的是攻击者构建了特殊的 URL，当服务器接收到请求后，`从 URL 中获取数据，拼接到 HTML 后返回，从而导致了恶意代码的执行`。

DOM 型(不经过 server)指的是攻击者构建了特殊的 URL，用户打开网站后，`js 脚本从 URL 中获取数据`，从而导致了恶意代码的执行。

XSS 攻击的预防可以从两个方面入手，一个是恶意代码提交的时候，一个是浏览器执行恶意代码的时候。

对于第一个方面，如果我们对存入数据库的数据都进行的转义处理，但是一个数据可能在多个地方使用，有的地方可能不需要转义，由于我们没有办法判断数据最后的使用场景，所以直接在输入端进行恶意代码的处理，其实是不太可靠的。

因此我们可以从浏览器的执行来进行预防，一种是使用纯前端的方式，不用服务器端拼接后返回。另一种是对需要插入到 HTML 中的代码做好充分的转义。对于 DOM 型的攻击，主要是前端脚本的不可靠而造成的，我们对于数据获取渲染和字符串拼接的时候应该对可能出现的恶意代码情况进行判断。

还有一些方式，比如使用 `CSP ，CSP 的本质是建立一个白名单，告诉浏览器哪些外部资源可以加载和执行，从而防止恶意代码的注入攻击。`

还可以对一些敏感信息进行保护，比如 cookie 使用 http-only ，使得脚本无法获取。也可以使用验证码，避免脚本伪装成用户执行一些操作。

详细资料可以参考：
[《前端安全系列（一）：如何防止 XSS 攻击？》](https://juejin.im/post/5bad9140e51d450e935c6d64)

#### 103. 什么是 CSP？

CSP 指的是**内容安全策略**，它的本质是建立一个白名单，告诉浏览器哪些外部资源可以加载和执行。我们只需要配置规则，如何拦截由浏览器自己来实现。

通常有两种方式来开启 CSP，一种是设置 **HTTP 首部中的 Content-Security-Policy**，一种是设置 **meta 标签**的方式 <meta http-equiv="Content-Security-Policy">

详细资料可以参考：
[《内容安全策略（CSP）》](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CSP)
[《前端面试之道》](https://juejin.im/book/5bdc715fe51d454e755f75ef/section/5bdc721851882516c33430a2)

#### 105. 什么是 Samesite Cookie 属性？

Samesite Cookie 表示**同站 cookie，避免 cookie 被第三方所利用**。

将 Samesite 设为 strict ，这种称为严格模式，表示这个 cookie 在任何情况下都不可能作为第三方 cookie。

将 Samesite 设为 Lax ，这种模式称为宽松模式，如果这个请求是个 GET 请求，并且这个请求改变了当前页面或者打开了新的页面，那么这个 cookie 可以作为第三方 cookie，其余情况下都不能作为第三方 cookie。

使用这种方法的缺点是，因为它不支持子域，所以子域没有办法与主域共享登录信息，每次转入子域的网站，都回重新登录。还有一个问题就是它的兼容性不够好。

#### 110. Object.defineProperty 介绍？

Object.defineProperty 函数一共有三个参数，第一个参数是需要定义属性的对象，第二个参数是需要定义的属性，第三个是该属性描述符。

一个属性的描述符有四个属性，分别是 value 属性的值，writable 属性是否可写，enumerable 属性是否可枚举，configurable 属性是否可配置修改。

详细资料可以参考：
[《Object.defineProperty()》](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Object/defineProperty)

#### 111. 使用 Object.defineProperty() 来进行数据劫持有什么缺点？

有一些对属性的操作，使用这种方法无法拦截，比如说通过下`标方式修改数组数据或者给对象新增属性`，vue 内部通过重写函数解决了这个问题。在 Vue3.0 中已经不使用这种方式了，而是通过使用 Proxy 对对象进行代理，从而实现数据劫持。使用 Proxy 的好处是它可以完美的监听到任何方式的数据改变，唯一的缺点是兼容性的问题，因为这是 ES6 的语法。

#### 112. 什么是 Virtual DOM？为什么 Virtual DOM 比原生 DOM 快？

我对 Virtual DOM 的理解是，

首先对我们将要插入到文档中的 DOM 树结构进行分析，使用 js 对象将其表示出来，比如一个元素对象，包含 TagName、props 和 Children 这些属性。然后我们将这个 js 对象树给保存下来，最后再将 DOM 片段插入到文档中。

当页面的状态发生改变，我们需要对页面的 DOM 的结构进行调整的时候，我们首先根据变更的状态，重新构建起一棵对象树，然后将这棵新的对象树和旧的对象树进行比较，记录下两棵树的的差异。

最后将记录的有差异的地方应用到真正的 DOM 树中去，这样视图就更新了。

我认为 Virtual DOM 这种方法对于我们需要有大量的 DOM 操作的时候，能够很好的提高我们的操作效率，通过在操作前确定需要做的最小修改，`尽可能的减少 DOM 操作带来的重流和重绘的影响`。其实 Virtual DOM 并不一定比我们真实的操作 DOM 要快，这种方法的目的是为了提高我们开发时的可维护性，在任意的情况下，都能保证一个尽量小的性能消耗去进行操作。

详细资料可以参考：
[《Virtual DOM》](https://juejin.im/book/5bdc715fe51d454e755f75ef/section/5bdc72e6e51d45054f664dbf)
[《理解 Virtual DOM》](https://github.com/y8n/blog/issues/5)
[《深度剖析：如何实现一个 Virtual DOM 算法》](https://github.com/livoras/blog/issues/13)
[《网上都说操作真实 DOM 慢，但测试结果却比 React 更快，为什么？》](https://www.zhihu.com/question/31809713/answer/53544875)

#### 113. 如何比较两个 DOM 树的差异？

两个树的完全 diff 算法的时间复杂度为 O(n^3) ，但是在前端中，我们很少会跨层级的移动元素，所以我们`只需要比较同一层级的元素进行比较`，这样就可以将算法的时间复杂度降低为 O(n)。

算法首先会对新旧两棵树进行一个深度优先的遍历，这样每个节点都会有一个序号。在深度遍历的时候，每遍历到一个节点，我们就将这个节点和新的树中的节点进行比较，如果有差异，则将这个差异记录到一个对象中。

在对列表元素进行对比的时候，由于 TagName 是重复的，所以我们不能使用这个来对比。我们需要给每一个子节点加上一个 key，列表对比的时候使用 key 来进行比较，这样我们才能够复用老的 DOM 树上的节点。

#### 116. offsetWidth/offsetHeight,clientWidth/clientHeight 与 scrollWidth/scrollHeight 的区别？

clientWidth/clientHeight 返回的是元素的内部宽度，它的值只包含 content + padding，如果有滚动条，不包含滚动条。
clientTop 返回的是`上边框的宽度`。
clientLeft 返回的`左边框的宽度`。

offsetWidth/offsetHeight 返回的是元素的布局宽度，它的值包含 content + padding + border 包含了滚动条。
offsetTop 返回的是当前元素相对于其 `offsetParent 元素`的顶部的距离。
offsetLeft 返回的是当前元素相对于其 `offsetParent 元素`的左部的距离。

scrollWidth/scrollHeight 返回值包含 content + padding + 溢出内容的尺寸。(`整个页面尺寸`)
scrollTop 属性返回的是一个元素的内容垂直滚动的像素数。(`已经滚动了多高`)
scrollLeft 属性返回的是元素滚动条到元素左边的距离。(`已经滚动了多宽`)

详细资料可以参考：
[《最全的获取元素宽高及位置的方法》](https://juejin.im/post/5bc9366d5188255c4834e75a)
[《用 Javascript 获取页面元素的位置》](http://www.ruanyifeng.com/blog/2009/09/find_element_s_position_using_javascript.html)

#### 117. 谈一谈你理解的函数式编程？

简单说，"函数式编程"是一种"编程范式"（programming paradigm），也就是如何编写程序的方法论。

它具有以下特性：`闭包`和`高阶函数`、`惰性计算`、`递归`、函数是`"第一等公民"`、只用"表达式"、`纯函数`。

详细资料可以参考：
[《函数式编程初探》](http://www.ruanyifeng.com/blog/2012/04/functional_programming.html)

#### 118. 异步编程的实现方式？

相关资料：

回调函数
优点：简单、容易理解
缺点：不利于维护，代码耦合高

事件监听（采用时间驱动模式，取决于某个事件是否发生）：
优点：容易理解，可以绑定多个事件，每个事件可以指定多个回调函数
缺点：事件驱动型，流程不够清晰

发布/订阅（观察者模式）
类似于事件监听，但是可以通过‘消息中心’，了解现在有多少发布者，多少订阅者

Promise 对象
优点：可以利用 then 方法，进行链式写法；可以书写错误时的回调函数；
缺点：编写和理解，相对比较难

Generator 函数
优点：函数体内外的数据交换、错误处理机制
缺点：流程管理不方便

async 函数
优点：内置执行器、更好的语义、更广的适用性、返回的是 Promise、结构清晰。
缺点：错误处理机制

回答：

js 中的异步机制可以分为以下几种：

第一种最常见的是使用回调函数的方式，使用回调函数的方式有一个缺点是，多个回调函数嵌套的时候会造成回调函数地狱，上下两层的回调函数间的代码耦合度太高，不利于代码的可维护。

第二种是 Promise 的方式，使用 Promise 的方式可以将嵌套的回调函数作为链式调用。但是使用这种方法，有时会造成多个 then 的链式调用，可能会造成代码的语义不够明确。

第三种是使用 generator 的方式，它可以在函数的执行过程中，将函数的执行权转移出去，在函数外部我们还可以将执行权转移回来。当我们遇到异步函数执行的时候，将函数执行权转移出去，当异步函数执行完毕的时候我们再将执行权给转移回来。因此我们在 generator 内部对于异步操作的方式，可以以同步的顺序来书写。使用这种方式我们需要考虑的问题是何时将函数的控制权转移回来，因此我们需要有一个自动执行 generator 的机制，比如说 co 模块等方式来实现 generator 的自动执行。

第四种是使用 async 函数的形式，async 函数是 generator 和 promise 实现的一个自动执行的语法糖，它内部自带执行器，当函数内部执行到一个 await 语句的时候，如果语句返回一个 promise 对象，那么函数将会等待 promise 对象的状态变为 resolve 后再继续向下执行。因此我们可以将异步逻辑，转化为同步的顺序来书写，并且这个函数可以自动执行。

#### 120. get 请求传参长度的误区

误区：我们经常说 get 请求参数的大小存在限制，而 post 请求的参数大小是无限制的。

实际上 `HTTP 协议从未规定 GET/POST 的请求长度限制是多少`。对 get 请求参数的限制是来源与`浏览器或 web 服务器`，`浏览器或 web 服务器限制了 url 的长度`。为了明确这个概念，我们必须再次强调下面几点:

- 1.`HTTP 协议未规定 GET 和 POST 的长度限制`
- 2.GET 的最大长度显示是因为`浏览器和 web 服务器限制了 URI 的长度`
- 3.不同的浏览器和 WEB 服务器，限制的最大长度不一样
- 4.要支持 IE，则最大长度为 2083byte，若只支持 Chrome，则最大长度 `8182byte(8KB)`

#### 122. get 和 post 请求在缓存方面的区别

相关知识点：

get 请求类似于查找的过程，用户获取数据，可以不用每次都与数据库连接，所以`可以使用缓存`。

post 不同，post 做的一般是修改和删除的工作，所以必须与数据库交互，所以不能使用缓存。因此 get 请求适合于请求缓存。

回答：

缓存一般只适用于那些不会更新服务端数据的请求。一般 get 请求都是查找请求，不会对服务器资源数据造成修改，而 post 请求一般都会对服务器数据造成修改，所以，一般会对 get 请求进行缓存，很少会对 post 请求进行缓存。

详细资料可以参考：
[《HTML 关于 post 和 get 的区别以及缓存问题的理解》](https://blog.csdn.net/qq_27093465/article/details/50479289)

#### 123. 图片的懒加载和预加载

相关知识点：

预加载：提前加载图片，当用户需要查看时可直接从本地缓存中渲染。

懒加载：懒加载的主要目的是作为服务器前端的优化，减少请求数或延迟请求数。

两种技术的本质：两者的行为是相反的，一个是提前加载，一个是迟缓甚至不加载。 懒加载对服务器前端有一定的缓解压力作用，预加载则会增加服务器前端压力。

回答：

懒加载也叫延迟加载，指的是在长网页中延迟加载图片的时机，当用户需要访问时，再去加载，这样可以提高网站的首屏加载速度，提升用户的体验，并且可以减少服务器的压力。它适用于图片很多，页面很长的电商网站的场景。懒加载的实现原理是，将页面上的图片的 src 属性设置为空字符串，将图片的真实路径保存在一个自定义属性中，当页面滚动的时候，进行判断，如果图片进入页面可视区域内，则从自定义属性中取出真实路径赋值给图片的 src 属性，以此来实现图片的延迟加载。

预加载指的是将所需的资源提前请求加载到本地，这样后面在需要用到时就直接从缓存取资源。通过预加载能够减少用户的等待时间，提高用户的体验。我了解的预加载的最常用的方式是使用 js 中的 image 对象，通过为 image 对象来设置 scr 属性，来实现图片的预加载。

这两种方式都是提高网页性能的方式，两者主要区别是一个是提前加载，一个是迟缓甚至不加载。懒加载对服务器前端有一定的缓解压力作用，预加载则会增加服务器前端压力。

详细资料可以参考：
[《懒加载和预加载》](https://juejin.im/post/5b0c3b53f265da09253cbed0)
[《网页图片加载优化方案》](https://juejin.im/entry/5a73f38cf265da4e99575be3)
[《基于用户行为的图片等资源预加载》](https://www.zhangxinxu.com/wordpress/2016/06/image-preload-based-on-user-behavior/)

#### 124. mouseover 和 mouseenter 的区别？

当鼠标移动到元素上时就会触发 mouseenter 事件，类似 mouseover，它们两者之间的差别是 `mouseenter 不会冒泡`。

由于 mouseenter 不支持事件冒泡，导致在一个元素的子元素上进入或离开的时候会触发其 mouseover 和 mouseout 事件，但是却不会触发 mouseenter 和 mouseleave 事件。

详细资料可以参考：
[《mouseenter 与 mouseover 为何这般纠缠不清？》](https://github.com/qianlongo/zepto-analysis/issues/1)

#### 125. js 拖拽功能的实现

相关知识点：

首先是三个事件，分别是 mousedown，mousemove，mouseup
当鼠标点击按下的时候，需要一个 tag 标识此时已经按下，可以执行 mousemove 里面的具体方法。
clientX，clientY 标识的是鼠标的坐标，分别标识横坐标和纵坐标，并且我们用 offsetX 和 offsetY 来表示元素的元素的初始坐标，移动的举例应该是：
鼠标移动时候的坐标-鼠标按下去时候的坐标。
也就是说定位信息为：
鼠标移动时候的坐标-鼠标按下去时候的坐标+元素初始情况下的 offetLeft.

回答：

一个元素的拖拽过程，我们可以分为三个步骤，第一步是鼠标按下目标元素，第二步是鼠标保持按下的状态移动鼠标，第三步是鼠标抬起，拖拽过程结束。

这三步分别对应了三个事件，mousedown 事件，mousemove 事件和 mouseup 事件。只有在鼠标按下的状态移动鼠标我们才会执行拖拽事件，因此我们需要在 mousedown 事件中设置一个状态来标识鼠标已经按下，然后在 mouseup 事件中再取消这个状态。`在 mousedown 事件中我们首先应该判断，目标元素是否为拖拽元素，如果是拖拽元素，我们就设置状态并且保存这个时候鼠`标的位置。然后在 mousemove 事件中，我们通过判断鼠标现在的位置和以前位置的相对移动，来确定拖拽元素在移动中的坐标。
`最后 mouseup 事件触发后，清除状态，结束拖拽事件。`

详细资料可以参考：
[《原生 js 实现拖拽功能基本思路》](https://blog.csdn.net/LZGS_4/article/details/43523465)

回答：

setInterval 的作用是每隔一段指定时间执行一个函数，但是这个执行不是真的到了时间立即执行，它真正的作用是`每隔一段时间将事件加入事件队列中去，只有当当前的执行栈为空的时候，才能去从事件队列中取出事件执行`。所以可能会出现这样的情况，就是当前执行栈执行的时间很长，导致事件队列里边积累多个定时器加入的事件，当执行栈结束的时候，这些事件会依次执行，因此就不能到间隔一段时间执行的效果。

针对 setInterval 的这个缺点，我们可以使用 setTimeout 递归调用来模拟 setInterval，这样我们就`确保了只有一个事件结束了，我们才会触发下一个定时器事件`，这样解决了 setInterval 的问题。

详细资料可以参考：
[《用 setTimeout 实现 setInterval》](https://www.jianshu.com/p/32479bdfd851)
[《setInterval 有什么缺点？》](https://zhuanlan.zhihu.com/p/51995737)

#### 127. let 和 const 的注意点？

- 1.声明的变量只在声明时的代码块内有效
- 2.不存在声明提升
- 3.存在暂时性死区，如果在变量声明前使用，会报错
- 4.不允许重复声明，重复声明会报错

#### 129. 什么是尾调用，使用尾调用有什么好处？

尾调用指的是函数的最后一步调用另一个函数。我们代码执行是基于执行栈的，所以当我们在一个函数里调用另一个函数时，我们会保留当前的执行上下文，然后再新建另外一个执行上下文加入栈中。使用尾调用的话，因为已经是函数的最后一步，所以这个时候我们可以`不必再保留当前的执行上下文`，从而节省了内存，这就是尾调用优化。但是 ES6 的尾调用优化只在严格模式下开启，正常模式是无效的。

#### 130. Symbol 类型的注意点？

- 1.Symbol 函数前不能使用 new 命令，否则会报错。
- 2.Symbol 函数可以接受一个字符串作为参数，表示对 Symbol 实例的描述，主要是为了在控制台显示，或者转为字符串时，比较容易区分。
- 3.Symbol 作为属性名，该属性不会出现在 for...in、for...of 循环中，也不会被 Object.keys()、Object.getOwnPropertyNames()、JSON.stringify() 返回。
- 4.Object.getOwnPropertySymbols 方法返回一个数组，成员是当前对象的所有用作属性名的 Symbol 值。
- 5.Symbol.for 接受一个字符串作为参数，然后搜索有没有以该参数作为名称的 Symbol 值。如果有，就返回这个 Symbol 值，否则就新建并返回一个以该字符串为名称的 Symbol 值。
- 6.Symbol.keyFor 方法返回一个已登记的 Symbol 类型值的 key。

#### 131. Set 和 WeakSet 结构？

- 1.ES6 提供了新的数据结构 Set。它类似于数组，但是成员的值都是唯一的，没有重复的值。
- 2.WeakSet 结构与 Set 类似，也是不重复的值的集合。但是 WeakSet 的成员只能是对象，而不能是其他类型的值。WeakSet 中的对象都是弱引用，即`垃圾回收机制不考虑 WeakSet 对该对象的引用`，

#### 134. Reflect 对象创建目的？

- 1.将 Object 对象的一些明显属于语言内部的方法（比如 Object.defineProperty，放到 Reflect 对象上。
- 2.修改某些 Object 方法的返回结果，让其变得更合理。
- 3.让 Object 操作都变成函数行为。
- 4.Reflect 对象的方法与 Proxy 对象的方法一一对应，只要是 Proxy 对象的方法，就能在 Reflect 对象上找到对应的方法。这就让 Proxy 对象可以方便地调用对应的 Reflect 方法，完成默认行为，作为修改行为的基础。也就是说，不管 Proxy 怎么修改默认行为，你总可以在 Reflect 上获取默认行为。

#### 135. require 模块引入的查找方式？

当 Node 遇到 require(X) 时，按下面的顺序处理。

（1）如果 X 是内置模块（比如 require('http')）
　　 a. 返回该模块。
　　 b. 不再继续执行。

（2）如果 X 以 "./" 或者 "/" 或者 "../" 开头
　　 a. 根据 X 所在的父模块，确定 X 的绝对路径。
　　 b. 将 X 当成文件，依次查找下面文件，只要其中有一个存在，就返回该文件，不再继续执行。
X
X.js
X.json
X.node

c. 将 X 当成目录，依次查找下面文件，只要其中有一个存在，就返回该文件，不再继续执行。
X/package.json（main 字段）
X/index.js
X/index.json
X/index.node

（3）如果 X 不带路径
　　 a. 根据 X 所在的父模块，确定 X 可能的安装目录。
　　 b. 依次在每个目录中，将 X 当成文件名或目录名加载。

（4）抛出 "not found"

详细资料可以参考：
[《require() 源码解读》](http://www.ruanyifeng.com/blog/2015/05/require.html)

#### 138. 如何检测浏览器所支持的最小字体大小？

用 JS 设置 DOM 的字体为某一个值，然后再取出来，如果值设置成功，就说明支持。

#### 140. 单例模式模式是什么？

单例模式保证了全局只有一个实例来被访问。比如说常用的如弹框组件的实现和全局状态的实现。

#### 141. 策略模式是什么？ 查表

策略模式主要是用来将方法的实现和方法的调用分离开，外部通过不同的参数可以调用不同的策略。我主要在 MVP 模式解耦的时候
用来将视图层的方法定义和方法调用分离。

#### 142. 代理模式是什么？

代理模式是为一个对象提供一个代用品或占位符，以便控制对它的访问。比如说常见的`事件代理。`

#### 143. 中介者模式是什么？

中介者模式指的是，多个对象通过一个中介者进行交流，而不是直接进行交流，这样能够将通信的各个对象解耦。

#### 144. 适配器模式是什么？

适配器用来解决两个接口不兼容的情况，不需要改变已有的接口，通过包装一层的方式实现两个接口的正常协作。假如我们需要一种
新的接口返回方式，但是老的接口由于在太多地方已经使用了，不能随意更改，这个时候就可以使用适配器模式。比如我们需要一种
自定义的时间返回格式，但是我们又不能对 js 时间格式化的接口进行修改，这个时候就可以使用适配器模式。

更多关于设计模式的资料可以参考：
[《前端面试之道》](https://juejin.im/book/5bdc715fe51d454e755f75ef/section/5bdc74186fb9a049ab0d0b6b)
[《JavaScript 设计模式》](https://juejin.im/post/59df4f74f265da430f311909#heading-3)
[《JavaScript 中常见设计模式整理》](https://juejin.im/post/5afe6430518825428630bc4d)

#### 145. 观察者模式和发布订阅模式有什么不同？

发布订阅模式其实属于广义上的观察者模式

在观察者模式中，观察者 Observer 需要直接订阅目标事件 Observable。在目标发出内容改变的事件后，直接接收事件并作出响应。

而在发布订阅模式中，发布者和订阅者之间`多了一个调度中心`。调度中心一方面从发布者接收事件，另一方面向订阅者发布事件，订阅者需要在调度中心中订阅事件。通过调度中心实现了发布者和订阅者关系的解耦。使用发布订阅者模式更利于我们代码的可维护性。

详细资料可以参考：
[《观察者模式和发布订阅模式有什么不同？》](https://www.zhihu.com/question/23486749)

#### 147. Vue 的各个生命阶段是什么？

Vue 一共有 8 个生命阶段，分别是创建前、创建后、加载前、加载后、更新前、更新后、销毁前和销毁后，每个阶段对应了一个生命周期的钩子函数。

（1）beforeCreate 钩子函数，在实例初始化之后，在数据监听和事件配置之前触发。因此在这个事件中我们是获取不到 data 数据的。

（2）created 钩子函数，在实例创建完成后触发，此时可以访问 data、methods 等属性。但这个时候组件还没有被挂载到页面中去，所以这个时候访问不到 $el 属性。一般我们可以在这个函数中进行一些页面初始化的工作，比如通过 ajax 请求数据来对页面进行初始化。

（3）beforeMount 钩子函数，在组件被挂载到页面之前触发。在 beforeMount 之前，会找到对应的 template，并编译成 render 函数。

（4）mounted 钩子函数，在组件挂载到页面之后触发。此时可以通过 DOM API 获取到页面中的 DOM 元素。

（5）beforeUpdate 钩子函数，在响应式数据更新时触发，发生在虚拟 DOM 重新渲染和打补丁之前，这个时候我们可以对可能会被移除的元素做一些操作，比如移除事件监听器。

（6）updated 钩子函数，虚拟 DOM 重新渲染和打补丁之后调用。

（7）beforeDestroy 钩子函数，在实例销毁之前调用。一般在这一步我们可以销毁定时器、解绑全局事件等。

（8）destroyed 钩子函数，在实例销毁之后调用，调用后，Vue 实例中的所有东西都会解除绑定，所有的事件监听器会被移除，所有的子实例也会被销毁。

当我们使用 keep-alive 的时候，还有两个钩子函数，分别是 activated 和 deactivated 。用 keep-alive 包裹的组件在切换时不会进行销毁，而是缓存到内存中并执行 deactivated 钩子函数，命中缓存渲染后会执行 actived 钩子函数。

详细资料可以参考：
[《vue 生命周期深入》](https://juejin.im/entry/5aee8fbb518825671952308c)
[《Vue 实例》](https://cn.vuejs.org/v2/guide/instance.html)

#### 150. vue-router 中的导航钩子函数

（1）全局的钩子函数 beforeEach 和 afterEach

beforeEach 有三个参数，to 代表要进入的路由对象，from 代表离开的路由对象。next 是一个必须要执行的函数，如果不传参数，那就执行下一个钩子函数，如果传入 false，则终止跳转，如果传入一个路径，则导航到对应的路由，如果传入 error ，则导航终止，error 传入错误的监听函数。

（2）单个路由独享的钩子函数 beforeEnter，它是在路由配置上直接进行定义的。

（3）组件内的导航钩子主要有这三种：beforeRouteEnter、beforeRouteUpdate、beforeRouteLeave。它们是直接在路由组
件内部直接进行定义的。

详细资料可以参考：
[《导航守卫》](https://router.vuejs.org/zh/guide/advanced/navigation-guards.html#%E5%85%A8%E5%B1%80%E5%89%8D%E7%BD%AE%E5%AE%88%E5%8D%AB)

#### 151. $route 和 $router 的区别？

$route 是“`路由信息对象`”，包括 path，params，hash，query，fullPath，matched，name 等路由信息参数。
而 $router 是“`路由实例`”对象包括了路由的跳转方法，钩子函数等。

#### 152. vue 常用的修饰符？

.prevent: 提交事件不再重载页面；.stop: 阻止单击事件冒泡；.self: 当事件发生在该元素本身而不是子元素的时候会触发；

#### 157. 开发中常用的几种 Content-Type ？

（1）application/x-www-form-urlencoded

浏览器的原生 form 表单，如果不设置 enctype 属性，那么最终就会以 application/x-www-form-urlencoded 方式提交数据。该种方式提交的数据放在 body 里面，数据按照 key1=val1&key2=val2 的方式进行编码，key 和 val 都进行了 URL 转码。

（2）multipart/form-data

该种方式也是一个常见的 POST 提交方式，通常表单上传文件时使用该种方式。

（3）application/json

告诉服务器消息主体是序列化后的 JSON 字符串。

（4）text/xml

该种方式主要用来提交 XML 格式的数据。

详细资料可以参考：
[《常用的几种 Content-Type》](https://honglu.me/2015/07/13/%E5%B8%B8%E7%94%A8%E7%9A%84%E5%87%A0%E7%A7%8DContent-Type/)

#### 165. 如何确定页面的可用性时间，什么是 Performance API？

Performance API 用于精确度量、控制、增强浏览器的性能表现。这个 API 为测量网站性能，提供以前没有办法做到的精度。

使用 getTime 来计算脚本耗时的缺点，首先，getTime 方法（以及 Date 对象的其他方法）都只能精确到`毫秒级别`（一秒的千分之一），想要得到更小的时间差别就无能为力了。其次，这种写法只能获取代码运行过程中的时间进度，无法知道一些后台事件的时间进度，比如浏览器用了多少时间从服务器加载网页。

为了解决这两个不足之处，ECMAScript 5 引入“高精度时间戳”这个 API，部署在 performance 对象上。它的精度可以达到 `微秒级别`。

navigationStart：当前浏览器窗口的前一个网页关闭，发生 `unload 事件时的 Unix 毫秒时间戳`。如果没有前一个网页，则等于 fetchStart 属性。

loadEventEnd：返回当前网页 `load 事件的回调函数运行结束时`的 Unix 毫秒时间戳。如果该事件还没有发生，返回 0。

根据上面这些属性，可以计算出网页加载各个阶段的耗时。比如，网页加载整个过程的耗时的计算方法如下：

首屏时间
js
var t = performance.timing;
var pageLoadTime = t.loadEventEnd - t.navigationStart;

详细资料可以参考：
[《Performance API》](http://javascript.ruanyifeng.com/bom/performance.html)

#### 167. js 语句末尾分号是否可以省略？

在 ECMAScript 规范中，语句结尾的分号并不是必需的。但是我们一般最好不要省略分号，因为加上分号一方面有
利于我们代码的可维护性，另一方面也可以避免我们在对代码进行压缩时出现错误。

#### 168. Object.assign()

Object.assign() 方法用于将所有`可枚举属性的值`从一个或多个源对象复制到目标对象。它将返回目标对象。

#### 170. js for 循环注意点

js
for (var i = 0, j = 0; i < 5, j < 9; i++, j++) {
console.log(i, j);
}

// 当判断语句含有多个语句时，`以最后一个判断语句的值为准，因此上面的代码会执行 10 次`。
// `当判断语句为空时，循环会一直进行`。

#### 171. 一个列表，假设有 100000 个数据，这个该怎么办？

我们需要思考的问题：该处理是否必须同步完成？数据是否必须按顺序完成？

解决办法：

（1）将数据`分页`，利用分页的原理，每次服务器端只返回一定数目的数据，浏览器每次只对一部分进行加载。

（2）使用`懒加载`的方法，每次加载一部分数据，其余数据当需要使用时再去加载。

（3）使用数组`分块`技术，基本思路是为要处理的项目创建一个队列，然后设置定时器每过一段时间取出一部分数据，然后再使用定时器取出下一个要处理的项目进行处理，接着再设置另一个定时器。

#### 172. js 中倒计时的纠偏实现？

在前端实现中我们一般通过 setTimeout 和 setInterval 方法来实现一个倒计时效果。但是使用这些方法会存在时间偏差的问题，这是由于 js 的程序执行机制造成的，setTimeout 和 setInterval 的作用是`隔一段时间将回调事件加入到事件队列中，因此事件并不是立即执行的`，`它会等到当前执行栈为空的时候再取出事件执行`，因此事件等待执行的时间就是造成误差的原因。

一般解决倒计时中的误差的有这样两种办法：

（1）第一种是通过`前端定时向服务器发送请求获取最新的时间差`，以此来校准倒计时时间。

（2）第二种方法是前端根据偏差时间来自动调整间隔时间的方式来实现的。这一种方式首先是以 setTimeout 递归的方式来实现倒计时，然后`通过一个变量来记录已经倒计时的秒数`。每一次函数调用的时候，首先将变量加一，然后根据这个变量和每次的间隔时间，我们就可以计算出此时无偏差时应该显示的时间。然后将当前的真实时间与这个时间相减，这样我们就可以得到时间的偏差大小，`因此我们在设置下一个定时器的间隔大小的时候，我们就从间隔时间中减去这个偏差大小`，以此来实现由于程序执行所造成的时间误差的纠正。

详细资料可以参考：
[《JavaScript 前端倒计时纠偏实现》](https://juejin.im/post/5badf8305188255c8e728adc)
