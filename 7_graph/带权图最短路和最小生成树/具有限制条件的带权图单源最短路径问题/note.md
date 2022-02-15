这类题**已经不再是由 普通的一维 dist 控制入堆条件了** 一般是 dijkstra 模板+ **多维的 dist** `根据限制条件控制是否再入堆`
考虑遍历最小费用，如果同一个城市出现更多费用，但是可以获得更多时间的时候，
这类题一般是 dijkstra 模板 + dist 控制再入堆的条件：如果还没看过这个点，或者当前的限制条件比之前更优，则加入堆

```Python
pq: List[Tuple[Cost, ID, Time]] = [(passingFees[0], 0, 0)]
```

两个方法:

1. 多维 dist 数组 dist[id][limit]=cost
   `1928. 规定时间内到达终点的最小花费.py`
   `2093_到达城市的的最小价格.py`
   `LCP 35. 电动车游城市.py `
   **多一个加油限制**
   两个维度的 dist 数组
   使用两个维度来控制是将状态否继续入队
   ```Python
   if cur_cost < dist[cur_id][cur_time]:
      dist[cur_id][cur_time] = cur_cost
      ...
   ```
2. dist 数组不存 cost 而是直接存 limit
   `1928. 规定时间内到达终点的最小花费.py`
   `2093_到达城市的的最小价格.py`

   ```Python
    if cur_time < visited[cur_id]:
      visited[cur_id] = cur_time
      ...

    if cur_discount < dist[cur_id]:
      dist[cur_id] = cur_discount
      ...

   ```

`787. K 站中转内最便宜的航班.py`
**多一个辗转次数限制**
使用 k 来控制入队

`LCP 35. 电动车游城市.py `
**多一个加油限制**
两个维度的 dist 数组
使用两个维度来控制是将状态否继续入队
