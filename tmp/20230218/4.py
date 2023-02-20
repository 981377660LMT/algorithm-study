from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 对任一由 n 个小写英文字母组成的字符串 word ，我们可以定义一个 n x n 的矩阵，并满足：

# lcp[i][j] 等于子字符串 word[i,...,n-1] 和 word[j,...,n-1] 之间的最长公共前缀的长度。
# 给你一个 n x n 的矩阵 lcp 。返回与 lcp 对应的、按字典序最小的字符串 word 。如果不存在这样的字符串，则返回空字符串。

# 对于长度相同的两个字符串 a 和 b ，如果在 a 和 b 不同的第一个位置，字符串 a 的字母在字母表中出现的顺序先于 b 中的对应字母，则认为字符串 a 按字典序比字符串 b 小。例如，"aabd" 在字典上小于 "aaca" ，因为二者不同的第一位置是第三个字母，而 'b' 先于 'c' 出现。

from typing import List


def getLCP(s: str) -> List[List[int]]:
    n = len(s)
    lcp = [[0] * (n + 1) for _ in range(n + 1)]
    for i in range(n - 1, -1, -1):
        for j in range(n - 1, -1, -1):
            if s[i] == s[j]:
                lcp[i][j] = lcp[i + 1][j + 1] + 1
    return [row[:-1] for row in lcp[:-1]]


def max(x, y):
    if x > y:
        return x
    return y


def min(x, y):
    if x < y:
        return x
    return y


from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self._parent))

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


# 并查集
class Solution:
    def findTheString(self, lcp: List[List[int]]) -> str:
        n = len(lcp)
        # 无法组成
        uf = UnionFindArray(2 * n + 2)
        visited = [False] * n
        for r in range(n):
            for c in range(n):
                if r == c and lcp[r][c] != n - r:
                    return ""
                if lcp[r][c] > min(n - r, n - c):
                    return ""
                # # s[r:]和s[c:]的最长公共前缀长度为lcp[r][c],合并并记录
                # # 不等,相等
                # for i in range(lcp[r][c]):
                #     uf.union(r + i, c + i)
                # 不等
                if lcp[r][c] == 0:
                    if uf.isConnected(r, c):
                        return ""
                    uf.union(r, c + n)
                    uf.union(c, r + n)
                # 结束不等
                next1, next2 = r + lcp[r][c], c + lcp[r][c]
                if next1 < n and next2 < n:
                    if uf.isConnected(next1, next2):
                        return ""
                    uf.union(next1, next2 + n)
                    uf.union(next2, next1 + n)

        groups = uf.getGroups()

        uf = UnionFindArray(n)
        # 必须在一块
        for g in groups.values():
            if len(g) >= 2:
                for i in range(1, len(g)):

                    uf.union(g[i - 1], g[i])


print(Solution().findTheString(lcp=[[4, 0, 2, 0], [0, 3, 0, 1], [2, 0, 2, 0], [0, 1, 0, 1]]))
print(Solution().findTheString(lcp=[[4, 3, 2, 1], [3, 3, 2, 1], [2, 2, 2, 1], [1, 1, 1, 1]]))
print(Solution().findTheString(lcp=[[4, 3, 2, 1], [3, 3, 2, 1], [2, 2, 2, 1], [1, 1, 1, 3]]))
