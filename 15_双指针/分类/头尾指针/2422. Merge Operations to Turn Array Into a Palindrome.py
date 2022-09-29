# 合并相邻项以形成回文的最少操作次数
from collections import deque
from typing import List


class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        queue, res = deque(nums), 0
        while len(queue) > 1:
            if queue[0] == queue[-1]:
                queue.popleft()
                queue.pop()
            elif queue[0] < queue[-1]:
                cur = queue.popleft()
                queue[0] += cur
                res += 1
            else:
                cur = queue.pop()
                queue[-1] += cur
                res += 1
        return res


print(Solution().minimumOperations(nums=[4, 3, 2, 1, 2, 3, 1]))
