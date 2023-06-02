https://nyaannyaan.github.io/library/string/suffix-automaton.hpp
https://baobaobear.github.io/post/20200220-sam/
https://w.atwiki.jp/uwicoder/pages/2842.html

https://ouuan.github.io/post/%E5%90%8E%E7%BC%80%E8%87%AA%E5%8A%A8%E6%9C%BAsam%E5%AD%A6%E4%B9%A0%E7%AC%94%E8%AE%B0/
https://etaoinwu.com/blog/%E6%84%9F%E6%80%A7%E7%90%86%E8%A7%A3-sam/
—个确定有限状态自动机(DFA)由以下五部分构成:

1. 字符集(∑)，该自动机只能输入这些字符。
2. 状态集合(Q)。如果把一个 DFA 看成一张有向图，那么 DFA 中的状态就相当于图上的顶点。
3. 起始状态(start)， start ∈Q，是一个特殊的状态。起始状态一般用 s 表示，为了避免混淆，本文中使用 start。
4. 接受状态集合(F)， F ∈Q，是一堆特殊的状态。
5. 转移函数(6)，6 是一个接受两个参数返回一个值的函数，其中第一个参数和返回值都是一个状态，第二个参数是字符集中的一个字符。如果把一个 DFA 看成一张有向图，那么 DFA 中的转移函数就相当于顶点间的边，而每条边上都有一个字符。

SAM 上的每一个状态去表示一个集合等价类，转移函数也相应地更改为对应的等价类
这样的 SAM，从起始状态到某个状态可能有多条路径，每条路径都对应一个字符串，那么我们称这个状态 对应 着这些字符串。

TODO 没理解
