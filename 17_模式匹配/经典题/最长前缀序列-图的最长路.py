# 此处采用图的最长路

# O(n⋅m)
class Solution:
    def solve(self, words):
        def dfs(cur: str) -> int:
            """有向图的最长路"""
            return 1 + max((dfs(next) for next in adjMap[cur]), default=0)

        adjMap = {word: set() for word in words}
        for word in words:
            pre = word[:-1]
            if pre in adjMap:
                adjMap[pre].add(word)

        return max((dfs(word) for word in words), default=0)


print(Solution().solve(words=["abc", "ab", "x", "xy", "abcd"]))

# We can form the following sequence: ["ab", "abc", "abcd"].


# 1.排序+Trie树
# 2. 图的最长路
