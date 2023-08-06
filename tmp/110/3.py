from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串),更加通用的并查集写法,支持动态添加"""

    __slots__ = ("part", "_parent", "_rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = dict()
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self._parent:
            self.add(key)
            return key

        while self._parent.get(key, key) != key:
            self._parent[key] = self._parent[self._parent[key]]
            key = self._parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def add(self, key: T) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while self._parent[x] != x:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
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

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


class UnionFindMap2(Generic[T]):
    """不自动合并 需要手动add添加元素"""

    __slots__ = ("part", "_parent", "_rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self.add(item)

    def add(self, key: T) -> "UnionFindMap2[T]":
        if key in self._parent:
            return self
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return self

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2 or root1 not in self._parent or root2 not in self._parent:
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        """此处不自动add"""
        if key not in self._parent:
            return key

        if key != self._parent[key]:
            root = self.find(self._parent[key])
            self._parent[key] = root
        return self._parent[key]

    def isConnected(self, key1: T, key2: T) -> bool:
        if key1 not in self._parent or key2 not in self._parent:
            return False
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent


class UnionFindGraph:
    """并查集维护无向图每个连通块的边数和顶点数"""

    __slots__ = ("n", "part", "_parent", "vertex", "edge")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self.vertex = [1] * n  # 每个联通块的顶点数
        self.edge = [0] * n  # 每个联通块的边数

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.edge[rootX] += 1  # 两个顶点已经在同一个连通块了,这个连通块的边数+1
            return False
        if self.vertex[rootX] > self.vertex[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self.vertex[rootY] += self.vertex[rootX]
        self.edge[rootY] += self.edge[rootX] + 1
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
        return list(set(self.find(i) for i in range(self.n)))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


# 给你一个下标从 0 开始长度为 n 的数组 nums 。

# 每一秒，你可以对数组执行以下操作：

# 对于范围在 [0, n - 1] 内的每一个下标 i ，将 nums[i] 替换成 nums[i] ，nums[(i - 1 + n) % n] 或者 nums[(i + 1) % n] 三者之一。
# 注意，所有元素会被同时替换。


# 请你返回将数组 nums 中所有元素变成相等元素所需要的 最少 秒数。


class Solution:
    def minimumSeconds(self, nums: List[int]) -> int:
        counter = Counter(nums)
        maxCount = max(counter.values())
        maxKey = [nums[i] for i in range(len(nums)) if counter[nums[i]] == maxCount][0]
        if maxCount == 1:
            return 0
        # 模拟,变成相等元素
        res = -1
        visited = [False] * len(nums)
        queue = deque()
        for i in range(len(nums)):
            if nums[i] == maxKey:
                visited[i] = True
                queue.append(i)
        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                if not visited[(cur + 1) % len(nums)]:
                    visited[(cur + 1) % len(nums)] = True
                    queue.append(cur - 1)
                if not visited[(cur - 1) % len(nums)]:
                    visited[(cur - 1) % len(nums)] = True
                    queue.append(cur - 1)
            res += 1
        return res


# nums = [2,1,3,3,2]
# nums = [1,2,1,2]
print(Solution().minimumSeconds(nums=[1, 2, 1, 2]))
print(Solution().minimumSeconds(nums=[2, 1, 3, 3, 2]))
print(Solution().minimumSeconds(nums=[5, 5, 5, 5]))
print(Solution().minimumSeconds(nums=[2, 2, 3, 4, 2]))
# [11,4,10]
print(Solution().minimumSeconds(nums=[11, 4, 10]))  # 1
