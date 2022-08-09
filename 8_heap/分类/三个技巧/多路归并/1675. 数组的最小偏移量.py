from heapq import heapify, heappop, heappush
from typing import List

# 多路归并


class Solution:
    def minimumDeviation(self, nums: List[int]) -> int:
        n = len(nums)
        arr = [[num] for num in nums]

        for i, v in enumerate(arr):
            if v[0] & 1:
                arr[i].append(v[0] * 2)
            else:
                cur = v[0]
                while cur & 1 == 0:
                    cur //= 2
                    arr[i].append(cur)
                arr[i].reverse()

        pq = [(arr[i][0], i, 0) for i in range(n)]
        heapify(pq)

        res = int(1e10)
        max_ = max(item[0] for item in pq)

        while pq:
            cur, row, col = heappop(pq)
            cand = max_ - cur
            res = cand if cand < res else res
            if col + 1 < len(arr[row]):
                cand = arr[row][col + 1]
                max_ = cand if cand > max_ else max_
                nextItem = (cand, row, col + 1)
                heappush(pq, nextItem)
            else:
                break

        return res
