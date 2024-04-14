https://zhuanlan.zhihu.com/p/583236322
https://zhuanlan.zhihu.com/p/482150152
完全背包与同余最短路

- 同余最短路问题是指在给定的带权有向图中，找到一条最短的路径，
  使得路径上的所有边的权值在模 m 意义下同余于一个给定的值 k。

- ModShortestPath

  “同余最短路”一般是指一种和不定方程相关的建模方式
  对于一堆系数 a0, a1, ..., an-1 (0<=a0<=a1<=...<=an-1,最小的非零 a0 称为 base)
  !我们可以用这个叫同余最短路的方法，来确定他们的`线性组合∑ai*xi`可能取到的值(xi 非负)
  更确切地说,可以处理出一个数组 dist[0,1,...,base-1]
  !这里 dist[i]记录的是最小的 x,满足 x=i(mod base)且 x 能被上述系数线性表出
  具体的方式是把余 0,余 1,...,余 base-1 看成一个个结点,每个点表示一个剩余类
  对每个剩余类,考虑所有的转移:
  !也就是 x 加上每一个 ai,即连一条从 x 到 (x+ai) mod base 的边，边权为 ai
  这样建完图从 0 号点开始跑最短路,每转移一条长度为 d 的边就对应着线性组合加了一个 d
  得到最后的 dist 数组

  在每个剩余类中， dist[i],dist[i]+base,...,dist[i]+k\*base(k>=0) 都能到达
  如果 dist[i]==INF,则表示这个剩余类不可达

---

非常类似 完全背包(无穷背包)

一般的完全背包是 `O(容量*物品个数)`，同余最短路的背包是 `O(物品价值*物品个数log(物品个数))`

---

TODO：
https://github.com/EndlessCheng/codeforces-go/blob/86d1fd150c7b53861a52fa81ce456666bf547691/copypasta/graph.go#L2008
https://zhuanlan.zhihu.com/p/583236322
https://www.cnblogs.com/alex-wei/p/17531487.html (同余最短路的转圈技巧)
https://oi-wiki.org/graph/mod-shortest-path/
https://zhuanlan.zhihu.com/p/672216458
