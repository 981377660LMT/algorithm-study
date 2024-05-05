from typing import List
import heapq

# 数组 下标从 0 开始 ，其中 piles[i] 表示第 i 堆石子中的石子数量。另给你一个整数 k ，请你执行下述操作 恰好 k 次
# 返回执行 k 次操作后，剩下石子的 最小 总数。
class Solution:
    def minStoneSum(self, piles: List[int], k: int) -> int:
        pq = [-v for v in piles]
        heapq.heapify(pq)
        for _ in range(k):
            heapq.heappush(pq, heapq.heappop(pq) // 2)
        return -sum(pq)


print(Solution().minStoneSum(piles=[5, 4, 9], k=2))
