1. ['1', '2', '3'].map(parseInt) what & why ?
   parseInt('1', 0) //radix 为 0 时，且 string 参数不以“0x”和“0”开头时，按照 10 为基数处理。这个时候返回 1
   parseInt('2', 1) //基数为 1（1 进制）表示的数中，最大值小于 2，所以无法解析，返回 NaN
   parseInt('3', 2) //基数为 2（2 进制）表示的数中，最大值小于 3，所以无法解析，返回 NaN
2. 介绍模块化发展历程
   模块化主要是用来抽离公共代码，隔离作用域，避免变量冲突等。

   IIFE => CommonJS => UMD =>ES Modules
   commonJS:特点: require、module.exports、exports
   文件即模块，**每个文件通过 module 来表示**，用 require 来引用其他依赖，用 module.exports 来导出自身
   通过 require 去引用文件时，**会将文件执行一遍后**，将其执行结果通过浅克隆的方式，写入全局内存。后续再 require 该路径，就直接从内存里取出，不需要重新执行对应的文件
   **是服务器编程范式**，因为服务器上所有文件都在硬盘里,通过同步加载的方式即可
   ES Modules
   ES6 的模块机制在依赖模块时并**不会先去预加载整个脚本**，而是生成一个只读引用，并且静态解析依赖，等到执行代码时，再去依赖里取出实际需要的模块
   编译时加载，不允许在里边引用变量，必须为真实的文件路径。可以通过调用 import()语句，会生成一个 promise 去加载对应的文件，这样子就是**运行时加载**，可以在路径里边编写变量

   **CommonJS 与 ES6 Modules 规范的区别**
   加载：CommonJS 模块是**运行时加载**，ES6 Modules 是**编译时输出**接口
   输出：CommonJS **输出是值的拷贝**；ES6 Modules **输出的是值的引用**，被输出模块的内部的改变会影响引用的改变
   CommonJs 导入的模块路径可以是一个表达式，因为它使用的是 require()方法；而 ES6 Modules 只能是字符串
   CommonJS this 指向当前模块，ES6 Modules this 指向 undefined
   且 ES6 Modules 中没有这些顶层变量：arguments、require、module、exports、**filename、**dirname

3. call 和 apply 的区别是什么，哪个性能更好一些
   call 比 apply 的性能要好，平常可以多用 call, call 传入参数的格式正是内部所需要的格式，参考 call 和 apply 的性能对比
   看到算法步骤中，apply 多了 CreateListFromArrayLike 的调用，其他的操作几乎是一样的（甚至 apply 仍然多了点操作）。从草案的算法描述来看，call 性能 > apply 性能。

4. 箭头函数与普通函数（function）的区别是什么？构造函数（function）可以使用 new 生成实例，那么箭头函数可以吗？为什么

5. 为什么普通 for 循环的性能远远高于 forEach 的性能，请解释其中的原因
   for 循环没有任何额外的函数调用栈和上下文；
   forEach 函数签名实际上是 array.forEach(function(currentValue, index, arr), thisValue)
   它不是普通的 for 循环的语法糖，还有诸多参数和上下文需要在执行的时候考虑进来，这里可能拖慢性能；
6. 数组里面有 10 万个数据，取第一个元素和第 10 万个元素的时间相差多少
   数组可以直接根据索引取的对应的元素，所以不管取哪个位置的元素的时间复杂度都是 O(1)
   得出结论：消耗时间几乎一致，差异可以忽略不计
7. null 和 undefined 的区别
   null 表示一个"无"的对象，也就是该处不应该有值；而 undefined 表示未定义。
   在转换为数字时结果不同，Number(null)为 0，而 undefined 为 NaN。

   null：
   作为函数的参数，表示该函数的参数不是对象
   作为对象原型链的终点

   undefined:
   变量被声明了，但没有赋值时，就等于 undefined
   调用函数时，应该提供的参数没有提供，该参数等于 undefined
   对象没有赋值属性，该属性的值为 undefined
   函数没有返回值时，默认返回 undefined

8. typeof 为什么对 null 错误的显示
   这只是 JS 存在的一个悠久 Bug。在 JS 的最初版本中使用的是 32 位系统，为了性能考虑使用低位存储变量的类型信息，000 开头代表是对象然而 null 表示为全零，所以将它错误的判断为 object 。
