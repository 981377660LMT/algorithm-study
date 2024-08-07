拓扑排序的两个入参
`adjList/adjMap` 和 `deg`

- 注意拓扑排序里加 indegree 的时候
  如果用哈希表+`数组` 可以重复加
  如果用哈希表+`集合` 不能重复加(需要去重)
- 注意拓扑排序里所有顶点入度必须初始化

ps:

1. 将后序遍历的结果进行反转，就是拓扑排序的结果。
   后序遍历的这一特点很重要，之所以拓扑排序的基础是后序遍历，
   **是因为一个任务必须在等到所有的依赖任务都完成之后才能开始开始执行。**

---

https://en.wikipedia.org/wiki/C3_linearization
C3 superclass linearization
(MRO: Method Resolution Order)
python 的多继承解决方法冲突的算法
类似于查找拓扑排序
