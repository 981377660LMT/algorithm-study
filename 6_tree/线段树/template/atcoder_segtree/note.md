使用案例
https://zhuanlan.zhihu.com/p/459880950
https://zhuanlan.zhihu.com/p/459579152
https://zhuanlan.zhihu.com/p/459679512

1. 单点修改，区间查询
   普通的线段树
2. 区间修改，单点查询
   DualSegmentTree/SegmentTreeDual
3. 区间修改，区间查询
   LazySegmentTree/SegmentTreeLazy

   TODO:优化 e 和 id 的处理.e 和 id 必须是函数,但是比较时可以不用函数,而是与幺元对象比较(少调用一次函数).

---

https://panzhongxian.cn/cn/2021/04/02-inside-the-v8-engine-5-tips-on-how-to-write-optimized-code/

- 内联 (Inlining)
- 隐藏类 (Hidden class)
  由于使用字典来查找对象属性在内存中的位置非常低效，因此 V8 使用了一种技巧：隐藏类。隐藏类的工作方式类似于 Java 之类的语言中使用的固定对象布局（类），不同之处在于隐藏类是在运行时创建的。
- 编写优化 JavaScript 的 5 个技巧

1.  对象属性的顺序：始终以相同的顺序实例化对象属性，以便可以共享隐藏的类以及随后优化的代码。
2.  动态属性：实例化后向对象添加属性将强制更改隐藏类，已经为上一个隐藏类优化过的方法都会因此而变慢。替代措施是在构造函数中分配对象的所有属性。
3.  方法：由于内联缓存 (inline caching) 的原因，重复执行相同方法将比仅执行一次许多不同方法的代码运行更快。
4.  数组：避免键不是递增数字的稀疏数组。里面没有所有元素的稀疏数组是一个“哈希表”。访问这样数组中的元素耗费更大。另外，避免预先分配大数组，在使用中逐渐增长更好。最后，不要删除数组中的元素，它使键更稀疏。
5.  标记值 (tagged values)：V8 使用 32 位表示对象和数字。它使用 32 bit 中的一个 bit 来标记变量是对象（flag = 1），还是的整数（flag = 0），这个整数被称为”小整数(SMI, SMall Integer)“，因为它只有 31 bit。如果数值大于 31 位，**则 V8 会将数字装箱 (box)，将其变成 double 并创建一个新对象以将数字放入其中**。尽可能使用 31 bit 带符号的数字，以避免对 JS 对象进行昂贵的装箱操作。