9. 一句话描述一下 this
   对于函数而言指向最后调用函数的那个对象，是函数运行时内部自动生成的一个内部对象，只能在函数内部使用；对于全局来说，this 指向 window。
10. 使用箭头函数时需要注意什么？
    不能用于构造函数
    不要用于事件绑定

```JS
const btn = document.getElementById('btn');
btn.addEventListener('click', function() {
  console.log(this) // <button id="btn">按钮</button>
})

const btn = document.getElementById('btn');
btn.addEventListener('click', () => {
  console.log(this) // window
})

```

11. 闭包产生的本质就是，当前环境中存在指向父级作用域的引用。
12. Array(3)和 Array(3, 4)的区别？

```JS
console.log(Array(3)) // [empty x 3]
console.log(Array(3, 4)) // [3, 4]
console.log(Array.of(3, 4)) // [3, 4]
```

13. Proxy get 的第三个参数 receiver(接受者)
    如果 target 对象中指定了 getter，receiver 则为 getter 调用时的 this 值。

```JS
var obj = {
  fn: function () {
    console.log('lindaidai')
  }
}
var obj1 = new Proxy(obj, {
  get (target, key, receiver) {
    console.log(receiver === obj1) // true
    console.log(receiver === target) // false
    return target[key]
  }
})
obj1.fn()

```

14. 使用 delete 删除数组元素，其长度是不会改变的。
    把数组理解为是一个特殊的对象，其中的每一项转换为对象的伪代码为：

```JS
key: value
// 对应:
0: 1
1: 2
2: 3
length: 3

```

delete 操作符删除一个数组元素时，相当于移除了数组中的一个属性，被删除的元素已经不再属于该数组。但是这种改变并不会影响数组的 length 属性。 15. 15. 8

15. “自动插入分号”，简称 ASI (Automatic Semicolon Insertion)

```JS
function getName () {
  return
  {
    name: 'LinDaiDai'
  }
}
console.log(getName())

相当于
function getName () {
  return;
  {
    name: 'LinDaiDai'
  }
}
console.log(getName())


因此最终的结果也就是undefined。
```

16. 按位取反，为什么~2 = -3
    计算机中用补码存储数据
    原码：最高位为符号位，剩余位（数据位）为 x 的绝对值。
    反码：如果 x 为正数，则与原码相同；如果 x 为负数，符号位保持不变，数据位取反。
    补码：如果 x 为正数，则与原码相同；**如果 x 为负数，符号位保持不变，数据位取反，然后加 1（若符号位有进位，则舍弃进位）。**

    1. 0000 0010
    2. 1111 1101
    3. 1000 0011 `2->3 补码转源码`
    4. -3

17. [,,,]的长度
    3
    尾后逗号的概念
18. 如何判断当前脚本运行在浏览器还是 node 环境中？
19. 如何判断一个对象是否为空对象 Object.keys()：
20. 闭包的内存如何释放

```JS
const fn = function() {
    let num = 0;

    return function() {
        return num += 1;
    }
}

let f1 = fn();

f1(); // 1
f1(); // 2
f1(); // 3

f1 = null;  // 释放

f1 = fn();

f1(); // 1
f1(); // 2
f1(); // 3
```

21. hooks 两条规范
    1. 只能用于 React 函数组件和自定义 hook 中
    2. 只能用于顶层代码，不能用于分支语句
       eslint 插件 eslint-plugin-react-hooks
       为什么？
       内部实现是数组
22. React Hooks 的坑
    1. useState 初始化值只初始化一次(参考实现)
    2. useEffect 内部不能修改 state
    3. useEffect 依赖不能为非 state 的引用类型，否则死循环
23. 原型链上接近神的男人->Object.prototype
24. 原型链判断

```JS
Object.prototype.__proto__; //null
Function.prototype.__proto__; //Object.prototype
Object.__proto__; //Function.prototype
Object instanceof Function; //true
Function instanceof Object; //true
Function.prototype === Function.__proto__; //true


```

25. 介绍一下`__proto__`和 prototype
    `__proto__`是每个对象都有的属性
    prototype 是函数才有的属性
