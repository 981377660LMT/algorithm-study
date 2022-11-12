dynamic-connectivity 动态连通性
https://en.wikipedia.org/wiki/Dynamic_connectivity
https://cp-algorithms.com/data_structures/deleting_in_log_n.html

在计算和图论中，`动态连接结构`是一种数据结构，它动态地维护有关图的连接组件的信息。
图的顶点集合 V 是固定的，但边的集合 E 可以改变。这三种情况，按照难易程度依次是：

- 边只添加到图中（这可以称为增量连接）；
- 边只从图中删除（这可以称为递减连通性）；
- 可以添加或删除边（这可以称为完全动态连接）。

在每次添加/删除一条边之后，动态连接结构应该自行调整，以便它可以快速回答“ x 和 y 之间是否有路径？”形式的查询。（相当于：“顶点 x 和 y 是否属于同一个连通分量？”）。

**golang 实现**

- 可持久化并查集：https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/union_find.go#L349
- 可撤销并查集 回滚并查集
  https://github.dev/EndlessCheng/codeforces-go/blob/cca30623b9ac0f3333348ca61b4894cd00b753cc/copypasta/union_find.go#L441
