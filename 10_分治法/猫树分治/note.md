https://ddosvoid.github.io/2021/08/04/%E7%8C%AB%E6%A0%91/

线段树的 op 无法快速进行，怎么办？

一种常见 trick，就是在分治时维护区间的前缀和后缀信息，用左区间的后缀和右区间的前缀求出跨过 mid 的答案。
和同样时空复杂度的 ST 表相比，猫树不需要信息满足可重复贡献，只需要满足可合并，因此适用范围更广。

https://www.luogu.com.cn/problem/P6240
https://www.luogu.com.cn/problem/P4755
https://www.luogu.com.cn/problem/SP1043
https://www.luogu.com.cn/problem/SP2916
