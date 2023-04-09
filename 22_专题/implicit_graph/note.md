在线拓扑排序/在线并查集/在线最小生成树
https://maspypy.github.io/library/graph/online_toposort.hpp
https://maspypy.github.io/library/graph/online_mst.hpp
https://maspypy.github.io/library/graph/online_unionfind.hpp

不预先给出图，
而是通过函数 findUnused 交互来找到下一个候选人。

---

离线就是普通的 bfs 需要把图先建出来。
在线就是不需要建图,bfs 过程可以中止，甚至是异步的。
通过调用 api 来逐步 bfs, 对应力扣上的说法就是可以做成 class 来调用。
函数可以换成生成器，整个 bfs 可以通过交互完成,可以出交互题。

```python
class Bfs:
  def __init__(self, start: int):
    pass
  def setUsed(self, u):
    pass
  def findUnused(self) -> Optional[int]:
    pass
  def getDist(self, u) -> List[int]:
    pass
```