26. use strict 效果?
    严格模式指定一个 script 标签内的代码在严格条件下执行。
    “use strict” 指令只允许出现在脚本或函数的开头。
    严格模式下不允许使用未声明的变量
    函数参数不能重名
    函数内 this 默认不再指向 window，默认为 undefined。
    delete 删除不可配置的属性报错
27. 什么是闭包
    闭包是函数和声明该函数的词法环境的组合。也就是说，它是一个内部函数，可以访问外部函数或封闭函数的变量。
28. null 与 undefined
    null:指示变量没有值
    undefined:表示变量本身不存在

29. 什么是 typed Array
    typed Array 是 ECMAScript 6 API 中用于**处理二进制数据的类数组对象**,
30. 什么是 ArrayBuffer
    ArrayBuffer 对象用于表示一般的、固定长度的**原始二进制数据缓冲区**

```JS
let buffer = new ArrayBuffer(16); // create a buffer of length 16
alert(buffer.byteLength); // 16

为了操作 ArrayBuffer，我们需要使用一个“ view”对象。

//Create a DataView referring to the buffer
let view = new DataView(buffer);
```

31. 是否所有对象都有原型
    用户创建的基本对象或使用 new 关键字创建的对象没有。
32. 如果我们把两个数组相加会发生什么
    如果将两个数组相加，它会将它们都转换为字符串并连接起来

```JS
console.log(['a'] + ['b']);  // "ab"
console.log([] + []); // ""
console.log(![] + []); // "false", because ![] returns false.
```

33. 如何创建复制到剪贴板按钮

```JS
document.querySelector("#copy-button").onclick = function() {
  // Select the content
  document.querySelector("#copy-input").select();
  // Copy to the clipboard
  document.execCommand('copy');
};

```

34. 如何在网页上禁用右键

```HTML
<body oncontextmenu="return false;">
```

35. shim 和 polyfill 的区别是什么
    shim 是一个库，它只使用旧环境的方法为旧环境带来新的 API。它不一定局限于 web 应用程序。
    而 polyfill 是一段代码(或插件) ，它提供了开发人员希望**浏览器**本身提供的技术。
36. 什么是尾递归(Tail Call)
    当函数调用是尾部调用时，程序或代码**不会为递归创建额外的堆栈帧**。

```JS
function factorial(n, acc = 1) {
  if (n === 0) {
    return acc
  }
  return factorial(n - 1, n * acc)
}
console.log(factorial(5)); //120

```

37. 如何检测一个函数是否被调用为构造函数

```JS
function Myfunc() {
   if (new.target) {
      console.log('called with new');
   } else {
      console.log('not called with new');
   }
}

new Myfunc(); // called with new
Myfunc(); // not called with new
Myfunc.call({}); not called with new

```

38. isNaN 和 number.isNaN 有什么区别？
    全局函数 isNaN 先将参数转换为 Number 如果结果值为 NaN，则返回 true。

```JS
isNaN(‘hello’);   // true
Number.isNaN('hello'); // false
```

建议用 Number.isNaN

// TODO

