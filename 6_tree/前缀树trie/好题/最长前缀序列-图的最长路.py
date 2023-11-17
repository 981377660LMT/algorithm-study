# 此处采用图的最长路

# O(n⋅m)
from collections import defaultdict
from functools import lru_cache
from typing import List


class Solution:
    def solve(self, words: List[str]) -> int:
        @lru_cache(None)
        def dfs(cur: int) -> int:
            """有向图的最长路"""
            return 1 + max((dfs(next) for next in adjMap[cur]), default=0)

        id_ = defaultdict(lambda: len(id_))
        adjMap = defaultdict(set)
        for word in words:
            pre = word[:-1]
            adjMap[id_[pre]].add(id_[word])
        return max((dfs(id_[word]) for word in words), default=0)


print(Solution().solve(words=["abc", "ab", "x", "xy", "abcd"]))

# We can form the following sequence: ["ab", "abc", "abcd"].


# 1. 排序+Trie树
# 2. 图的最长路
