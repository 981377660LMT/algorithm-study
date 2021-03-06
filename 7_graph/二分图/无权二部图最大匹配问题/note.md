**无向图**匹配

1. **二分图(BipartiteGraph)**匹配问题两种方法
   最大匹配/完全匹配
   **最大流算法**(MaxFlow):新建虚拟的源点和汇点，无向图转换为有向图,所有边容量为 1,流入汇点的每一个流量即为一个匹配
   **匈牙利算法**(Hangurian)：从 左侧非匹配点出发 从右向左走永远走匹配边 终止于另外一个非匹配点:增广路径
   则最大匹配数加一
   匈牙利算法 BFS/DFS **关键** O(VE)

2. 最大匹配：一个二分图，最多的一一匹配

3. 完全匹配（所有的顶点都是匹配点）一定是最大匹配，最大匹配不一定是完全匹配

4. 可以使用最大流算法解决匹配问题，即所有边的容量都为 1，最大流即为最大匹配数

5. LCP 4 覆盖问题：

- 首先将棋盘所有的点认为是类似国际象棋的棋盘的样子（黑白相间），然后就可以建图。
- 可以认为每个黑白连起来的 2x1 的格子都是一条边，连接了二分图中的两个点（黑点和白点）
- 这样问题就可以转化为最大匹配问题

6. **匈牙利算法**解决最大流(hungarian algorithm)

- 从左边开始，往右边去连第一个还没有匹配的点
- 如果右边的点是一个匹配点，从右往左的边，永远走匹配边
- 最后将匹配的边换成未匹配的边，未匹配的边再变成匹配的边
- 匹配边和非匹配边交替出现：交替路
- **终止与另外一个非匹配点**：即找到了一条增广路径
- 有增广路径，意味着最大匹配数可以加一
- 在交替过程中，由于我们起始于非匹配点，终止与非匹配点，所以中间经过的非匹配边的数目一定比匹配边的数目大 1

7. 匈牙利算法总结：对左侧每一个尚未匹配的点，不断地寻找可以增广的交替路

8. 可以利用 BFS/DFS 寻找增广路径

9. BFS 队列中只存储左边的点

10. 经典问题：Lintcode 1576. 最佳匹配

匈牙利算法的流程，其中找妹子是个递归的过程，最最关键的字就是“ 腾”字
其原则大概是：有机会上，没机会创造机会也要上

**最小点覆盖问题**
另外一个关于二分图的问题是求最小点覆盖：我们想找到最少的一些点，使二分图所有的边都至少有一个端点在这些点之中。倒过来说就是，`删除包含这些点的边，可以删掉所有边。`
**一个二分图中的最大匹配数等于这个图中的最小点覆盖数。**
图的覆盖是一些顶点（或边）的集合，使得图中的每一条边（每一个顶点）都至少接触集合中的一个顶点（边）

**转换:矩阵就是二分图；矩阵第 i 行第 j 列(i,j)处如果为 1，表示二分图中左边右边这两个点相连**
（洛谷 P1129） [ZJOI2007]矩阵游戏
把矩阵转化为二分图（`左侧集合代表各行，右侧集合代表各列，某位置为 1 则该行和该列之间有边`）
行交换操作：选择矩阵的任意两行，交换这两行:相当于给行重命名
我们可以在保持当前二分图结构不变的情况下，把右侧点的编号进行改变，这与交换的效果是一样的。

（vijos1204） CoVH 之柯南开锁
锁是由 M\*N 个格子组成, 其中某些格子凸起(灰色的格子).` 每一次操作可以把某一行或某一列的格子给按下去.`
如果柯南能在组织限定的次数内将所有格子都按下去, 那么他就能够进入总部. 但是 OIBH 组织不是吃素的, 他们的限定次数恰是最少次数.
请您帮助柯南计算出开给定的锁所需的最少次数.
`按下一行或一列，其实就是删掉与某个点相连的所有边`。现在要求最少的操作次数，想想看，这不就是求最小点覆盖数吗？所以直接套匈牙利算法即可。代码略。
