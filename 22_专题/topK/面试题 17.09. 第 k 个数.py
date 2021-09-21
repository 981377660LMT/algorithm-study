# 有些数的素因子只有 3，5，7，请设计一个算法找出第 k 个数
# 一个小顶堆，然后每次从小顶堆取一个
import heapq
from typing import List


class Solution:
    def getKthMagicNumber(self, k: int) -> int:
        pq = [1]
        visited = set()
        last = 1

        while len(visited) < k:
            cur = heapq.heappop(pq)
            if cur in visited:
                continue
            visited.add(cur)
            heapq.heappush(pq, cur * 3)
            heapq.heappush(pq, cur * 5)
            heapq.heappush(pq, cur * 7)
            last = cur
        return last


print(Solution().getKthMagicNumber(5))


class Solution2:
    def majorityElement(self, nums: List[int]) -> int:
        val = 0
        count = 0
        for num in nums:
            if num == val:
                count += 1
            elif count == 0:
                val = num
                count = 1
            else:
                count -= 1
        return val if nums.count(val) > len(nums) / 2 else -1


print(3 / 2)

print(Solution2().majorityElement([3, 2, 3]))