39. js 垃圾回收
40. JavaScript 中数组是如何存储的？**快速组/慢数组。**
    https://blog.csdn.net/wanderlustLee/article/details/100929118
    同种类型数据的数组分配连续的内存空间
    存在非同种类型数据的数组使用哈希映射分配内存空间

    `快数组`:连续的存储空间；数组长度可变，是内部通过扩容和收缩机制实现，类似 Java 中的 ArrayList 扩容形式，达到一定的阈值则拷贝内存到一个更大的空间中。**新容量 = 旧容量 + 50% + 16**
    `慢数组`:是一个以数字为键的 HashTable，他不用开辟大块的连续空间从而节省内存

    `快数组慢数组何时变换？`

    - 如果不指定容量新创建的数组都是快数组形式的，如果指定容量，预分配数组长度小于等于 1024，也是以快数组形式存放，如果大于 1024，就会使用哈希表来存放。
    - 对数组赋值时使用远超当前数组的容量（超过的容量由 kMaxGap 决定，为 1024），这样会出现大量空洞，这时候要对数组分配大量空间则将可能造成存储空间的浪费，为了空间的优化，**会转化为慢数组。**

    在遍历方面，快数组比慢数组快了很多。快数组就是以空间换时间的方式，申请了大块连续内存，提高效率。慢数组以时间换空间，不必申请连续的空间，节省了内存，但需要付出效率变差的代价。

    **扩展：ArrayBuffer**
    在后面，JS 在 ES6 也推出了**分配连续内存的数组**，这就 ArrayBuffer。ArrayBuffer 会从内存中申请设定的二进制大小的空间，但是并不能直接操作 ArrayBuffer，需要通过 TypedArray/DataView 构建一个视图，通过视图来操作这个内存。

    ```TS
    interface ArrayBuffer {
        /**
         * Read-only. The length of the ArrayBuffer (in bytes).
        */
        readonly byteLength: number;

        /**
         * Returns a section of an ArrayBuffer.
        */
        slice(begin: number, end?: number): ArrayBuffer;
    }
    ```

    ```TS
    这行代码就申请了1kb的内存区域。
    var bf = new ArrayBuffer(1024);
    var a = new Uint8Array(bf);  // 数组的长度也是1024
    var b = new Int32Array(bf);  // 数组的长度是1024/4=256

    type Test = Exclude<keyof Array<any>, keyof Uint8Array>
    可以看到 TypedArray确实是长度固定的
    type Test = typeof Symbol.unscopables | "pop" | "push" | "concat" | "shift" | "splice" | "unshift" | "flatMap" | "flat"
    ```

41. 说说对原生 JavaScript 的理解？
    多范式语言
    事件驱动 / 全是异步 IO
42. 高性能的 JavaScript 开发在语法层面你觉得有哪些可以提升性能？
    假如支持海象表达式的话...会提升一些
43. 在 JavaScript 中如何实现对象的私有属性?
    IIFE 实现
44. async / await 和 Promise 的区别?
    await 会等待异步代码执行，**会阻塞代码**（使用时要考虑性能）
    async / await 在调试方面会更加方便

https://github.com/lin-123/javascript/blob/cn/README.md

45. aribnb
    由于 Symbols 和 BigInts 不能被正确的 polyfill。所以不应在不能原生支持这些类型的环境或浏览器中使用他们。

    将你的所有缩写放在对象声明的前面。因为这样能更方便地知道有哪些属性用了缩写。

    对象浅拷贝时，更推荐使用扩展运算符（即 ... 运算符），而不是 Object.assign

    当需要动态生成字符串时，使用模板字符串而不是字符串拼接。

    不要对参数重新赋值:参数重新赋值会导致意外行为，尤其是对 arguments。这也会导致优化问题，特别是在 V8 引擎里。

    **用 Number.isNaN 代替全局的 isNaN**
    **用 Number.isFinite 代替 isFinite**
    全局 `isNaN 强制把非数字转成数字`， 然后对于任何强转后为 NaN 的变量都返回 true

    ```JS
    Number.isNaN(Number('1.2.3')) // true
    isNaN('1.2.3')  // true
    ```

    把 const 和 let 分别放一起。(先 const 后 let)

    除非外部库或框架需要使用特定的非静态方法，否则类方法应该使用 this 或被写成静态方法(**不用 this 就用静态**)

    一个空的构造函数或只是代表父类的构造函数是不需要写的。

    **多个返回值用对象的解构，而不是数组解构。**

    不要使用链式声明变量;**链式声明变量会创建隐式全局变量**

    ```JS
        // bad
    (function example() {
      // JavaScript 将这一段解释为
      // let a = ( b = ( c = 1 ) );
      // let 只对变量 a 起作用; 变量 b 和 c 都变成了全局变量
      let a = b = c = 1;
    }());

    console.log(a); // undefined
    console.log(b); // 1
    console.log(c); // 1
    ```

    布尔值要用缩写，而字符串和数字要明确使用比较操作符。

    ```JS
      // bad
      if (isValid === true) {
        // ...
      }

      // good
      if (isValid) {
        // ...
      }

      // bad
      if (name) {
        // ...
      }

      // good
      if (name !== '') {
        // ...
      }

      // bad
      if (collection.length) {
        // ...
      }

      // good
      if (collection.length > 0) {
        // ...
      }
    ```

    用圆括号来组合操作符

    不要用选择操作符代替控制语句。

    ```JS
    // bad
    !isRunning && startRunning();

    // good
    if (!isRunning) {
      startRunning();
    }
    ```

    用移位运算要小心。数字是用 64-位表示的，但移位运算常常返回的是 32 为整形

    **不要用前置或后置下划线**

