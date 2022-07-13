# 有一个长度为n的序列A ,你可以执行以下操作任意次:选择一个i(i＋k ≤n)i(i+k≤n)，并swap(a[i], a[i＋k)
# 问:是否可以使得最后得到的序列A 满足不下降。

# !并查集分组
# 思路:
# 考虑对相距为k的点之间用并查集进行合并操作，最后我们会得到若干个连通块，我们可以对连通块内部进行任意排序，
# 对于本题来说是想得到一个递增的数组，故对每个连通块进行从小到大排序，然后在判断整合后的数组是否是递增的即可。
import sys
import os
from collections import defaultdict
from typing import DefaultDict, List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
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

    def __str__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


def main() -> None:
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    uf = UnionFindArray(n)
    for i in range(n):
        if i + k < n:
            uf.union(i, i + k)

    groups = uf.getGroups()
    for group in groups.values():
        group.sort(reverse=True, key=nums.__getitem__)  # 编号排序

    res = []
    for i in range(n):
        root = uf.find(i)
        res.append(nums[groups[root].pop()])

    target = sorted(nums)
    if res == target:
        print("Yes")
    else:
        print("No")


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
