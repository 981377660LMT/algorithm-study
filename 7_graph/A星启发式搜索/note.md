https://www.bilibili.com/video/BV1bv411y79P?spm_id_from=333.999.0.0

1. Roughly speaking, A\* is a BFS with customized priority for selecting a node to expand; the priority is to select the lowest f()=heuristic()+cost()
   严格地说，A\*是一个具有自定义优先级的 BFS(`增强版 dijkstra`)，用于选择要扩展的节点；优先顺序是选择最低的 `f()=heuristic()+cost()`
2. A heuristic function is admissible if: it's always a lower bound of the actual cost from current state to target state (in this example, manhatten distance is admissible)
   启发式函数是容许误差的。它总是目标状态的实际成本的下限(在这个例子中，曼哈顿是很好的)。一个重要的定理：如果启发式函数是可接受的，那么 A\*就可以找到最优路径。
3. An important theorem: If heuristic function is admissible, then A\* is guaranteed to find the optimal path
   一个重要的定理：如果启发式函数是容许误差的，则 A\*被保证找到最优路径

---

迭代加深搜索：解决无上限的搜索问题
https://www.luogu.com.cn/blog/LawrenceSivan/uva12558-ai-ji-fen-shuo-egyptian-fractions-hard-version-die-dai-jia-sh

埃及分数 Egyptian Fractions (HARD version)
https://www.luogu.com.cn/problem/UVA12558
