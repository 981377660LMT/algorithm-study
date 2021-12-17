# 1 <= s.length <= 16
# 枚举分割点

# n个拆分点 每个拆分点都可以选择拆分或者不拆分 2^(n-1) =>n*2^(n)
class Solution:
    def maxUniqueSplit(self, s: str) -> int:
        def backtrack(index: int, splitCount: int) -> None:
            if index >= length:
                nonlocal maxSplit
                maxSplit = max(maxSplit, splitCount)
                return

            for i in range(index, length):
                substr = s[index : i + 1]
                if substr not in visited:
                    visited.add(substr)
                    backtrack(i + 1, splitCount + 1)
                    visited.remove(substr)

        length = len(s)
        visited = set()
        maxSplit = 1
        backtrack(0, 0)
        return maxSplit

