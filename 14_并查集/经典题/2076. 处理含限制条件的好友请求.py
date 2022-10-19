from collections import defaultdict
from typing import DefaultDict, List
from 可撤销并查集 import RevocableUnionFindArray


# restrictions[i] = [xi, yi] 意味着用户 xi 和用户 yi 不能 成为 朋友
# 如果第 j 个好友请求 成功 ，那么 result[j] 就是 true
# 2 <= n <= 1000
# 0 <= restrictions.length <= 1000
# 1 <= requests.length <= 1000


class UnionFind:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
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
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


class Solution:
    def friendRequests(
        self, n: int, restrictions: List[List[int]], requests: List[List[int]]
    ) -> List[bool]:
        """并查集，连接之前遍历限制找到对应组的边，如果边重合，那么就不能连这条边 O(n^2)"""
        res = []
        uf = UnionFind(n)
        for user1, user2 in requests:
            root1, root2 = uf.find(user1), uf.find(user2)
            if root1 == root2:
                res.append(True)
            else:
                for user3, user4 in restrictions:
                    root3, root4 = uf.find(user3), uf.find(user4)
                    if (root1 == root3 and root2 == root4) or (root1 == root4 and root2 == root3):
                        res.append(False)
                        break
                else:
                    uf.union(user1, user2)
                    res.append(True)

        return res

    def friendRequests2(
        self, n: int, restrictions: List[List[int]], requests: List[List[int]]
    ) -> List[bool]:
        """可撤销并查集"""
        res = []
        uf = RevocableUnionFindArray(n)
        for user1, user2 in requests:
            uf.union(user1, user2)
            ok = True
            for user3, user4 in restrictions:
                if uf.isConnected(user3, user4):
                    uf.revocate()
                    ok = False
                    break
            res.append(ok)
        return res


print(Solution().friendRequests(n=3, restrictions=[[0, 1]], requests=[[0, 2], [2, 1]]))
print(Solution().friendRequests2(n=3, restrictions=[[0, 1]], requests=[[0, 2], [2, 1]]))
# 输出：[true,false]
# 解释：
# 请求 0 ：用户 0 和 用户 2 可以成为朋友，所以他们成为直接朋友。
# 请求 1 ：用户 2 和 用户 1 不能成为朋友，因为这会使 用户 0 和 用户 1 成为间接朋友 (1--2--0) 。
