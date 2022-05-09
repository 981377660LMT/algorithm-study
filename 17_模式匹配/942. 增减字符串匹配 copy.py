from collections import deque
from typing import List


class Solution:
    def diStringMatch(self, s: str) -> List[int]:
        """由范围 [0,n] 内所有整数组成的 n + 1"""
        queue = deque(range(len(s) + 1))
        return [queue.popleft() if c == "I" else queue.pop() for c in s] + [queue[0]]


# 如果不要求[0,n]也一样
