from typing import List
from collections import defaultdict, deque
from heapq import heappush, heappop

# 2 <= n <= 104
# 每个节点都有一个交通信号灯，每 change 分钟改变一次，从绿色变成红色，再由红色变成绿色，
# 循环往复。所有信号灯都 同时 改变。你可以在 任何时候 进入某个节点，
# 但是 只能 在节点 信号灯是绿色时 才能离开。如果信号灯是  绿色 ，
# 你 不能 在节点等待，必须离开。


# 给你 n、edges、time 和 change ，返回从节点 1 到节点 n 需要的 第二短时间
# 穿过任意一条边的时间是 time 分钟。
# 交通信号灯，每 change 分钟改变一次，
# 在 启程时 ，所有信号灯刚刚变成 绿色 。
# 第二小的值 是 严格大于 最小值的所有值中最小的值。


# bfs求次短路
# 实际上 所有边权都一样可以不用pq 而是普通的deque


class Solution:
    def secondMinimum(self, n: int, edges: List[List[int]], time: int, change: int) -> int:
        adjMap = defaultdict(set)
        for u, v in edges:
            u, v = u - 1, v - 1
            adjMap[u].add(v)
            adjMap[v].add(u)

        # 每个点保存多个距离而不是只有一个
        dist = [[] for _ in range(n)]
        dist[0] = [0]
        pq = [(0, 0)]
        minTime = None

        while pq:
            cost, cur = heappop(pq)

            if cur == n - 1:
                if minTime is None:
                    minTime = cost
                elif minTime < cost:
                    return cost

            if (cost // change) & 1:
                # 到下一个绿灯开始
                cost = (cost // change + 1) * change
            cost += time

            for next in adjMap[cur]:
                # 入队条件
                if not dist[next] or len(dist[next]) == 1 and dist[next][0] < cost:
                    dist[next].append(cost)
                    heappush(pq, (cost, next))

        raise Exception("No path")


print(
    Solution().secondMinimum(n=5, edges=[[1, 2], [1, 3], [1, 4], [3, 4], [4, 5]], time=3, change=5)
)
