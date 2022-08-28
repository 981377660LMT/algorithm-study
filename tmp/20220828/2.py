from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def removeStars(self, s: str) -> str:
        stack = []
        for char in s:
            if char != "*":
                stack.append(char)
            else:
                if stack and stack[-1] != "*":
                    stack.pop()
                else:
                    stack.append(char)
        return "".join(stack)
