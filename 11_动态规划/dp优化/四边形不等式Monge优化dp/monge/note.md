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

---

## Totally Monotone 和 Monotone 有什么区别？

"Monotone"（单调）和"Totally Monotone"（全单调）是两种不同的数学性质，它们都可以用来描述函数或矩阵。

1. "Monotone"（单调）：通常用来描述一个函数或序列，如果它的值随着输入的增加而单调增加或单调减少，那么我们就说这个函数或序列是单调的。例如，函数 f(x) = x^2 在 x >= 0 时就是单调的，因为当 x 增加时，f(x)的值也增加。

2. "Totally Monotone"（全单调）：是一种更强的单调性质，通常用来描述矩阵或二元函数。对于一个矩阵或函数，如果它满足对于任意的 i < j 和 k < l，如果矩阵的[i][k]元素大于矩阵的[j][l]元素，那么矩阵的[i][l]元素也大于矩阵的[j][k]元素，那么我们就说这个矩阵或函数是全单调的。

总的来说，全单调是一种更强的单调性质，它要求函数或矩阵在多个维度上都保持单调性。而单调只要求函数或序列在一个维度上保持单调性。