46. ArrayBuffer 对象
    ArrayBuffer 对象、TypedArray 视图和 DataView 视图是 JavaScript **操作二进制数据**的一个接口
    这个接口的原始设计目的，与 WebGL 项目有关。所谓 WebGL，就是指浏览器与显卡之间的通信接口，为了满足 JavaScript 与显卡之间大量的、实时的数据交换，它们之间的数据通信必须是二进制的，而不能是传统的文本格式
    存在一种机制，可以像 C 语言那样，直接操作字节，将 4 个字节的 32 位整数，以二进制形式原封不动地送入显卡，脚本的性能就会大幅提升。
    **二进制数组由三类对象组成**
    （1）ArrayBuffer 对象：代表内存之中的`一段二进制数据`，`可以通过“视图”进行操作`。“视图”部署了数组接口，这意味着，可以用数组的方法操作内存。

    ArrayBuffer 和 Buffer 有何区别？
    ArrayBuffer[1] 对象用来表示通用的、固定长度的原始二进制数据缓冲区，是一个字节数组，可读但不可直接写。

    Buffer[2] 是 Node.JS 中用于操作 ArrayBuffer 的视图，是 `TypedArray[3] 的一种`。

    Blob 是一种二进制对象(包括字符，文件等等)，es6 对其进行了补充
    File 是基于 Blob 的一种二进制文件对象,扩展了 Blob，es6 同样进行了补充
    ArrayBuffer 是 ES6 新引入的二进制缓冲区
    Buffer 是 Nodejs 内置的二进制缓冲区，Buffer 相当于 ES6 中 Uint8Array(属于 TypedArray)的一种扩展

    （2）TypedArray 视图：共包括 9 种类型的视图，比如 Uint8Array（无符号 8 位整数）数组视图, Int16Array（16 位整数）数组视图, Float32Array（32 位浮点数）数组视图等等。

    （3）DataView 视图：**可以自定义复合格式**的视图，比如第一个字节是 Uint8（无符号 8 位整数）、第二、三个字节是 Int16（16 位整数）、第四个字节开始是 Float32（32 位浮点数）等等，此外还可以自定义字节序。

    `简单说，ArrayBuffer 对象代表原始的二进制数据，TypedArray 视图用来读写简单类型的二进制数据，DataView 视图用来读写复杂类型的二进制数据。`

47. 多线程通信:共享内存 SharedArrayBuffer
    JavaScript 是单线程的，Web worker 引入了多线程：主线程用来与用户互动，Worker 线程用来承担计算任务。每个线程的数据都是隔离的，通过 postMessage()通信
    ES2017 引入 SharedArrayBuffer，允许 Worker 线程与主线程共享同一块内存。SharedArrayBuffer 的 API 与 ArrayBuffer 一模一样，唯一的区别是后者无法共享数据。

    ```JS
    // 主线程

    // 新建 1KB 共享内存
    const sharedBuffer = new SharedArrayBuffer(1024);

    // 主线程将共享内存的地址发送出去
    w.postMessage(sharedBuffer);

    // 在共享内存上建立视图，供写入数据
    const sharedArray = new Int32Array(sharedBuffer);
    ```

48. Atomics 对象
    SharedArrayBuffer API 提供 Atomics 对象，保证所有共享内存的操作都是“原子性”的，并且可以在所有线程内同步
    Atomics.store()，Atomics.load()

    store()方法用来向共享内存写入数据，load()方法用来从共享内存读出数据。比起直接的读写操作，它们的好处是保证了读写操作的原子性。

    Atomics.wait()，Atomics.notify()

49. 类数组有哪些
    `有 length 属性的对象就叫类数组`(见 TS ArrayLike 接口)

    ```JS
    const arrayLike = { length: 2, a: 12, 0: 9 }
    console.log(Array.from(arrayLike))

    [ 9, undefined ]
    ```

    arguments NodeList HTMLCollection

