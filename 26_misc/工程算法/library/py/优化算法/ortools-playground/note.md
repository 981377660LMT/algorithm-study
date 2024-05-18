关键词：运筹学、线性规划、整数规划、网络流

OR-Tools 主要包含以下 4 个方面的求解器：

- 约束优化(Constraint Programming)：用于计算可行解，与 z3 求解器的功能类似。​​ 
  ​https://developers.google.cn/optimization/cp​​
  约束编程，即 constraint programming，简称 CP。CP 是找到一个可行的解决方案，而不是找到一个最优的解决方案，它关注约束和变量，一般没有目标函数，即使有目标也仅仅是增加约束条件，将一大堆可能的解决方案缩小到一个更易于管理的子集。
- 线性规划(Linear and Mixed-Integer Programming)：与 PuLP 库的功能类似。​​ 
  ​https://developers.google.cn/optimization/lp​​
- 车辆路线图(Vehicle Routing)：计算给定约束条件下最佳车辆路线的专用库。​
  ​​https://developers.google.cn/optimization/routing​​
- 网络流(Graph Algorithms)：用于求解最大流量和最小成本流等问题。​​ 
  ​https://developers.google.cn/optimization/flow​​

---

[OR-Tools 简介](https://zhuanlan.zhihu.com/p/547925740)
[OR-Tools 入门教程](https://zhuanlan.zhihu.com/p/551807837)
[Google OR-Tools 使用技巧](https://zhuanlan.zhihu.com/p/629643347)
[Google OR-Tools 搜索参数、源码编译、cython 实现及一些使用 Tips](https://zhuanlan.zhihu.com/p/374530559)
[谷歌 OR-Tools CP-SAT 求解器原理解析](https://zhuanlan.zhihu.com/p/631406803)
[Google OR-Tools 中 AddElement 约束的使用](https://zhuanlan.zhihu.com/p/632897768)

---

https://www.zhihu.com/column/c_157167393
[ortools 系列：运筹优化工具 google ortools 简介](https://zhuanlan.zhihu.com/p/55089642)

---

[OR-Tools 官网](https://developers.google.cn/optimization?hl=zh-cn)
[OR-Tools 官档中文用法大全（CP、LP、VRP、Flows 等）](https://blog.51cto.com/u_11866025/5833945)

TODO:
https://atcoder.jp/contests/abc354/submissions/53613731
https://atcoder.jp/contests/abc354/submissions/53619029
https://blog.51cto.com/u_11866025/5833945
