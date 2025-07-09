# 3607. 电网维护-倒序查询
# 给你一个整数 c，表示 c 个电站，每个电站有一个唯一标识符 id，从 1 到 c 编号。
#
# 这些电站通过 n 条 双向 电缆互相连接，表示为一个二维数组 connections，其中每个元素 connections[i] = [ui, vi] 表示电站 ui 和电站 vi 之间的连接。直接或间接连接的电站组成了一个 电网 。
#
# 最初，所有 电站均处于在线（正常运行）状态。
#
# 另给你一个二维数组 queries，其中每个查询属于以下 两种类型之一 ：
#
# [1, x]：请求对电站 x 进行维护检查。如果电站 x 在线，则它自行解决检查。如果电站 x 已离线，则检查由与 x 同一 电网 中 编号最小 的在线电站解决。如果该电网中 不存在 任何 在线 电站，则返回 -1。
#
# [2, x]：电站 x 离线（即变为非运行状态）。
#
# 返回一个整数数组，表示按照查询中出现的顺序，所有类型为 [1, x] 的查询结果。
# 注意：电网的结构是固定的；离线（非运行）的节点仍然属于其所在的电网，且离线操作不会改变电网的连接性。
#
# !删除点、查询联通分量中的最小编号

from typing import Callable, List, Optional


INF = int(1e18)


class UnionFindArray:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(
        self, key1: int, key2: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        if beforeUnion is not None:
            beforeUnion(root1, root2)
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if self._data[key] < 0:
            return key
        self._data[key] = self.find(self._data[key])
        return self._data[key]

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


class Solution:
    def processQueries(
        self, n: int, connections: List[List[int]], queries: List[List[int]]
    ) -> List[int]:
        uf = UnionFindArray(n)
        for u, v in connections:
            u, v = u - 1, v - 1
            uf.union(u, v)

        offlineTime = [INF] * n
        for i in range(len(queries) - 1, -1, -1):
            t, x = queries[i]
            if t == 2:
                x -= 1
                offlineTime[x] = i

        groupMin = [INF] * n
        for i in range(n):
            if offlineTime[i] == INF:
                root = uf.find(i)
                groupMin[root] = min(groupMin[root], i)

        res = []
        for i in range(len(queries) - 1, -1, -1):
            t, x = queries[i]
            x -= 1
            root = uf.find(x)
            if t == 1:
                if i < offlineTime[x]:
                    res.append(x + 1)
                else:
                    res.append(groupMin[root] + 1 if groupMin[root] != INF else -1)
            else:
                if offlineTime[x] == i:
                    groupMin[root] = min(groupMin[root], x)

        res.reverse()
        return res