50. 如何理解声明式(declarative)与命令式(imperative)
    声明式(declarative)是结果导向的，命令式(imperative)是过程导向的。它们都有自己适用的场景和局限，于是现实中的编程语言常常都有两者的身影。
    命令式编程（Imperative）：详细的命令机器怎么（How）去处理一件事情以达到你想要的结果（What）；
    声明式编程（ Declarative）：只告诉你想要的结果（What），机器自己摸索过程（How）
    React 是声明式的

命令式编程（Imperative)会一步一步的告诉程序该怎么运行

```JS
var array = [1,2,3]
var output = []
for(var i = 0; i < array.length; i++)
{
  var tmp = array[i] * 2
  output.push (tmp)
}
console.log (output) //=> [2,4,6]
```

如果使用声明式编程（ Declarative）则会是这样

```JS
var array = [1,2,3,]
var output = array.map (function (n)
{
  return n * 2
})
console.log (output) //=> [2,4,6]
```

通常情况下我们常用的大部分编程语言：c，java，c++等都是命令式编程语言。而像正则表达式（regular expressions）或者逻辑语言（Prolog）则为声明式语言。
声明式语言包包括数据库查询语言（SQL，XQuery），正则表达式，逻辑编程，`函数式编程`和组态管理系统。

51. 说出结果

```JS
function test(person) {
  person.age = 26
  person = {
    name: 'hzj',
    age: 18
  }
  return person
}
const p1 = {
  name: 'fyq',
  age: 19
}
const p2 = test(p1)
console.log(p1) // -> ?
console.log(p2) // -> ?

p1：{name: “fyq”, age: 26}
p2：{name: “hzj”, age: 18}

函数传参的时候传递的是`对象在堆中的内存地址值`，test函数中的实参person是p1对象的内存地址，通过调用person.age = 26确实改变了p1的值，但随后person变成了另一块内存空间的地址，并且在最后将这另外一份内存空间的地址返回，赋给了p2。
```

52. null 是对象吗？为什么？
    虽然 typeof null 会输出 object，但是这只是 JS 存在的一个悠久 Bug。在 `JS 的最初版本中使用的是 32 位系统`，为了性能考虑使用`低位存储变量的类型信`息，`000 开头代表是对象然而 null 表示为全零`，所以将它错误的判断为 object 。

53. '1'.toString()为什么可以调用？

```JS
其实在这个语句运行的过程中做了这样几件事情：
var s = new Object('1');
s.toString();
s = null;

创建Object类实例; 调用实例方法;执行完方法立即销毁这个实例。
整个过程体现了`基本包装类型`的性质
```

54. instanceof 能否判断基本数据类型？
    一般不能。

```JS
class PrimitiveNumber {
  static [Symbol.hasInstance](x) {
    return typeof x === 'number'
  }
}
console.log(111 instanceof PrimitiveNumber) // true

```

55. [] == ![]结果是什么？为什么？
    解析:
    == 中，左右两边都需要转换为数字然后进行比较。
    []转换为数字为 0。
    ![] 首先是转换为布尔值，由于[]作为一个引用类型转换为布尔值为 true,
    因此![]为 false，进而在转换成数字，变为 0。
    0 == 0 ， 结果为 true
56. 对象转原始类型是根据什么流程运行的？
    如果 Symbol.toPrimitive()方法，优先调用再返回
    调用 valueOf()，如果转换为原始类型，则返回
    调用 toString()，如果转换为原始类型，则返回
    如果都没有返回原始类型，会报错

57. V8 执行代码过程

    1. 首先通过词法分析和语法分析生成 AST
    2. 将 AST 转换为字节码
    3. 由解释器逐行执行字节码(JIT)，遇到热点代码(HotSpot)启动编译器进行编译，生成对应的机器码并保存, 以优化执行效率。

    代码执行的时间越久，那么执行效率会越来越高，因为有越来越多的字节码被标记为热点代码，遇到它们时直接执行相应的机器码，**不用再次将转换为机器码**。

