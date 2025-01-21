1. 动态维护dfs序

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
