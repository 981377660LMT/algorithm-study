from collections import deque
from typing import Generator

# 如果可以通过将 A 中的两个小写字母精确地交换位置 K 次得到与 B 相等的字符串，我们称字符串 A 和 B 的相似度为 K（K 为非负整数）。
# 给定两个字母异位词 A 和 B ，返回 A 和 B 的相似度 K 的最小值。
# 1 <= A.length == B.length <= 20
# A 和 B 只包含集合 {'a', 'b', 'c', 'd', 'e', 'f'} 中的小写字母。
# 1 <= s1.length <= 20


class Solution:
    def kSimilarity(self, s1: str, s2: str) -> int:
        def genNexts(pre: str) -> Generator[str, None, None]:
            """ "每次交换至少对齐一个字母，没有对齐字母的操作是没意义的."""
            first = next((i for i in range(len(pre)) if pre[i] != s2[i]), len(pre))
            for j in range(first + 1, len(pre)):
                if pre[j] == s2[first] and pre[first] != s2[first]:
                    yield pre[:first] + pre[j] + pre[first + 1 : j] + pre[first] + pre[j + 1 :]

        queue, visited = deque([(s1, 0)]), set([s1])
        while queue:
            cur, cost = queue.popleft()
            if cur == s2:
                return cost
            for next_ in genNexts(cur):
                if next_ not in visited:
                    visited.add(next_)
                    queue.append((next_, cost + 1))
        return -1


print(Solution().kSimilarity(s1="abc", s2="bca"))
