https://ddosvoid.github.io/2021/08/04/%E7%8C%AB%E6%A0%91/
https://immortalco.blog.uoj.ac/blog/2102
https://www.cnblogs.com/tzcwk/p/msfz.html
https://www.luogu.com.cn/blog/yizhixiaoyun/cat-tree

猫树问题可以适用于离线解决以下类型的数据结构问题：

与序列有关，且询问是一段区间
序列静态，即，不涉及修改操作

---

线段树的 op 无法快速进行，但是求前后缀很快，怎么办？

一种常见 trick，就是在分治时维护区间的前缀和后缀信息，用左区间的后缀和右区间的前缀求出跨过 mid 的答案。
和同样时空复杂度的 ST 表相比，猫树不需要信息满足可重复贡献，只需要满足可合并，因此适用范围更广。

https://www.luogu.com.cn/problem/CF750E
https://www.luogu.com.cn/problem/P6240
https://www.luogu.com.cn/problem/P4755
https://www.luogu.com.cn/problem/SP1043
https://www.luogu.com.cn/problem/SP2916

---

https://oi-wiki.org/ds/seg/#%E6%8B%93%E5%B1%95---%E7%8C%AB%E6%A0%91
线段树建树的时候需要做 O(n) 次合并操作，而每一次区间询问需要做 O(\log{n}) 次合并操作，询问区间和这种东西的时候还可以忍受，但是当我们需要询问区间线性基这种**合并复杂度高达 O(\log^2{w}) 的信息的话，此时就算是做 O(\log{n}) 次合并有些时候在时间上也是不可接受的。**

而所谓「猫树」就是一种不支持修改，仅仅支持快速区间询问的一种静态线段树。
构造一棵这样的静态线段树**需要 O(n\log{n}) 次合并操作，但是此时的查询复杂度被加速至 O(1) 次合并操作。**
**在处理线性基这样特殊的信息的时候甚至可以将复杂度降至 O(n\log^2{w})。**

---

猫树的核心思想为，将区间分为 O(logn)层，每层处理每个点到区间中心点的答案
