1. 双数组 Trie 以及 双数组 AC 自动机(HANLP 分词中使用)
   https://www.hankcs.com/program/java/triedoublearraytriejava.html
   https://zhuanlan.zhihu.com/p/527783449
2. 正则转 NFA,NFA 转 DFA, DFA 最小化
   后缀自动机是后缀字符串建立 Trie 的最小 DFA
3. 后缀自动机，广义后缀自动机，编辑距离自动机
   https://www.zhihu.com/people/ideas-yd/posts
   https://julesjacobs.com/2015/06/17/disqus-levenshtein-simple-and-fast.html
4. HanLP《自然语言处理入门》笔记
   https://www.zhihu.com/column/young-doctor

5. Trie 的变体
   https://zhuanlan.zhihu.com/p/527783449

   - Patricia Trie （Compact Trie）
     Patricia Trie 也叫 Compact Trie 是一种简单针对单链路径压缩的 Trie。压缩程度不高，但是很简单。
     https://zhuanlan.zhihu.com/p/444061702
     类似后缀树的压缩
     Patricia Trie 优化了空间性能，但依然没有完全解决稀疏数据的空间浪费问题——减少了数组的创建，但数组依然可能会有大量的空链接
   - Double Array Trie
     它有两个数组 base 和 check。
     https://zhuanlan.zhihu.com/p/35193582
   - MARISA Trie
     据说是压缩率最高，效率很好的的 Trie 变体。
   - Louds trie
     https://zhuanlan.zhihu.com/p/38194127

华为笔试：
给定一些 x 轴正半轴上的点(N<=10) ,每个点的起点数组 offset，每个点周期 period(<=256，offset < period)，
这些点从 offset[i]开始每隔 period[i]就会出现一次，但最大不会超过 INT_MAX
求使得每个点至少出现一次的最小窗口长度和左端点，如果有多个结果相同取左端点最小的

---

https://maspypy.github.io/library/test/5_atcoder/abc224h.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc227g.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc234ex.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc241e.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc301e.test.cpp
https://atcoder.jp/contests/abc312/tasks/abc312_f
https://maspypy.github.io/library/test/5_atcoder/abc314g.test.cpp
https://atcoder.jp/contests/abc339/submissions/49947896
https://maspypy.github.io/library/test/5_atcoder/abc339f.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc349g.test.cpp
https://maspypy.github.io/library/test/5_atcoder/abc350f.test.cpp
https://maspypy.github.io/library/test/5_atcoder/arc153b.test.cpp
https://maspypy.github.io/library/test/3_yukicoder/1323.test.cpp
https://maspypy.github.io/library/mod/prefix_sum_of_binom.hpp
https://maspypy.github.io/library/test/3_yukicoder/2242.test.cpp
https://www.acwing.com/blog/content/3494/
https://maspypy.github.io/library/graph/stable_matching.hpp
https://maspypy.github.io/library/test/3_yukicoder/2292.test.cpp
https://codeforces.com/contest/852/submission/221983258

https://github.dev/QuBenhao/LeetCode#Rust
