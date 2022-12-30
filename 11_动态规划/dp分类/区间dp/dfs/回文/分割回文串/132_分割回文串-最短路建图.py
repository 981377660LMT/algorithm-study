# 1 <= s.length <= 2000

# 此解法O(n^2)
from collections import defaultdict, deque
from functools import lru_cache

# dfs/bfs都可求最短路


class Solution(object):
    def minCut(self, s: str) -> int:
        """
        :给你一个字符串 s，请你将 s 分割成一些子串，使每个子串都是回文。
        返回符合要求的 最少分割次数 。
        """
        n = len(s)
        adjMap = defaultdict(set)

        # 建图
        for start in range(n):
            left, right = start, start
            while left >= 0 and right < n and s[left] == s[right]:
                adjMap[left].add(right + 1)
                left -= 1
                right += 1

            left, right = start, start + 1
            while left >= 0 and right < n and s[left] == s[right]:
                adjMap[left].add(right + 1)
                left -= 1
                right += 1

        # 求最短路
        queue = deque([0])
        visited = set([0])
        depth = 0
        while queue:
            curLen = len(queue)
            for _ in range(curLen):
                cur = queue.popleft()
                if cur == n:
                    return depth - 1
                for next in adjMap[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            depth += 1
        return n - 1

    def minCut2(self, s: str) -> int:
        """
        :给你一个字符串 s，请你将 s 分割成一些子串，使每个子串都是回文。
        返回符合要求的 最少分割次数 。
        """
        n = len(s)
        adjMap = defaultdict(set)

        # 建图
        for start in range(n):
            left, right = start, start
            while left >= 0 and right < n and s[left] == s[right]:
                adjMap[left].add(right + 1)
                left -= 1
                right += 1

            left, right = start, start + 1
            while left >= 0 and right < n and s[left] == s[right]:
                adjMap[left].add(right + 1)
                left -= 1
                right += 1

        @lru_cache(None)
        def dfs(index: int) -> int:
            if index >= n:
                return 0
            return min(1 + dfs(next) for next in adjMap[index])

        return dfs(0) - 1


# python切片 不超时 js超时
print(Solution().minCut2(s="aab"))
