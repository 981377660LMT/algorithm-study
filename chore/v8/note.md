## 如何利用 v8 来优化 js 代码

1. 在 V8 出现之前，所有的 JavaScript 虚拟机所采用的都是解释执行的方式，这是
   JavaScript 执行速度过慢的一个主要原因。而 V8 率先引入了即时编译（JIT）的双轮驱动
   的设计，这是一种权衡策略，混合编译执行和解释执行这两种手段，给 JavaScript 的执行
   速度带来了极大的提升。
2. JavaScript 是一门非常优秀的语言，特别是“原型继承机制”和“函数是一等公民”这两个设计。
3. V8 并没有采用某种单一的技术，而是混合编译执行和解释执行这两种手段，我们
   **把这种混合使用编译器和解释器的技术称为 JIT（Just In Time）技术。**
4. **频繁使用大的临时变量**，导致了新生代空间很快被装满，从而频繁触发垃圾回收

```JS
function strToArray(str) {
   let i = 0
   const len = str.length
   let arr = new Uint16Array(str.length)
   for (; i < len; ++i) {
      arr[i] = str.charCodeAt(i)
   }
   return arr;
}

function foo() {
   let i = 0
   let str = 'test V8 GC'
   while (i++ < 1e5) {
      strToArray(str);
   }
}

foo()
```

这段代码就会频繁创建临时变量，这种方式很快就会造成新生代内存内装满，从而频繁触发
垃圾回收。**为了解决频繁的垃圾回收的问题，你可以考虑将这些临时变量设置为全局变量。**
