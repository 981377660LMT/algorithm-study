1. 怪盗基德的滑翔翼
   初始时，怪盗基德可以在任何一幢建筑的顶端。他可以选择一个方向逃跑，但是不能中途改变方向
   他只能往下滑行,请问，他最多可以经过多少幢不同建筑的顶部
   `左右LIS最大值`
2. 登山
   就是不连续浏览海拔相同的两个景点，并且一旦开始下山，就不再向上走了
   队员们希望在满足上面条件的同时，尽可能多的浏览景点，你能帮他们找出最多可能浏览的景点数么？
   `求最长的山脉子序列`
   `枚举分割点，左右LIS相加的最大值`
3. 合唱队形
   `求最长的山脉子序列`

4. 友好城市
   每对友好城市都向政府申请在河上开辟一条直线航道连接两个城市，但是由于河上雾太大，政府决定避免任意两条航道交叉，以避免事故。
   编程帮助政府做出一些批准和拒绝申请的决定，使得在保证任意两条航线不相交的情况下，被批准的申请尽量多
   `不相交：排序后没有逆序，一个维度排序，求另一个维度的最长上升子序列`

5. 最大上升子序列和
   `就是LIS模板`

6. 拦截导弹
   但是这种导弹拦截系统有一个缺陷：虽然它的第一发炮弹能够到达`任意的高度`，但是以后每一发炮弹都不能高于前一发的高度。
   输入导弹依次飞来的高度（雷达给出的高度数据是不大于 30000 的正整数，导弹数不超过 1000），计算这套系统最多能拦截多少导弹，如果要拦截所有导弹最少要配备多少套这种导弹拦截系统。
   `不断删除下降子序列，问最少删除几次(答案是 最长上升子序列的长度)`
   dilworth 定理：`偏序集上最小链划分中链的数量等于其反链长度的最大值。`
7. 导弹防御系统
   一套防御系统的导弹拦截高度要么`一直 严格单调 上升要么一直 严格单调 下降`。给定即将袭来的一系列导弹的高度，请你求出至少需要多少套防御系统，就可以将它们全部击落。
   `dfs爆搜 求最小深度 要记录全局最小值`
   `为什么不用bfs?:空间会爆炸`
8. 最长公共上升子序列
