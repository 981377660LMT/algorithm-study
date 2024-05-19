关键词：运筹学、线性规划、整数规划、网络流

OR-Tools 主要包含以下 4 个方面的求解器：

- 约束优化(Constraint Programming)：用于计算可行解，与 z3 求解器的功能类似。​​ 
  ​https://developers.google.cn/optimization/cp​​
  约束编程，即 constraint programming，简称 CP。CP 是找到一个可行的解决方案，而不是找到一个最优的解决方案，它关注约束和变量，一般没有目标函数，即使有目标也仅仅是增加约束条件，将一大堆可能的解决方案缩小到一个更易于管理的子集。
  **要求：纯整数编程问题**
  (CP-SAT 求解器要求所有的约束和目标必须是整数)
- 线性规划(Linear and Mixed-Integer Programming)：与 PuLP 库的功能类似。​​ 
  ​https://developers.google.cn/optimization/lp​​
  线性规划，即 Linear optimization，用来计算一个问题的最佳解决方案，问题被建模为一组线性关系。 OR-Tools 提供的求解 LP 问题的方法是 MPSolver
- 车辆路线图(Vehicle Routing)：计算给定约束条件下最佳车辆路线的专用库。​
  ​​https://developers.google.cn/optimization/routing​​
- 网络流(Graph Algorithms)：用于求解最大流量和最小成本流等问题。​​ 
  ​https://developers.google.cn/optimization/flow​​

**尽量使用大写 api，atcoder 上只有大写 api**

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

不错的：
[OR-Tools 官网](https://developers.google.cn/optimization?hl=zh-cn)
[OR-Tools 官档中文用法大全（CP、LP、VRP、Flows 等）](https://blog.51cto.com/u_11866025/5833945)

OR-Tools 简介

- 约束优化(CP)
  入门示例：获取一个解、多个解
  密码拼图
  数学运算
  八皇后问题
  解数独
  护士排班问题
  cp_model 无法求解的问题
  - 只能创建整数变量
  - 不支持非线性运算
- 线性规划(LP)
  入门示例
  MIP 求解器求解整数规划问题

  - LP 求解器（例如上面使用的 GLOP）变量结果是可以会得到非整数的，对于严格整数规划问题，除了使用 CP-SAT 求解器外还可以使用 MIP 求解器（例如 SCIP）。

  斯蒂格勒的饮食问题
  医院每天护士人数分配
  数据包络分析(DEA)

  - 是一种用于评估相对效率的非参数方法。它是一种多变量线性规划技术，用于评估一组具有多个输入和输出的单位（如企业、组织或个人）的相对效率

- 背包与装箱问题
  背包问题
  多背包问题
  装箱问题
  双 11 购物的凑单问题
- 分配问题
  基础示例
  小组任务分配
  小组任务分配 2
  工作量约束的任务分配
- 车辆路线图(VRPs)
  旅行推销员问题
  示例: 钻一块电路板
  Dimensions 维度
  车辆路径问题
  常见设置
  基本搜索限制
  first_solution_strategy 设置
  本地搜索选项
  定义初始路径
  自定义每辆车的起始和结束位置
  允许任意的开始和结束位置
  带承重限制的 VRP 问题
  惩罚并放弃访问点
  物流配送相关的 VRP 问题：每一个项目创建一个取货和送货请求，添加每个项目必须由同一辆车辆取货和送货的约束，添加每个项目在送货前必须先取货
  带有时间窗口的 VRP 问题：每个位置只能在特定的时间范围内访问
- 网络流（Network Flows）
- 最大流量（Maximum Flows）
- 最小成本流（Minimum Cost Flows）
