火灾扩散/病毒扩散
都是非常麻烦的 bfs
模拟 bfs 扩散

关键是 `visited` 和 `queue`

---

- visited 是 0/1/2/3 的二维数组
  0:未访问,1:人访问过,2:洪水访问过,3:墙壁
- queue1 和 queue2 是人/火的当前层队列
  在 while queue1 中层序扩散
