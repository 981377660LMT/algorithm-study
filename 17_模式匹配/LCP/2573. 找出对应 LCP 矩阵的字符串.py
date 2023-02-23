# 并查集+LCP的性质
# 从LCP矩阵还原出最小字典序的字符串
# 如果不存在或者存在矛盾,返回空字符串
# n<=1000

# https://leetcode.cn/problems/find-the-string-with-lcp/solution/tan-xin-xuan-neng-xuan-de-zui-xiao-zi-fu-3tt2/
# !LCP 值从根本上来说就是字符间的相等关系和不等关系，与字符串之间的大小关系无关
# LCP 信息是高度冗余的

# 重要结论:
#  !如果 LCP[i][j] > 0 那么 第i个字符等于第j个字符
#  !如果 LCP[i][j] = 0 那么 第i个字符不等于第j个字符
# 因此可以使用并查集寻找相同字符的组,然后顺序遍历分配字符,最后求出LCP验证是否正确


from typing import DefaultDict, List
from collections import defaultdict
from getLCP import getLCP


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


class Solution:
    def findTheString(self, lcp: List[List[int]]) -> str:
        n = len(lcp)
        uf = UnionFindArray(n)
        for i in range(n):
            for j in range(i + 1, n):
                if lcp[i][j] > 0:
                    uf.union(i, j)

        groups = uf.getGroups()
        assign = [-1] * n  # ord - 97
        id = 0
        for i in range(n):
            if assign[i] != -1:
                continue
            root = uf.find(i)
            for g in groups[root]:
                assign[g] = id
            id += 1
        if id > 26:
            return ""

        res = "".join(chr(assign[i] + 97) for i in range(n))
        resLcp = getLCP(res)
        for r in range(n):
            for c in range(n):
                if resLcp[r][c] != lcp[r][c]:
                    return ""
        return res


# lcp = [[4,0,2,0],[0,3,0,1],[2,0,2,0],[0,1,0,1]]
print(Solution().findTheString([[4, 0, 2, 0], [0, 3, 0, 1], [2, 0, 2, 0], [0, 1, 0, 1]]))
