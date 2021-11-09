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
