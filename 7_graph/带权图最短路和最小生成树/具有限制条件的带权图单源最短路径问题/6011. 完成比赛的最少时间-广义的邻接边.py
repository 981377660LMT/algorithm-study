from heapq import heappop, heappush
from typing import List


# 1 <= numLaps <= 1000
# 1 <= tires.length <= 105

# 总结：
# 这道题一开始想贪心的解法(贪心ptsd)，sortedList弄了好久，
# 最后才意识到是dp 状态由圈数唯一决定 但是怎么求每个圈的最小时间花费呢?
# 关键是要意识到最多是乘18次

# 这题也可以最短路dijk
# 预处理不换轮子，连走i轮需要的最短时间，minCost数组的值代表每个点和他邻接点相连边的边权
# 优先队列起点为连走i轮的所有情况，要求起点到 numLaps 的最短路径
# dist数组控制入队条件，然后直接使用dijk算法求到达 numLaps 的最短路即可

# 主要是建图不会建
INF = int(1e20)

# 时间复杂度 O((V+E)logV)，由于图最多20*V条边，所以这个复杂度是可以接受的


class Solution:
    def minimumFinishTime(self, tires: List[List[int]], changeTime: int, numLaps: int) -> int:
        # 建边:不换轮子 连走i轮需要的最短时间
        minCost = [INF] * 20
        for a0, q in tires:
            minCost[0], curCost = min(a0, minCost[0]), a0
            for i in range(1, 20):
                curCost += a0 * q
                if curCost > 1e6:
                    break
                minCost[i] = min(curCost, minCost[i])
                a0 *= q

        # dijk求最短路，
        # 注意pq可以不用heapify， 因为heappop方法返回数组的第一个元素(堆顶)，而这里的pq第一个元素必定是cost最小的，要不要heapify对第一轮搜索没有影响
        pq = [(minCost[i], i + 1) for i in range(min(len(minCost), numLaps))]
        dist = [INF] * (numLaps + 1)
        while pq:
            cost, step = heappop(pq)
            if step == numLaps:
                return cost

            for nextStep, weight in enumerate(minCost, start=1):
                nextCost = cost + weight + changeTime
                if step + nextStep <= numLaps and nextCost < dist[step + nextStep]:
                    dist[step + nextStep] = nextCost
                    heappush(pq, (nextCost, step + nextStep))

        return -1


# 21 25
print(Solution().minimumFinishTime(tires=[[2, 3], [3, 4]], changeTime=5, numLaps=4))
print(Solution().minimumFinishTime(tires=[[1, 10], [2, 2], [3, 4]], changeTime=6, numLaps=5))

