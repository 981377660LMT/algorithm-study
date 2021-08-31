# 有些数的素因子只有 3，5，7，请设计一个算法找出第 k 个数
# 一个小顶堆，然后每次从小顶堆取一个
import heapq


class Solution:
    def getKthMagicNumber(self, k: int) -> int:
        pq = [1]
        visited = set()
        last = 1

        while k:
            cur = heapq.heappop(pq)
            if cur in visited:
                continue
            heapq.heappush(pq, cur * 3)
            heapq.heappush(pq, cur * 5)
            heapq.heappush(pq, cur * 7)
            k -= 1
            visited.add(cur)
            last = cur
        return last


print(Solution().getKthMagicNumber(5))
