# 2045. 到达目的地的第二短时间
# https://leetcode.cn/problems/second-minimum-time-to-reach-destination/description/
# 2 <= n <= 1e4
# 每个节点都有一个交通信号灯，每 change 分钟改变一次，从绿色变成红色，再由红色变成绿色，
# 循环往复。所有信号灯都 同时 改变。你可以在 任何时候 进入某个节点，
# 但是 只能 在节点 信号灯是绿色时 才能离开。如果信号灯是  绿色 ，
# 你 不能 在节点等待，必须离开。
# 给你 n、edges、time 和 change ，返回从节点 1 到节点 n 需要的 第二短时间
# 穿过任意一条边的时间是 time 分钟。
# 交通信号灯，每 change 分钟改变一次，
# 在 启程时 ，所有信号灯刚刚变成 绿色 。
# !第二小的值 是 严格大于 最小值的所有值中最小的值。
#
# dijkstra求次短路


from typing import List
from heapq import heappush, heappop


INF = int(1e20)


class Solution:
    def secondMinimum(self, n: int, edges: List[List[int]], time: int, change: int) -> int:
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            u, v = u - 1, v - 1
            adjList[u].append(v)
            adjList[v].append(u)

        dist1, dist2 = [INF] * n, [INF] * n
        dist1[0] = 0
        pq = [(0, 0)]

        while pq:
            curDist, cur = heappop(pq)
            if curDist > dist2[cur]:
                continue
            if (curDist // change) & 1:
                curDist = (curDist // change + 1) * change  # 到下一个绿灯开始
            curDist += time
            for next in adjList[cur]:
                cand = curDist
                if curDist < dist1[next]:
                    dist1[next], cand = curDist, dist1[next]
                    heappush(pq, (dist1[next], next))
                elif dist1[next] < curDist < dist2[next]:
                    dist2[next] = cand
                    heappush(pq, (dist2[next], next))

        return dist2[n - 1]


if __name__ == "__main__":
    print(
        Solution().secondMinimum(
            n=5, edges=[[1, 2], [1, 3], [1, 4], [3, 4], [4, 5]], time=3, change=5
        )
    )