58. Pull vs Push
    Pull 和 Push 是数据生产者和数据的消费者两种不同的交流方式。
    每一个 JavaScript 函数都是一个 "拉" 体系，函数是数据的生产者，调用函数的代码通过 ''拉出" 一个单一的返回值来消费该数据。ES6 介绍了 iterator 迭代器 和 Generator 生成器 — 另一种 "拉" 体系，调用 iterator.next() 的代码是消费者，可从中拉取多个值。
    Promise(承诺) 是当今 JS 中最常见的 "推" 体系，一个 Promise (数据的生产者)发送一个 resolved value (成功状态的值)来执行一个回调(数据消费者)，但是不同于函数的地方的是：Promise 决定着何时数据才被推送至这个回调函数。
59. Observable vs Promise
    Observable（可观察对象）是基于推送（Push）运行时执行（lazy）的多值集合。
    方法|单值|多值
    -----|-----|-----
    pull|函数|遍历器
    push|Promise|Observable

    **Promise**
    返回单个值
    不可取消的

    **Observable**
    随着时间的推移发出多个值
    可以取消的
    支持 map、filter、reduce 等操作符 d
    延迟执行，当订阅的时候才会开始执行

60. JavaScript 里最大的安全的整数为什么是 2 的 53 次方减一？
    sign（S）：符号位，0 是正数，1 是负数 1

exponent（E）：指数位 11

fraction（F）：有效数字，IEEE754 规定，`在计算机内部保存有效数字时`，`默认第一位总是1`，所以舍去，只保留后面的部分。比如保存 1.01，只存 01，等读取的时候再把第一位 1 加上去。`所以52位有效数字实际可以存储53位`。 52
十进制 5，写成二进制是 101，相当于二进制 1.01× 2∧2，写成 64 位浮点数是
0 10000000001 0100000...0（F 剩下的全是 0，一共 52 位）

该二进制数换算为十进制数字即为 2∧52+2∧51+…+2∧1+2∧0 一个等比数列求和的计算~结果为 2∧53-1
所以，双精度浮点数可以表示的最大安全整数就是 2∧53-1 了！

61. JavaScript 里 Infinity 怎么算呢

62. encodeURIComponent() 函数会编码所有的字符。如果你想把 URI `当作请求参数传递，那么你可以使用本函数`
    如果你只是想编码一个带有特殊字符（比如中文）的 URI，这个 URI `用作请求地址，请使用 encodeURI 函数`。

```JS
// 原URI
var ftpUri = 'ftp://192.168.0.100/共享文件夹';

// 编码URI
var encodedFtpUri = encodeURI(ftpUri);
console.log(encodedFtpUri); // ftp://192.168.0.100/%E5%85%B1%E4%BA%AB%E6%96%87%E4%BB%B6%E5%A4%B9

// 原 URI 组件
var origin = 'ftp://192.168.0.100/共享文件夹';

// 编码 URI 组件
var encodedUri = encodeURIComponent(origin);
document.writeln(encodedUri); // ftp%3A%2F%2F192.168.0.100%2F%E5%85%B1%E4%BA%AB%E6%96%87%E4%BB%B6%E5%A4%B9
```

如果参数 encodedURIString 无效，将引发 URIError 错误。

63. Symbol.hasInstance 用于判断某对象是否为某构造器的实例。当其他对象使用 instanceof 运算符，判断是否为该对象的实例时，会调用这个方法。

```JS
class MyArray {
  static [Symbol.hasInstance]() {
    return Array.isArray(instance);
  }
}

[] instanceof new MyArray(); // true
```

64. Symbol 值只能通过 Symbol 函数生成
    注意，Symbol 函数前不能使用 new 命令实例化，
    这是因为生成的 Symbol 是一个原始类型的值，不是对象。
    基本上，它是一种类似于字符串的数据类型。
    使用 instanceof 检测实例与 Symbol 之间的关系没用。

```JS
const symbol = Symbol('foo');

console.log(symbol instanceof Symbol);
// false
```

65. 由于每一个 Symbol 值都是不相等的，这意味着 Symbol 值可以作为 标识符，用于对象的属性名，就能保证不会出现同名的属性。`利用 Symbol 值的唯一特性，作为类库某些对象的属性名，这样可以避免使用者命名冲突导致的覆盖问题`

```JS
let mySymbol = Symbol();
// 第二种写法 字面量
let b = {
  [mySymbol]: 'Hello!',
};

```

