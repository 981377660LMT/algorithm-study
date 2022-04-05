from collections import deque
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def triangularSum(self, nums: List[int]) -> int:
        queue = deque(nums)
        while len(queue) > 1:
            nextQueue = deque()
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                if queue:
                    next = queue[0]
                    nextQueue.append((cur + next) % 10)
            queue = nextQueue
        return queue[0]


print(Solution().triangularSum([1, 2, 3, 4, 5]))
