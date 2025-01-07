1. [Index 1,600,000,000 Keys with Automata and Rust](https://burntsushi.net/transducers/)
   有限状态机还可以用于紧凑地表示有序集合或字符串映射，这些集合或映射可以被非常快速地搜索。
   本文中介绍的技术也是 Lucene 表示其倒排索引一部分的方式。
   [Using Finite State Transducers in Lucene](https://blog.mikemccandless.com/2010/12/using-finite-state-transducers-in.html)

2. [Lucene 倒排索引之 FST](https://zhuanlan.zhihu.com/p/671225495)

   XX

3. [字典数据结构-FST(Finite State Transducers)](https://zhuanlan.zhihu.com/p/366849553)

   - SortedDict 倒排索引的方式，需要完整存储每一个 term。term 数目多达上千万时，占用的内存将不可接受
   - 常用字典数据结构：有序数组、HashMap、SkipList、Trie、DoubleArrayTrie、FST
   - lucene 从 4 开始大量使用的数据结构是 FST（Finite State Transducer）。FST 有两个优点：1）空间占用小。通过对词典中单词前缀和后缀的重复利用，压缩了存储空间；2）查询速度快。O(len(str))的查询时间复杂度。

4. [关于 Lucene 的词典 FST 深入剖析](https://www.shenyanchao.cn/blog/2018/12/04/lucene-fst/)

5. [trie、FSA、FST（转）](https://www.cnblogs.com/ajianbeyourself/p/11259984.html)
   trie，FSA，FST 都是用来解决有限状态机的存储，trie 是树，它进一步演化为 FSA 和 FST，这两者是图

   - FSA 和 trie 的区别：
     trie 树只共享前缀，而 FSA 可以共享前缀和后缀
     https://steflerjiang.github.io/2017/03/18/%E4%BD%BF%E7%94%A8Automata%E6%9D%A5%E7%B4%A2%E5%BC%951-600-000-000%E4%B8%AA%E9%94%AE/
     FSA 一般的构建方法是，DFA 最小化

   - FSA 和 FST 的区别：
     FST 和 FSA 很像，但是对于一个 key，FSA 只回答了”yer or no”，FST 不仅回答”yes or no”，还好返回和这个 key 相关的一个值。
     http://examples.mikemccandless.com/fst.py?terms=mop%2F0%0D%0Amoth%2F1%0D%0Apop%2F2%0D%0Astar%2F3%0D%0Astop%2F4%0D%0Atop%2F5%0D%0A&cmd=Build+it%21

6. [Finite state machines as data structure for representing](https://news.ycombinator.com/item?id=10551280)
