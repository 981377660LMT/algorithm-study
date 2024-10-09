整体的步骤是三步：

一，先把正规式转换为 NFA（非确定有穷自动机）(Thompson 算法)；,

二，在把 NFA 通过“子集构造法”转化为 DFA，

三，在把 DFA 通过“分割法”进行最小化(Hopcroft 算法)。

---

1. 正则表达式
   https://taodaling.github.io/blog/2020/08/04/%E6%AD%A3%E5%88%99%E8%A1%A8%E8%BE%BE%E5%BC%8F/#heading-%E8%87%AA%E5%8A%A8%E6%9C%BA
2. 算法实现
   https://github.com/spaghetti-source/algorithm/tree/master/string

   DFA/NFA
   https://github.com/hos-lyric/libra/blob/60b8b56ecae5860f81d75de28510d94336f5dad9/string/finite_automaton.cpp#L191
   https://nxtech.gitbook.io/regexp/advanced/dfa-nfa

3. DFA 转正则表达式的方法
   https://github.com/tdzl2003/leetcode_live/blob/master/regexp/20220102.md

4. 其他
   [正则表达式是如何运作的？](https://zhuanlan.zhihu.com/p/608908724)
   [一句话总结 NFA 转 DFA 算法](https://zhuanlan.zhihu.com/p/39452141)
   [DFA 最小化算法-三步](https://www.zhihu.com/question/39767421/answer/338794446)

5. javascript parser generator 选型
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
