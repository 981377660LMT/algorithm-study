# 给定一个无重复元素的数组，求交换成有序数组的最少次数(minSwap)
# 可以选择 同一层 上任意两个位置，交换这两个位置的值。
# 返回按 `严格递增顺序` 排序所需的最少操作数目。


from typing import List


def minSwapToSortedArray(nums: List[int]) -> int:
    mp = {num: i for i, num in enumerate(nums)}
    target = sorted(nums)
    uf = UnionFind()
    for i, num in enumerate(target):
        uf.union(i, mp[num])
    groups = uf.getGroups().values()
    return sum(len(g) - 1 for g in groups if len(g) > 1)  # 每个环需要交换的次数为g-1


from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFind(Generic[T]):
    """当元素不是数组index时(例如字符串),更加通用的并查集写法,支持动态添加"""

    __slots__ = ("part", "parent", "rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.rank = dict()
        for item in iterable or []:
            self.add(item)

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
            self.add(key)
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

    def add(self, key: T) -> bool:
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


if __name__ == "__main__":

    # 2471. 逐层排序二叉树所需的最少操作数目
    # Definition for a binary tree node.
    class TreeNode:
        def __init__(self, val=0, left=None, right=None):
            self.val = val
            self.left = left
            self.right = right

    class Solution:
        def minimumOperations(self, root: Optional["TreeNode"]) -> int:
            def dfs(root: Optional["TreeNode"], dep: int) -> None:
                if root is None:
                    return
                levels[dep].append(root.val)
                dfs(root.left, dep + 1)
                dfs(root.right, dep + 1)

            levels = defaultdict(list)
            dfs(root, 0)
            return sum(minSwapToSortedArray(level) for level in levels.values())
