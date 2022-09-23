# 有n张卡片(骨牌),每张卡片正面有一个数字p[i],背面也有一个数字q[i];
# 保证所有牌中正面和反面出现的数字都是1到n的排列,现在想要取一些牌
# !这些牌正反面必须包含1～n所有数字,求方案数.
# n<=2e5

# 置换环
# !思路:我们可以把每张牌看作 p[i]连向q[i]的边，
# 由于p,q都是一个排列，意味着图中每个点的入度为1，出度也为1，也就形成了一个环。
# 问题就转化为了给你若干个不联通的环，其中每个环你可以删除若干条边，
# 但要满足每一个点都和一条边相连，即每两条相邻的边需要至少选一条。问满足的方案数

# !1. 并查集找环
# !2. 环上dp (可以删若干遍，但是每个点必须和一条边相连的方案数)
# !3. 每个环乘法原理计数

import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def main() -> None:

    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    # dp[i][0/1] 表示 第i个数取或不取的方案数
    # 第一个数不选 最后一个数就必须选
    dp1 = [[0, 0] for _ in range(n)]
    dp1[0][0] = 1
    for i in range(1, n):
        dp1[i][0] = dp1[i - 1][1]
        dp1[i][1] = (dp1[i - 1][0] + dp1[i - 1][1]) % MOD

    # 第一个数选 最后一个数可选可不选
    dp2 = [[0, 0] for _ in range(n)]
    dp2[0][1] = 1
    for i in range(1, n):
        dp2[i][0] = dp2[i - 1][1]
        dp2[i][1] = (dp2[i - 1][0] + dp2[i - 1][1]) % MOD

    uf = UnionFindMap()
    for u, v in zip(nums1, nums2):
        uf.union(u, v)
    groups = uf.getGroups()
    res = 1

    for group in groups.values():
        size = len(group)
        if size > 1:
            res *= dp1[size - 1][1] + dp2[size - 1][0] + dp2[size - 1][1]
            res %= MOD
    print(res)


if __name__ == "__main__":
    from collections import defaultdict
    from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar

    T = TypeVar("T", bound=Hashable)

    class UnionFindMap(Generic[T]):
        """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

        def __init__(self, iterable: Optional[Iterable[T]] = None):
            self.part = 0
            self.parent = dict()
            self.rank = defaultdict(lambda: 1)
            for item in iterable or []:
                self._add(item)

        def union(self, key1: T, key2: T) -> bool:
            """rank一样时 默认key2作为key1的父节点"""
            root1 = self.find(key1)
            root2 = self.find(key2)
            if root1 == root2:
                return False
            if self.rank[root1] > self.rank[root2]:
                root1, root2 = root2, root1
            self.parent[root1] = root2
            self.rank[root2] += self.rank[root1]
            self.part -= 1
            return True

        def find(self, key: T) -> T:
            if key not in self.parent:
                self._add(key)
                return key

            while self.parent.get(key, key) != key:
                self.parent[key] = self.parent[self.parent[key]]
                key = self.parent[key]
            return key

        def isConnected(self, key1: T, key2: T) -> bool:
            return self.find(key1) == self.find(key2)

        def getRoots(self) -> List[T]:
            return list(set(self.find(key) for key in self.parent))

        def getGroups(self) -> DefaultDict[T, List[T]]:
            groups = defaultdict(list)
            for key in self.parent:
                root = self.find(key)
                groups[root].append(key)
            return groups

        def _add(self, key: T) -> bool:
            if key in self.parent:
                return False
            self.parent[key] = key
            self.rank[key] = 1
            self.part += 1
            return True

        def __repr__(self) -> str:
            return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

        def __len__(self) -> int:
            return self.part

        def __contains__(self, key: T) -> bool:
            return key in self.parent

    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