66. Symbol 值作为对象属性名不能用点运算符，因为会转为字符串

67.

```JS
String.fromCharCode(65, 66, 67);
// ABC
```

68. Promise.allSettled
    当您有多个彼此不依赖的异步任务成功完成时，或者您总是想知道每个 Promise 的结果时，通常使用它。

    应用场景：
    `同时上传多张图片，实现异步并发（例如使用阿里云 OSS 同时批量上传多张图片）`

69. Promise.any
    只要`其中的一个 promise 成功`，就返回那个已经成功的 promise
    如果可迭代对象中没有一个 promise 成功（即所有的 promises 都失败/拒绝），就返回一个失败的 promise 和 AggregateError 类型的实例，它是  Error  的一个子类，用于把单一的错误集合在一起
    Promise.any 应用场景
    从最快的服务器检索资源
    显示第一张已加载的图片（来自 MDN）
    `Promise.any vs Promise.race`
    Promise.any() ：关注于 Promise 是否已经成功
    Promise.race() ： 主要关注 Promise 是否已经解决，无论它是被解决还是被拒绝
70. Promise.allSettled() 与 Promise.all() 各自的适用场景
    Promise.allSettled()  更适合：

    彼此不依赖，其中任何一个被 reject ，对其它都没有影响
    期望知道每个 promise 的执行结果

    Promise.all()  更适合：

    彼此相互依赖，其中任何一个被 reject ，其它都失去了实际价值

71. `yield惰性求值`

```JS
function* gen() {
  yield 123 + 456;
}
```

上面代码中，yield 后面的表达式 123 + 456，不会立即求值，只会在 next 方法将指针移到这一句时，才会求值。

72. Generator 原型方法
    Generator.prototype.next
    Generator.prototype.return
    Generator.prototype.throw
    三者的作用都是让 Generator 函数恢复执行，`并且使用不同的语句替换 yield 表达式。`

```JS
const generator = function*(x, y) {
  let result = yield x + y;
  return result;
};

const gen = generator(1, 2);

gen.next(); // Object {value: 3, done: false}

next() 是将 yield 表达式替换成一个值。
gen.next(1); // Object {value: 1, done: true}

// 相当于将 let result = yield x + y
// 替换成 let result = 1;

throw() 是将 yield 表达式替换成一个 throw 语句。
gen.throw(new Error('出错了')); // Uncaught Error: 出错了
// 相当于将 let result = yield x + y
// 替换成 let result = throw(new Error('出错了'));

return() 是将 yield 表达式替换成一个 return 语句。
gen.return(2); // Object {value: 2, done: true}

// 相当于将 let result = yield x + y
// 替换成 let result = return 2;
```

73. Proxy

- 代理的引用上下文问题

```JS
const target = {
  foo: function () {
    console.log(this === proxy);
  },
};

const handler = {};

const proxy = new Proxy(target, handler);

console.log(target.foo());
// false
console.log(proxy.foo());
// true
```

一旦 proxy 代理 target.foo，后者内部的 this 就是指向 proxy，而不是 target。

- Proxy 与 Object.defineProperty
  Object.defineProperty 的三个主要问题：

  无法监听数组变化，Vue 通过 Hack 改写八种数组方法实现
  只能劫持对象的属性，因此对需要双向绑定的属性需要显式地定义
  必须深层遍历嵌套的对象

  与 Proxy 的区别：
  Proxy 可以直接监听数组的变化
  Proxy 可以直接监听对象而非属性
  Proxy 直接可以劫持整个对象，并返回一个新的对象，不管是操作便利程度还是底层功能上都远强于 Object.defineProperty
  Proxy 有多达 13 中拦截方法，不限于 apply、ownKeys、deleteProperty、has 等等是 Object.defineProperty 不具备的

```JS
const pipe = (value) => {
  const stack = [];
  const proxy = new Proxy(
    {},
    {
      get(target, prop) {
        if (prop === 'execute') {
          return stack.reduce(function (val, fn) {
            return fn(val);
          }, value);
        }
        stack.push(window[porp]);
        return proxy;
      },
    }
  );
  return proxy;
};

const double = (n) => n * 2;
const pow = (n) => n * n;

console.log(pipe(3).double.pow.execute);
```
