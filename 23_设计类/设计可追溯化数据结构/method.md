https://blog.csdn.net/Ljnoit/article/details/119886465
https://github.com/csvoss/retroactive

1. 并查集
   部分可追溯化：
   不考虑删除操作(有删除操作时的部分可追溯化实际上就是动态图连通性问题)
   `因为 union 具有交换律，直接用并查集`
   完全可追溯化：LCT 维护删除时间最大生成树
2. 队列
   部分可追溯化：链表
   完全可追溯化：平衡树
3. 栈
   可追溯化：平衡树
4. 双端队列
   可追溯化：平衡树
5. 优先队列
   部分可追溯化：线段树
   [text](PartiallyRetroactivePriorityQueue.go)
   完全可追溯化：堆的完全可追溯化比较困难，因此使用的方法都是比较通用的，适用范围比较广的方法。

---

通用的技巧:

- NonRetroactive -> PartiallyRetroactive：

- PartiallyRetroactive -> FullyRetroactive：
  O(sqrt(m)) 操作分块
  对于每个块，将这个块之前的所有操作组成的操作序列用 partial
  的方法进行维护

---

r 是一个描述允许发生多远的过去追溯操作的参数
m 是对数据结构执行的追溯更新的总数。
