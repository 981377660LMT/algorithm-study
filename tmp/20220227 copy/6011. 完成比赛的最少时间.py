from math import log
from typing import List, Tuple
from collections import defaultdict

from sortedcontainers import SortedList

# 1 <= numLaps <= 1000
# 1 <= tires.length <= 105

# 总结：
# 这道题一开始想贪心的解法(贪心ptsd)，sortedList弄了好久，
# 最后才意识到是dp 状态由圈数唯一决定 但是怎么求每个圈的最小时间花费呢?

# 这题也可以最短路dijk
class Solution:
    def minimumFinishTime(self, tires: List[List[int]], changeTime: int, numLaps: int) -> int:
        """tires[i] = [fi, ri] 表示第 i 种轮胎如果连续使用，第 x 圈需要耗时 fi * ri(x-1) 秒"""
        """每一圈后，你可以选择耗费 changeTime 秒 换成 任意一种轮胎（也可以换成当前种类的新轮胎）。"""
        res = 0
        count = 0
        sortedList = SortedList()
        for _, tire in enumerate(tires):
            # 不换
            okCount = tire[0] * ()
            sortedList.add((tire[0], tire[0], tire[1]))

        # dp啊
        # 不换的最小耗时 以及对应的编号
        def dfs(arg):
            pass

        # pre = None
        # while count < numLaps:
        #     cost, f, r = sortedList[0]
        #     res += cost
        #     count += 1

        #     # if pre is not None:
        #     #     sortedList.add(pre)

        #     # changeTime+f <= cost 就换轮胎
        #     while count < numLaps and sortedList and cost * r <= sortedList[0][1] + changeTime:
        #         cost *= r
        #         res += cost
        #         count += 1

        #     if count >= numLaps:
        #         break

        #     # sortedList.pop(0)
        #     # 换轮胎
        #     res += changeTime
        #     pre = (f, f, r)

        # return res


# 21 25
print(Solution().minimumFinishTime(tires=[[2, 3], [3, 4]], changeTime=5, numLaps=4))
print(Solution().minimumFinishTime(tires=[[1, 10], [2, 2], [3, 4]], changeTime=6, numLaps=5))


# class Solution:
#     def minimumFinishTime(self, tires: List[List[int]], changeTime: int, numLaps: int) -> int:
#         n = len(tires)
#         g = [inf] * 25
#         for a, b in tires:
#             g[0], tot = min(a, g[0]), a
#             for j in range(1, 25):
#                 tot += a * b
#                 if tot > 1e6:
#                     break
#                 g[j] = min(tot, g[j])
#                 a *= b
#         while g[-1] == inf:
#             g.pop()
#         pq = [(g[i], i + 1) for i in range(min(len(g), numLaps))]
#         vis = [False] * (numLaps + 1)
#         while pq:
#             w, t = heappop(pq)
#             if t == numLaps:
#                 return w
#             if vis[t]:
#                 continue
#             vis[t] = True
#             for i in range(len(g)):
#                 if t + i + 1 > numLaps or vis[t + i + 1]:
#                     break
#                 ww = w + g[i] + changeTime
#                 heappush(pq, (ww, t + i + 1))
