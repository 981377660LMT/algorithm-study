1. 求图中的最小环:

- 枚举每条边删除,求 to 到 from 的最短路.
  `O(E*O(最短路复杂度))`
  https://maspypy.github.io/library/graph/mincostcycle.hpp
- 求最短路径树
  `O(V*O(最短路复杂度))`
  https://yukicoder.me/problems/no/1320/editorial

1. 求图中的最大环:
   https://leetcode.cn/problems/shortest-cycle-in-a-graph/solution/yi-tu-miao-dong-mei-ju-qi-dian-pao-bfspy-ntck/
   极端情况下，这会算出一个哈密顿回路，而它是 NP 完全问题，只能通过爆搜得到。
