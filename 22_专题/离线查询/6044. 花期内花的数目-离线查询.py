# 所有花按照花期开始时间排序，所有人按照访问时间排序。用小根堆维护当前开着的花的花期的结束时间。
# 时间复杂度 O((N+M)logN+MlogM)。
# 空间复杂度 O(N+M)。


from heapq import heappop, heappush
from typing import List

# 这道题的三种解法:

# 动态开点线段树/树状数组
# 离线查询+小根堆
# 离散化+差分


class Solution:
    def fullBloomFlowers(self, flowers: List[List[int]], persons: List[int]) -> List[int]:
        queries = sorted([(person, index) for index, person in enumerate(persons)])
        flowers = sorted(flowers)

        fi, res = 0, [0] * len(queries)
        pq = []
        for qi in range(len(queries)):
            while fi < len(flowers) and flowers[fi][0] <= queries[qi][0]:
                heappush(pq, flowers[fi][1])
                fi += 1
            while pq and pq[0] < queries[qi][0]:
                heappop(pq)
            res[queries[qi][1]] = len(pq)
        return res
