https://xyzl.blog.luogu.org/Shortest-Path-Tree-SPT

- 最短路径树 SPT，是一张连通图的生成树(每个点都联通)
  **从树的根节点 u 到任意点 v 的路径都为原图 G 中 u 到 v 的最短路径**
- 最短路径树 × 最小生成树的区别

  - 最小生成树只是满足全图联通且边权和最小，而最短路径树是满足`从根节点到任意点的最短路和原图最短路相同`
  - 最短路径树的边权和 ≥ 最小生成树的边权和

- SPT 的技巧

1. 删边:不在最短路径树上的边被删除对最短路不会有影响
2. 复原:求出最短路径树中每个结点的前驱边,然后从根节点开始进行树的遍历即可
3. 最小权值 SPT: 在最短路的条件下,保证每个条边的权值最小

```python
for next, weight, eid in adjList[cur]:
   cand = dist[cur] + weight
   if cand < dist[next]:
       dist[next] = cand
       preE[next] = eid
       preV[next] = cur
       heappush(pq, (dist[next], next))
   elif cand == dist[next]:  # 在最短路相等的情况下，扩展到同一个节点，后出堆的点连的边权值一定更小
       preE[next] = eid
       preV[next] = cur
```
