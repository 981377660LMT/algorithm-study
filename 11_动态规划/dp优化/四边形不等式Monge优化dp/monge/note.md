https://kmyk.github.io/monotone-matrix-visualizer/
https://topcoder-g-hatena-ne-jp.jag-icpc.org/spaghetti_source/20120915/1347668163.html
https://maspypy.github.io/library/convex/monge.hpp

[monge グラフ上の d-辺最短路の d=1,...,N における列挙](https://nyaannyaan.github.io/library/dp/monge-d-edge-shortest-path-enumerate.hpp)
[monge グラフ上の d-辺最短路](https://nyaannyaan.github.io/library/dp/monge-d-edge-shortest-path.hpp)
[monge グラフ上の最短路](https://nyaannyaan.github.io/library/dp/monge-shortest-path.hpp)
[monotone minima](https://nyaannyaan.github.io/library/dp/monotone-minima.hpp)

---

Monge 图：

```
给定一个n个顶点的DAG，图中的每个顶点i都有一条边指向编号j(i<j)，且边的权重满足Monge性质(也就是四边形不等式)，dist(a, b) + dist(c, d) ≤ dist(a, d) + dist(c, b)，这个图就是Monge图
```

---

- **SMAWK（totally monotone minima; TM minima)**:
  `O(ROW+COL)` 时间复杂度求出一个全单调矩阵的每行的最小值
- **Monotone Minima**:
  `O(ROW+COL*logROW)` 时间复杂度求出一个单调矩阵的每行的最小值
- **Monge**:
  `O(nlogmax)`求出 Monge 图的 d 边最短路
