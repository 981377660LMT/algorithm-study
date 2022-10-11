from typing import Callable, Generic, List, Protocol, TypeVar, Union

S = TypeVar("S")
"""线段树维护的值的类型"""

F = TypeVar("F")
"""懒标记的类型/操作的类型"""


class IOperation(Generic[S, F], Protocol):
    """线段树的操作的接口"""

    def e(self) -> S:
        """线段树的值的幺元"""
        ...

    def id(self) -> "F":
        """懒标记的幺元"""
        ...

    def op(self, leftData: "S", rightData: "S") -> "S":
        """线段树的合并操作"""
        ...

    def mapping(self, parentLazy: "F", childData: "S") -> "S":
        """父结点的懒标记更新子结点的值"""
        ...

    def composition(self, parentLazy: "F", childLazy: "F") -> "F":
        """父结点的懒标记更新子结点的懒标记"""
        ...


class AtcoderLazySegmentTree(Generic["S", "F"]):

    __slots__ = (
        "_n",
        "_size",
        "_log",
        "_data",
        "_lazy",
        "_e",
        "_id",
        "_op",
        "_mapping",
        "_composition",
    )

    def __init__(
        self,
        sizeOrArray: Union[int, List["S"]],
        operation: "IOperation[S, F]",
    ):
        self._n = len(sizeOrArray) if isinstance(sizeOrArray, list) else sizeOrArray
        self._log = (self._n - 1).bit_length()
        self._size = 1 << self._log
        self._e = operation.e
        self._id = operation.id
        self._op = operation.op
        self._mapping = operation.mapping
        self._composition = operation.composition
        self._data = [self._e() for _ in range(2 * self._size)]
        self._lazy = [self._id() for _ in range(self._size)]
        if isinstance(sizeOrArray, list):
            for i in range(self._n):
                self._data[self._size + i] = sizeOrArray[i]
        for i in range(self._size - 1, 0, -1):
            self._pushUp(i)

    def query(self, left: int, right: int) -> "S":
        assert 0 <= left and left <= right and right <= self._n
        if left == right:
            return self._e()
        left += self._size
        right += self._size
        for i in range(self._log, 0, -1):
            if ((left >> i) << i) != left:
                self._pushDown(left >> i)
            if ((right >> i) << i) != right:
                self._pushDown(right >> i)
        leftRes, rightRes = self._e(), self._e()
        while left < right:
            if left & 1:
                leftRes = self._op(leftRes, self._data[left])
                left += 1
            if right & 1:
                right -= 1
                rightRes = self._op(self._data[right], rightRes)
            left >>= 1
            right >>= 1
        return self._op(leftRes, rightRes)

    def queryAll(self) -> "S":
        return self._data[1]

    def update(self, left: int, right: int, f: "F") -> None:
        assert 0 <= left and left <= right and right <= self._n
        if left == right:
            return
        left += self._size
        right += self._size
        for i in range(self._log, 0, -1):
            if ((left >> i) << i) != left:
                self._pushDown(left >> i)
            if ((right >> i) << i) != right:
                self._pushDown((right - 1) >> i)
        left2, right2 = left, right
        while left < right:
            if left & 1:
                self._allApply(left, f)
                left += 1
            if right & 1:
                right -= 1
                self._allApply(right, f)
            left >>= 1
            right >>= 1
        left, right = left2, right2
        for i in range(1, self._log + 1):
            if ((left >> i) << i) != left:
                self._pushUp(left >> i)
            if ((right >> i) << i) != right:
                self._pushUp((right - 1) >> i)

    def maxRight(self, left: int, key: Callable[["S"], bool]) -> int:
        assert 0 <= left and left <= self._n
        assert key(self._e())
        if left == self._n:
            return self._n
        left += self._size
        for i in range(self._log, 0, -1):
            self._pushDown(left >> i)
        res = self._e()
        while True:
            while left % 2 == 0:
                left >>= 1
            if not (key(self._op(res, self._data[left]))):
                while left < self._size:
                    self._pushDown(left)
                    left = 2 * left
                    if key(self._op(res, self._data[left])):
                        res = self._op(res, self._data[left])
                        left += 1
                return left - self._size
            res = self._op(res, self._data[left])
            left += 1
            if (left & -left) == left:
                break
        return self._n

    def minLeft(self, right: int, key: Callable[["S"], bool]) -> int:
        assert 0 <= right and right <= self._n
        assert key(self._e())
        if right == 0:
            return 0
        right += self._size
        for i in range(self._log, 0, -1):
            self._pushDown((right - 1) >> i)
        res = self._e()
        while True:
            right -= 1
            while right > 1 and (right % 2):
                right >>= 1
            if not (key(self._op(self._data[right], res))):
                while right < self._size:
                    self._pushDown(right)
                    right = 2 * right + 1
                    if key(self._op(self._data[right], res)):
                        res = self._op(self._data[right], res)
                        right -= 1
                return right + 1 - self._size
            res = self._op(self._data[right], res)
            if (right & -right) == right:
                break
        return 0

    def _pushUp(self, root: int) -> None:
        self._data[root] = self._op(self._data[2 * root], self._data[2 * root + 1])

    def _pushDown(self, root: int) -> None:
        self._allApply(2 * root, self._lazy[root])
        self._allApply(2 * root + 1, self._lazy[root])
        self._lazy[root] = self._id()

    def _allApply(self, root: int, parentLazy: "F") -> None:
        self._data[root] = self._mapping(parentLazy, self._data[root])
        if root < self._size:
            self._lazy[root] = self._composition(parentLazy, self._lazy[root])


# https://atcoder.jp/contests/abc265/tasks/abc265_g
# !012逆序对
# 给定一个长度为N的数组A，其中每个元素都是0,1,2 中的一个。
# 之后要按顺序操作Q次,如果接收到的是(1,L, R)，就要输出[Az,.. .,AR]中所有的逆序对
# 如果接收到的是(2,L,R,S,T,U)，就要将[AL,...,AR]中所有的0变成S，所有的1变成T，所有的2变成U。
# n,q<=1e5
if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    class Operation(IOperation[int, int]):
        def e(self) -> int:
            return 0

        def id(self) -> int:
            return 0

        def op(self, leftData: int, rightData: int) -> int:
            return 1

        def mapping(self, parentLazy: int, childData: int) -> int:
            return parentLazy

        def composition(self, parentLazy: int, childLazy: int) -> int:
            return 1

    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    tree = AtcoderLazySegmentTree(nums, Operation())
