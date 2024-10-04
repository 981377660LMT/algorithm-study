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

hopocroft 算法**合并等价状态**

实际上，NFA 转 DFA 是一个繁琐的过程，如果正则采用 DFA 引擎，势必会消耗部分性能在 NFA 的转换上，而这个转换的效益很多时候远不比直接用使用 NFA 高效，同时 DFA 相对来说没有 NFA 直观，可操作空间也要比 NFA 少，所以大多数语言的采用 NFA 作为正则的引擎。

回溯灾难
正则采用的是 NFA 引擎，那么我们就必须面对它的不确定性，体现在正则上就是 回溯 的发生。
回溯失控，举个简单的例子，假设我们有一个字符串 aaa，以及正则表达式 a\*b

---

javascript parser generator 选型
https://segmentfault.com/a/1190000038554196

解析器:两个流派，自底向上和自顶向下

- codegen 的 parser
  https://pegjs.org/ （最好用的）
  https://gerhobbelt.github.io/jison/docs/

- Parsec
  自顶向下分析算法的一种
  parsec 支持 cfg，所以可以解析语法，跟手写递归下降是等价的
  手写递归下降要手写很多东西，parsec 只要组合一些常见的函数，所以写原型快，另外零依赖
  只需要编写每个语法单元的 parser，然后利用 parsec 库组合起来，就是一个完整的语法解析程序
  性能不是最好，所以仅限于原型
  parser combinator 和 parser generator 的区别，实际上就是将算法的状态表示为一个整数（Yacc/Bison 的做法）和一个类（ANTLR 的做法）的分别，看上去写法完全不同，然而写法并不决定表达能力。虽然我必须承认，从工程角度上，combinator 不需要单独写语法定义文件，调试起来确实比 Yacc/Bison 舒服很多。
  https://www.zhihu.com/question/35778359
  Parser Combinator 在语法解析的当中处于怎样的位置?
- `上下文无关文法`：就是说这个文法中**所有的产生式左边只有一个非终结符**
  S -> aSb

  S -> ab

  这个文法有两个产生式，每个产生式左边只有一个非终结符 S，这就是上下文无关文法，因为你只要找到符合产生式右边的串，就可以把它归约为对应的非终结符。

  aSb -> aaSbb

  S -> ab

  这就是上下文相关文法，因为它的第一个产生式左边有不止一个符号，所以你在匹配这个**产生式中的 S 的时候必需确保这个 S 有正确的“上下文”，也就是左边的 a 和右边的 b**，所以叫上下文相关文法。

---
