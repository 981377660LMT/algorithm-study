# 在线每次查询和修改是O(logn^2)的

from typing import List, Union


class PersistentUnionFind:
    """可持久化并查集:用可持久化数组维护并查集的parent数组和rank数组

    不能路径压缩，因为路径压缩时间复杂度是均摊的 这里最坏情况是O(nlogn)
    同样,基于均摊复杂度的珂朵莉树、Splay、替罪羊树都不能简单的可持久化

    需要按秩合并
    """

    __slots__ = ("curVersion", "_parent", "_rank")

    def __init__(self, n: int, updateTimes: int) -> None:
        self.curVersion = 0
        self._parent = PersistentArray.create(list(range(n + 10)), updateTimes)
        self._rank = PersistentArray.create([1] * (n + 10), updateTimes)

    def union(self, version: int, u: int, v: int) -> int:
        """如果u和v已经连通,则版本号不变,否则返回新版本号"""
        root1, root2 = self.find(version, u), self.find(version, v)
        if root1 == root2:
            return self.curVersion
        rank1, rank2 = self._rank.query(version, root1), self._rank.query(version, root2)
        if rank1 > rank2:
            root1, root2 = root2, root1
            rank1, rank2 = rank2, rank1
        self._parent.update(version, root1, root2)
        self._rank.update(version, root2, rank1 + rank2)
        self.curVersion += 1
        return self.curVersion

    def find(self, version: int, u: int) -> int:
        while True:
            pre = self._parent.query(version, u)
            if pre == u:
                return u
            u = pre

    def isConnected(self, version: int, u: int, v: int) -> bool:
        root1, root2 = self.find(version, u), self.find(version, v)
        return root1 == root2


class PersistentArray:
    __slots__ = (
        "curVersion",
        "_n",
        "_leftChild",
        "_rightChild",
        "_treeValue",
        "_roots",
        "_nodeId",
    )

    @staticmethod
    def create(sizeOrArray: Union[int, List[int]], updateTimes: int) -> "PersistentArray":
        """创建一个可持久化数组,并指定更新次数的上限"""
        isArray = isinstance(sizeOrArray, list)
        n = len(sizeOrArray) if isArray else sizeOrArray
        assert n > 0, f"length must be positive, but {n} received"
        if isArray:
            return PersistentArray(sizeOrArray, updateTimes)
        return PersistentArray([0] * n, updateTimes)

    def __init__(self, nums: List[int], updateTimes: int):
        n = len(nums)
        size = 4 * n + n.bit_length() * updateTimes

        self.curVersion = 0  # !从0开始编号,版本0表示初始状态数组
        self._n = n
        self._leftChild = [0] * size
        self._rightChild = [0] * size
        self._treeValue = [0] * size
        self._roots = [0] * (updateTimes + 1)

        self._nodeId = 0
        self._roots[0] = self._build(0, n - 1, nums)

    def query(self, version: int, index: int) -> int:
        """访问历史版本`version`的数组的`index`位置的值

        Args:
            version (int): 版本号 >= 0
            index (int): 位置 >= 0

        Returns:
            int: 数组的值
        """
        assert (
            0 <= version <= self.curVersion
        ), f"version must be in [0, {self.curVersion}], but {version} received"
        assert 0 <= index < self._n, f"index must be in [0, {self._n}), but {index} received"
        return self._query(self._roots[version], 0, self._n - 1, index)

    def update(self, version: int, index: int, value: int) -> int:
        """在历史版本`version`的数组上更新`index`位置的值为`value`，并返回新数组的版本号

        Args:
            version (int): 版本号 >= 0
            index (int): 位置 >= 0
            value (int): 更新后的值

        Returns:
            int: 新数组的版本号
        """
        assert (
            0 <= version <= self.curVersion
        ), f"version must be in [0, {self.curVersion}], but {version} received"
        assert 0 <= index < self._n, f"index must be in [0, {self._n}), but {index} received"
        rootId = self._update(self._roots[version], 0, self._n - 1, index, value)
        self.curVersion += 1
        self._roots[self.curVersion] = rootId
        return self.curVersion

    def _build(self, left: int, right: int, array: List[int]) -> int:
        node = self._nodeId
        self._nodeId += 1
        if left == right:
            self._treeValue[node] = array[left]
            return node

        mid = (left + right) // 2
        self._leftChild[node] = self._build(left, mid, array)
        self._rightChild[node] = self._build(mid + 1, right, array)
        return node

    def _query(self, curRoot: int, left: int, right: int, pos: int) -> int:
        if left == right:
            return self._treeValue[curRoot]

        mid = (left + right) // 2
        if pos <= mid:
            return self._query(self._leftChild[curRoot], left, mid, pos)
        return self._query(self._rightChild[curRoot], mid + 1, right, pos)

    def _update(self, preRoot: int, left: int, right: int, pos: int, value: int) -> int:
        # !copy preVersion
        node = self._nodeId
        self._nodeId += 1
        self._leftChild[node] = self._leftChild[preRoot]
        self._rightChild[node] = self._rightChild[preRoot]
        self._treeValue[node] = self._treeValue[preRoot]
        if left == right:
            self._treeValue[node] = value
            return node

        mid = (left + right) // 2
        if pos <= mid:
            self._leftChild[node] = self._update(self._leftChild[preRoot], left, mid, pos, value)
        else:
            self._rightChild[node] = self._update(
                self._rightChild[preRoot], mid + 1, right, pos, value
            )
        return node


Q = int(1e5)
if __name__ == "__main__":
    uf = PersistentUnionFind(10, updateTimes=Q)
    v0 = 0
    v1 = uf.union(v0, 1, 2)
    v2 = uf.union(v1, 2, 3)
    assert uf.isConnected(v2, 1, 3)
    assert uf.isConnected(v2, 1, 2)
    assert not uf.isConnected(v2, 1, 4)
    v3 = uf.union(v0, 1, 3)
    assert uf.isConnected(v3, 1, 3)
