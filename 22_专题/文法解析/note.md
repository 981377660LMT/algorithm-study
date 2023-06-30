TODO
https://github.com/spaghetti-source/algorithm/tree/master/string

DFA/NFA
https://github.com/hos-lyric/libra/blob/60b8b56ecae5860f81d75de28510d94336f5dad9/string/finite_automaton.cpp#L191
https://nxtech.gitbook.io/regexp/advanced/dfa-nfa

DFA 转正则表达式的方法
https://github.com/tdzl2003/leetcode_live/blob/master/regexp/20220102.md

---

https://ewind.us/2016/toy-html-parser/
一般一个完整的编译过程由三步组成：词法分析、语法分析和语义分析。
这三个流程各对应一个模块：词法分析器、语法分析器和语义计算模块。

语法分析器中的 LL 算法和 LR 算法:
LL 算法递归向下地处理代码语句，而 LR 算法则是自底向上地归约词法元素。

---

[正则表达式是如何运作的？](https://zhuanlan.zhihu.com/p/608908724)
[一句话总结 NFA 转 DFA 算法](https://zhuanlan.zhihu.com/p/39452141)
[DFA 最小化算法-三步](https://www.zhihu.com/question/39767421/answer/338794446)

DFA 可以看作是一种特殊的 NFA，其最大的特点就是确定性，即输入一个字符，一定会转换到确定的状态，而不会有其他的可能。

正则表达式转 NFA 主要基于以下 3 条规则
连接 选择 重复

Thompson 算法 中两种最基本的单元（或者说两种最基本的 NFA）:表示经过字符 a 过渡的下一个状态以及不需要任何字符输入的 ε 转换 过渡到下一个状态。
Thompson 提出的 Thompson 算法。其思想主要就是通过简单的表达式组合成复杂的表达式。

正则转 nfa 转 dfa 的三步

1. 拆分组合表达式，正则转 nfa

2. 子集构造算法，nfa 转 dfa
   如果状态 B, C, D 可以由同一(或多)个状态 A 加同一个 transition 得到, 那么状态 B, C, D 等价.
   Subset Construction 的完整过程就是沿着起始状态, 做一把 BFS, 通过 Duck Test 不断把可以折叠的鸟合并成一只鸭子, 再从这些鸭子节点上继续做 BFS, 直到所有的鸟都是鸭子.

所以为什么这个算法一定可以在有限时间内完成? 我也学到了一个新名词叫做不动点(Fixed-Point)运算. 因为假设 NFA 有 n 个状态, m 种 transition, 那么可以做的 DFA 折叠操作理论最大次数是 n 的组合*m+(n-1)的组合 m...直至 1*m, 次数有限, 自然不会停机.

3. dfa 最小化

hopocroft 算法合并等价状态

实际上，NFA 转 DFA 是一个繁琐的过程，如果正则采用 DFA 引擎，势必会消耗部分性能在 NFA 的转换上，而这个转换的效益很多时候远不比直接用使用 NFA 高效，同时 DFA 相对来说没有 NFA 直观，可操作空间也要比 NFA 少，所以大多数语言的采用 NFA 作为正则的引擎。

回溯灾难
正则采用的是 NFA 引擎，那么我们就必须面对它的不确定性，体现在正则上就是 回溯 的发生。
回溯失控，举个简单的例子，假设我们有一个字符串 aaa，以及正则表达式 a\*b
