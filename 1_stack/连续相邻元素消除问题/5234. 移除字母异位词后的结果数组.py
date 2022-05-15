from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def removeAnagrams(self, words: List[str]) -> List[str]:
        """相邻元素消除用栈，只考虑社么情况可以入栈"""
        stack = []
        for w in words:
            if not stack or sorted(w) != sorted(stack[-1]):
                stack.append(w)
        return stack

