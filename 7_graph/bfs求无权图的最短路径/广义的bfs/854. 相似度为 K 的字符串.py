from collections import deque
from typing import Generator

# 如果可以通过将 A 中的两个小写字母精确地交换位置 K 次得到与 B 相等的字符串，我们称字符串 A 和 B 的相似度为 K（K 为非负整数）。
# 给定两个字母异位词 A 和 B ，返回 A 和 B 的相似度 K 的最小值。

# 1 <= A.length == B.length <= 20
# A 和 B 只包含集合 {'a', 'b', 'c', 'd', 'e', 'f'} 中的小写字母。
class Solution:
    def kSimilarity(self, s1: str, s2: str) -> int:
        def getNexts(pre: str) -> Generator[str, None, None]:
            """"Each child node requires one swap to change from x and each child node has one character more similiar to B than x."""
            swap = 0
            while pre[swap] == s2[swap]:
                swap += 1
            for j in range(swap + 1, len(pre)):
                # if pre[j] == s2[swap]:
                if pre[j] == s2[swap] and pre[swap] != s2[swap]:
                    yield pre[:swap] + pre[j] + pre[swap + 1 : j] + pre[swap] + pre[j + 1 :]

        queue, visited = deque([(s1, 0)]), set([s1])
        while queue:
            cur, cost = queue.popleft()
            if cur == s2:
                return cost
            for next in getNexts(cur):
                if next not in visited:
                    visited.add(next)
                    queue.append((next, cost + 1))
        return -1


print(Solution().kSimilarity(s1="abc", s2="bca"))
