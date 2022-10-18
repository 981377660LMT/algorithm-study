from typing import List, Union


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
        assert 0 <= version <= self.curVersion
        assert 0 <= index < self._n
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
        assert 0 <= version <= self.curVersion
        assert 0 <= index < self._n
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


if __name__ == "__main__":
    nums = [59, 46, 14, 87, 41]
    v0 = 0
    persistentArray = PersistentArray.create(nums, updateTimes=3)
    assert persistentArray.query(0, 0) == 59
    assert persistentArray.query(0, 1) == 46
    assert persistentArray.query(0, 2) == 14
    assert persistentArray.query(0, 3) == 87
    assert persistentArray.query(0, 4) == 41

    v1 = persistentArray.update(0, 0, 100)
    assert persistentArray.query(v1, 0) == 100

    v2 = persistentArray.update(v1, 1, 200)
    assert persistentArray.query(v2, 1) == 200

    v3 = persistentArray.update(v0, 2, 300)
    for i in range(5):
        print(persistentArray.query(v3, i))
