https://github.com/tdzl2003/leetcode_live/tree/master/flow20220719
https://www.bilibili.com/video/BV1jt4y1t7pd

- 割的定义(分两组)：

设置两个集合 S 集合和 T 集合，且 s 点属于 S 集合，t 点属于 T 集合
其它所有点属于 S 集合或 T 集合之一
所有来自于 S 集合，指向 T 集合的边，权值之和为割的大小

- 能否套用最小割方法的判断：

1. 把**最少/最小代价**对应到最小割
2. 需要把所有元素**分成两组**
3. 不合法的方案，定义成**流量无穷大**的边，表示不能割这条边
   时间复杂度为 Dinic 复杂度`O(V^2*E)` 数据一般在 100-200 之间

- 建图过程：

根据可以做二元选择的对象定义所有的点
根据不合法方案连接无穷大的边，明确点的集合归属含义
根据点的集合归属含义，明确所有需要支付的成本

- 额外思路：
  如果问的是**最小花费**，直接用最小割
  如果问的是**最大利益**，先假设能获得所有收益，然后把得不到的收益也记作成本，最小割来解决最大化收益问题

1. 棋子问题

有 `N*M` 的棋盘，上面有 K 个棋子，位于 Ai 行 Bi 列。
一次操作可以消除一行的棋子或者消除一列的棋子。
问`最少`要多少次操作，才能消灭所有的棋子？

- 思考:
  `每个棋子由行消除还是列消除，网络中的点是 行 与 列`
  `不合法的方案，表示有棋子没有被消除。`
- 建模:
  `属于 S 集合表示行不消除，属于 T 集合表示列不消除`
- 不合法方案的边：
  `对于为 1 的单元格，行不消除且列不消除非法`
  `对应单元格建立行指向列的无穷大流量边`
- 支付成本：
  `s 到行，流量为 1，割这条边表示行消除`
  `列到 t，流量为 1，割这条边表示列消除`

2. 设备采购问题

有 N 种设备，采购价格分别是 A1,A1,...AN
和 M 个项目，项目收益分别是 B1,B1,...BM
另外有 K 个约束条件，Ci,Di 表示完成 Di 号项目，必须要拥有 Ci 号设备。
问`最高`能获得的利润是多少（利润等于收益-成本）

- 思考:
  `先假设能获得所有收益，然后把得不到的收益也记作成本`
  `设备：采购与不采购；项目：做与不做`
  `不合法方案：项目依赖的设备没有采购，但又要做对应项目`
- 建模:
  `属于 S 集合表示不采购设备，属于 T 集合表示完成项目`
- 不合法方案的边：
  `项目依赖设备，但设备未采购，不合法`
  `有依赖时，从设备指向项目，流量无穷大`
- 支付成本：
  `s 到采购设备，流量为采购成本(割这条边表示买设备，需要支付成本)`
  `完成项目到 t，流量为收益项目(割这条边表示不做项目，收益减少)`

3. 取数问题

在一个 `M*N` 的棋盘中，每个方格有一个正整数。
现在要从方格中取若干个数，使任意 2 个数所在方格没有公共边，并使取出的数总和`最大`。
试设计一个满足要求的取数算法。
(数据很小的时候可以轮廓线 dp)

- 思考:
  `先假设能获得所有数，然后建图使方案合法`
  `每个数选择与不选择`
  `不合法的方案：相邻的格子不能同时选`
- 建模:
  `属于 S 集合表示选取偶数格子，属于 T 集合表示选取奇数格子`
- 不合法方案的边：
  `相邻的格子不能同时选`
  `相邻的格子，从偶数格向奇数格建边，流量无穷大`
- 支付成本：
  `s 到偶数格子，流量为偶数格子的值(割这条边表示不选这个偶数格子，收益减少)`
  `奇数格子到 t，流量为奇数格子的值(割这条边表示不选这个奇数数格子，收益减少)`

4. 另一个棋子问题
   有 `N*M `的棋盘，每个方格有一个数，其中有正数也有负数。要求选定其中的若干行和若干列，
   满足：
   **要求负数的格子不能同时选择所在行和所在列。**
   如果方格所在的行或所在的列被选中，它的分数被加到总分上，不论是正数还是负数。
   如果行和列同时被选中，分数也只累加一次。
   试设计一个算法，问能得到的`最大`得分是多少。

- 思考:
  `先假设能获得所有和为正数的行的和+所有和为正数的列的和`
  `每个数选择与不选择`
  `不合法的方案：行和列同时选了负数格子`
- 建模:
  `属于 S 集合表示选择行，属于 T 集合表示选择列`
- 不合法方案的边：
  `1.行和列同时选了负数格子`
  `负数格子被行和列都选取了，从选择行向选择列建边，流量无穷大`
  `2.非负数格子被统计了两次`
  `非负数格子被行和列都选取了，从选择列向选择行建边，流量为非负数格子的值(割掉这条边表示在答案里少算一次)`
- 支付成本：
  `s 到行的和为正数的行，流量为行的和(割这条边表示不选这个行，收益减少)`
  `列的和为正数的列到 t，流量为奇数格子的值(割这条边表示不选这个列，收益减少)`
