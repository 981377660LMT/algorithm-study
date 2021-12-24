from typing import List
from heapq import heapify, heappop, heappush


class Solution:
    def lastStoneWeight(self, stones: List[int]) -> int:
        pq = [-s for s in stones]
        heapify(pq)

        while len(pq) >= 2:
            a, b = -heappop(pq), -heappop(pq)
            if a == b:
                continue
            heappush(pq, -(a - b))

        return -pq[0] if pq else 0


print(Solution().lastStoneWeight([2, 7, 4, 1, 8, 1]))
# 输出：1
# 解释：
# 先选出 7 和 8，得到 1，所以数组转换为 [2,4,1,1,1]，
# 再选出 2 和 4，得到 2，所以数组转换为 [2,1,1,1]，
# 接着是 2 和 1，得到 1，所以数组转换为 [1,1,1]，
# 最后选出 1 和 1，得到 0，最终数组转换为 [1]，这就是最后剩下那块石头的重量。
