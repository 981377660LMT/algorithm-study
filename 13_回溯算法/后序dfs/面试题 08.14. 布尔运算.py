from functools import lru_cache
from typing import List


class Solution:
    def countEval(self, s: str, result: int) -> int:
        @lru_cache(None)
        def dfs(s: str) -> List[int]:
            count = [0] * 2
            if s in '01':
                count[int(s)] = 1
                return count
            for i, c in enumerate(s):
                if c not in '01':
                    left = dfs(s[:i])
                    right = dfs(s[i + 1 :])
                    for leftValue, leftCount in enumerate(left):
                        for rightValue, rightCount in enumerate(right):
                            if c == '|':
                                cur = leftValue | rightValue
                            elif c == '&':
                                cur = leftValue & rightValue
                            elif c == '^':
                                cur = leftValue ^ rightValue
                            count[cur] += leftCount * rightCount
            return count

        return dfs(s)[result]

