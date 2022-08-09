# https://simezi-tan.hatenadiary.org/entry/20140920/1411169324

# ROW<=6
# COL<=100
# 求棋盘染色的方案数 使得
# 1.左上角和右下角是黑色
# 2.左上角和右下角可以通过黑色格子连通

# dfs(COL,state)  # 前COL列,此时行的连接状态为state(是一个元组) 用并查集维护
# !列间状态转移dp

import sys
from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)

T = TypeVar("T", int, str)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

    __slots__ = "parent"

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.parent = dict()
        for item in iterable or []:
            self._add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 连接到小的root上"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if root1 < root2:
            root1, root2 = root2, root1
        self.parent[root1] = root2
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
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __contains__(self, key: T) -> bool:
        return key in self.parent


#############################################

ROW, COL = map(int, input().split())
OFFSET = ROW + 1
initParent = list(range(OFFSET))
initParent[1] = 0  # 初始第一个元素的parent为0 即左上角涂黑
dp = defaultdict(int, {tuple(initParent): 1})  # !key为每个行的并查集父结点组成的元组
# print(dp, "init")
for c in range(COL):
    ndp = defaultdict(int)

    for preParent in dp:
        for state in range(1 << ROW):
            paint = set(i + 1 for i in range(ROW) if state & (1 << i))
            uf = UnionFindMap()  # 当前的并查集 注意开两倍表示前一行和后一行
            uf.parent = {r: p for r, p in enumerate(preParent)}
            # print(uf.parent)

            for r in range(1, OFFSET):
                # 这一行连接并查集
                if r in paint and r - 1 in paint:
                    uf.union(r + OFFSET, r - 1 + OFFSET)
                # 这一行与上一行连接并查集
                if r in paint:
                    uf.union(r + OFFSET, r)

            # print(uf.parent, state)
            curParent = [0]
            for r in range(1, OFFSET):
                p = uf.find(r + OFFSET) % OFFSET
                curParent.append(p)

            curParent = tuple(curParent)
            ndp[curParent] += dp[preParent]
            ndp[curParent] %= MOD

    dp = ndp
    print(dp, f"end {c}")

# print(dp)
res = 0
for state in dp:
    if state[-1] == 0:  # 最后一行与0相连
        res += dp[state]
        res %= MOD

print(res)


# TODO 哪里有问题
