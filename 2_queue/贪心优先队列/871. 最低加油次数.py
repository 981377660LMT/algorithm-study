from typing import List
import heapq

# 为了到达目的地，汽车所必要的最低加油次数是多少？如果无法到达目的地，则返回 -1 。
# 事后诸葛


class Solution:
    def minRefuelStops(self, target: int, startFuel: int, stations: List[List[int]]) -> int:
        res = 0
        pq = []
        curFuel = startFuel
        curPos = 0
        stations.append([target, 0])

        for _, (pos, add) in enumerate(stations):
            curFuel -= pos - curPos
            curPos = pos
            while curFuel < 0 and pq:
                curFuel += -heapq.heappop(pq)
                res += 1

            if curFuel < 0:
                return -1

            heapq.heappush(pq, -add)

        return res


print(Solution().minRefuelStops(100, 10, [[10, 60], [20, 30], [30, 30], [60, 40]]))

# 输出：2
# 解释：
# 我们出发时有 10 升燃料。
# 我们开车来到距起点 10 英里处的加油站，消耗 10 升燃料。将汽油从 0 升加到 60 升。
# 然后，我们从 10 英里处的加油站开到 60 英里处的加油站（消耗 50 升燃料），
# 并将汽油从 10 升加到 50 升。然后我们开车抵达目的地。
# 我们沿途在1两个加油站停靠，所以返回 2 。

