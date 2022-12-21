from collections import deque
from typing import List
from sortedcontainers import SortedList


class Solution:
    def minHeightShelves(self, books: List[List[int]], shelfWidth: int) -> int:
        n = len(books)
        heights = [0] + [h for _, h in books]
        preSum = [0] * (n + 1)
        dp = [0] * (n + 1)
        queue = deque()
        sl = SortedList()
        for i in range(1, n + 1):
            preSum[i] = preSum[i - 1] + books[i - 1][0]
            left = i

            while queue and heights[queue[-1][1]] <= heights[i]:
                sl.remove(dp[queue[-1][0] - 1] + heights[queue[-1][1]])
                left = queue.pop()[0]

            queue.append([left, i])
            sl.add(dp[queue[-1][0] - 1] + heights[i])
            while queue and preSum[i] - preSum[queue[0][0] - 1] > shelfWidth:
                sl.remove(dp[queue[0][0] - 1] + heights[queue[0][1]])
                queue[0][0] += 1
                if queue[0][0] > queue[0][1]:
                    queue.popleft()
                else:
                    sl.add(dp[queue[0][0] - 1] + heights[queue[0][1]])

            dp[i] = sl[0]

        return dp[n]


print(Solution().minHeightShelves(books=[[1, 1] for _ in range(100000)], shelfWidth=4))
