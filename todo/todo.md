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
   ∏
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

https://github.dev/QuBenhao/LeetCode#Rust

---

2^61-1 优化字符串哈希模版
https://qiita.com/keymoon/items/11fac5627672a6d6a9f6
一种比较科学的字符串哈希实现方法 - cup-pyy 的文章 - 知乎
https://zhuanlan.zhihu.com/p/784510655

https://maspypy.github.io/library/string/rollinghash.hpp 主要是 mod2^61-1 下的乘法
https://maspypy.github.io/library/string/rollinghash_2d.hpp
https://maspypy.github.io/library/alg/monoid/rollinghash.hpp

`18_哈希/字符串哈希/dynamic/hashString.go`
