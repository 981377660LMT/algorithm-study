对偶图编号：
https://www.cnblogs.com/llmmkk/p/15566132.html

https://blog.csdn.net/zxw0819/article/details/71436532?depth_1-utm_source=distribute.pc_relevant.none-task&utm_source=distribute.pc_relevant.none-task

https://bakapiano.github.io/2020/07/14/%E5%B9%B3%E9%9D%A2%E5%9B%BE%E8%BD%AC%E5%AF%B9%E5%81%B6%E5%9B%BE%E6%B1%82%E7%BD%91%E7%BB%9C%E6%B5%81/

https://leetcode.cn/problems/check-if-the-rectangle-corner-is-reachable/solutions/2860236/bfs-or-bing-cha-ji-by-tsreaper-ga61/

**平面图最小割等于对偶图最短路**

平面图：就是所有边都不相交的图(`网格图`)
对偶图：平面图的对偶图是一个新的图，其顶点对应原图的面，边对应原图的边，如果两个面有公共边，则对偶图中对应的两个顶点之间有一条边。
**最小割就是对偶图上的最短路（用最小的代价把一个图分成两半）**

作用：解决点数和边数较多的最小割问题(普通网络流的复杂度太大)。

- [ICPC-Beijing 2006] 狼抓兔子：锣鼓 P4001，最经典的例题。
- Cactus Wall：代码力量 1749E。

---

1. 矩形左下角和右上角不能联通：
   把矩形【左边界/上边界】视作节点 n。
   把矩形【下边界/右边界】视作节点 n+1。
   如果节点 n 和 n+1 不在并查集的同一个连通块中，则矩形左下角和右上角可以互相到达。
   [text](<100347. 判断矩形的两个角落是否可达.py>)

2. 最小割转最短路：
   https://www.luogu.com.cn/problem/P2046
3. 最大流->最小割->平面图最小割等于对偶图最短路
   [text](<abc413-G - Big Banned Grid-左上角到右下角是否可达.py>)
