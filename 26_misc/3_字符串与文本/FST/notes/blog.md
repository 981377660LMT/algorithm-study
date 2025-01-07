1. [Index 1,600,000,000 Keys with Automata and Rust](https://burntsushi.net/transducers/)
   有限状态机还可以用于紧凑地表示有序集合或字符串映射，这些集合或映射可以被非常快速地搜索。
   本文中介绍的技术也是 Lucene 表示其倒排索引一部分的方式。
   [Using Finite State Transducers in Lucene](https://blog.mikemccandless.com/2010/12/using-finite-state-transducers-in.html)

2. [Lucene 倒排索引之 FST](https://zhuanlan.zhihu.com/p/671225495)
3. [字典数据结构-FST(Finite State Transducers)](https://zhuanlan.zhihu.com/p/366849553)
   - 一般SortedDict 倒排索引的方式，需要完整存储每一个term。term 数目多达上千万时，占用的内存将不可接受
