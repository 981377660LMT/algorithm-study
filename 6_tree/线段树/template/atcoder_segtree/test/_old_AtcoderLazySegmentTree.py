# Deprecated

# !由于lazy模板通用性 效率不如自己维护数组的线段树
# !注意如果是单点查询,可以去掉所有pushUp函数逻辑
# !如果是单点修改,可以去掉所有懒标记逻辑

import sys
from typing import Callable, Generic, List, TypeVar, Union, overload

S = TypeVar("S")
"""线段树维护的值的类型"""

F = TypeVar("F")
"""懒标记的类型/更新操作的类型"""

E = Callable[[], S]
"""线段树的值的幺元"""

Id = Callable[[], F]
"""更新操作/懒标记的幺元"""

Op = Callable[[S, S], S]
"""线段树的合并操作"""

Mapping = Callable[[F, S], S]
"""父结点的懒标记更新子结点的值"""

Composition = Callable[[F, F], F]
"""父结点的懒标记更新子结点的懒标记"""

# check python version for Protocol
if sys.version_info >= (3, 8):
    from typing import Protocol
else:

    class DummyProtocol(Generic[S, F]):
        __slots__ = ()
        ...

    Protocol = DummyProtocol


class IOperation(Protocol[S, F]):
    """线段树的操作的接口"""

    def e(self) -> S:
        """线段树的值的幺元"""
        ...

    def id(self) -> "F":
        """更新操作/懒标记的幺元"""
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


class AtcoderLazySegmentTree(Generic[S, F]):

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

    @overload
    def __init__(
        self,
        sizeOrArray: Union[int, List["S"]],
        *,
        e: "E[S]",
        id: "Id[F]",
        op: "Op[S]",
        mapping: "Mapping[S, F]",
        composition: "Composition[F]",
    ) -> None:

        ...

    @overload
    def __init__(self, sizeOrArray: Union[int, List["S"]], *, operation: "IOperation[S, F]"):
        ...

    def __init__(self, sizeOrArray: Union[int, List["S"]], **kwargs):
        self._n = len(sizeOrArray) if isinstance(sizeOrArray, list) else sizeOrArray
        self._log = (self._n - 1).bit_length()
        self._size = 1 << self._log

        if "operation" in kwargs:
            operation = kwargs["operation"]
            self._e = operation.e
            self._id = operation.id
            self._op = operation.op
            self._mapping = operation.mapping
            self._composition = operation.composition
        else:
            self._e = kwargs["e"]
            self._id = kwargs["id"]
            self._op = kwargs["op"]
            self._mapping = kwargs["mapping"]
            self._composition = kwargs["composition"]

        self._data = [self._e() for _ in range(2 * self._size)]
        self._lazy = [self._id() for _ in range(self._size)]
        if isinstance(sizeOrArray, list):
            for i in range(self._n):
                self._data[self._size + i] = sizeOrArray[i]
        for i in range(self._size - 1, 0, -1):
            self._pushUp(i)

    def query(self, left: int, right: int) -> "S":
        """
        查询切片 `[left:right]` 内的值

        `0 <= left <= right <= n`

        alias: prod
        """
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
        """
        更新切片 `[left:right]` 内的值

        `0 <= left <= right <= n`

        alias: apply
        """
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
        """
        树上二分查询最大的 `right` 使得切片 `[left:right]` 内的值满足 `key`
        """
        assert 0 <= left and left <= self._n
        # assert key(self._e())  # FIXME
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
        """
        树上二分查询最小的 `left` 使得切片 `[left:right]` 内的值满足 `key`
        """
        assert 0 <= right and right <= self._n
        # assert key(self._e())  # FIXME
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

    #  propagate
    def _allApply(self, root: int, parentLazy: "F") -> None:
        self._data[root] = self._mapping(parentLazy, self._data[root])
        # !叶子结点不需要更新lazy
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

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    MOD = 998244353
    INF = int(4e18)

    # 线段树维护区间的6个值
    # ![0的个数,1的个数,2的个数,10逆序对的个数,20逆序对的个数,21逆序对的个数]
    # class Operation(IOperation[List[int], List[int]]):
    #     def e(self) -> List[int]:
    #         return [0] * 6

    #     def id(self) -> List[int]:
    #         return [0, 1, 2]

    #     def op(self, leftData: List[int], rightData: List[int]) -> List[int]:
    #         res = [0] * 6
    #         res[0] = leftData[0] + rightData[0]
    #         res[1] = leftData[1] + rightData[1]
    #         res[2] = leftData[2] + rightData[2]
    #         res[3] = leftData[3] + rightData[3] + leftData[1] * rightData[0]
    #         res[4] = leftData[4] + rightData[4] + leftData[2] * rightData[0]
    #         res[5] = leftData[5] + rightData[5] + leftData[2] * rightData[1]
    #         return res

    #     def mapping(self, parentLazy: List[int], childData: List[int]) -> List[int]:
    #         res = [0] * 6
    #         res[parentLazy[0]] += childData[0]
    #         res[parentLazy[1]] += childData[1]
    #         res[parentLazy[2]] += childData[2]
    #         counter = [[0] * 3 for _ in range(3)]  # !counter[i][j]表示(i,j)的对数
    #         counter[parentLazy[1]][parentLazy[0]] += childData[3]
    #         counter[parentLazy[2]][parentLazy[0]] += childData[4]
    #         counter[parentLazy[2]][parentLazy[1]] += childData[5]
    #         counter[parentLazy[0]][parentLazy[1]] += childData[0] * childData[1] - childData[3]
    #         counter[parentLazy[0]][parentLazy[2]] += childData[0] * childData[2] - childData[4]
    #         counter[parentLazy[1]][parentLazy[2]] += childData[1] * childData[2] - childData[5]
    #         res[3] = counter[1][0]
    #         res[4] = counter[2][0]
    #         res[5] = counter[2][1]
    #         return res

    #     def composition(self, parentLazy: List[int], childLazy: List[int]) -> List[int]:
    #         res = [0] * 3
    #         res[0] = parentLazy[childLazy[0]]
    #         res[1] = parentLazy[childLazy[1]]
    #         res[2] = parentLazy[childLazy[2]]
    #         return res
    def e() -> List[int]:
        return [0] * 6

    def id() -> List[int]:
        return [0, 1, 2]

    def op(leftData: List[int], rightData: List[int]) -> List[int]:
        res = [0] * 6
        res[0] = leftData[0] + rightData[0]
        res[1] = leftData[1] + rightData[1]
        res[2] = leftData[2] + rightData[2]
        res[3] = leftData[3] + rightData[3] + leftData[1] * rightData[0]
        res[4] = leftData[4] + rightData[4] + leftData[2] * rightData[0]
        res[5] = leftData[5] + rightData[5] + leftData[2] * rightData[1]
        return res

    def mapping(parentLazy: List[int], childData: List[int]) -> List[int]:
        res = [0] * 6
        res[parentLazy[0]] += childData[0]
        res[parentLazy[1]] += childData[1]
        res[parentLazy[2]] += childData[2]
        counter = [[0] * 3 for _ in range(3)]
        counter[parentLazy[1]][parentLazy[0]] += childData[3]
        counter[parentLazy[2]][parentLazy[0]] += childData[4]
        counter[parentLazy[2]][parentLazy[1]] += childData[5]
        counter[parentLazy[0]][parentLazy[1]] += childData[0] * childData[1] - childData[3]
        counter[parentLazy[0]][parentLazy[2]] += childData[0] * childData[2] - childData[4]
        counter[parentLazy[1]][parentLazy[2]] += childData[1] * childData[2] - childData[5]
        res[3] = counter[1][0]
        res[4] = counter[2][0]
        res[5] = counter[2][1]
        return res

    def composition(parentLazy: List[int], childLazy: List[int]) -> List[int]:
        res = [0] * 3
        res[0] = parentLazy[childLazy[0]]
        res[1] = parentLazy[childLazy[1]]
        res[2] = parentLazy[childLazy[2]]
        return res

    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    init = [[0] * 6 for _ in range(n)]
    for i in range(n):
        init[i][nums[i]] = 1

    # !4300ms
    # tree = AtcoderLazySegmentTree(init, operation=Operation())

    # !4100ms
    tree = AtcoderLazySegmentTree(init, e=e, id=id, op=op, mapping=mapping, composition=composition)
    for _ in range(q):
        kind, *rest = map(int, input().split())
        if kind == 1:
            left, right = rest
            left -= 1
            print(sum(tree.query(left, right)[3:]))
        else:
            left, right, s, t, u = rest
            left -= 1
            tree.update(left, right, [s, t, u])  # 0=>s,1=>t,2=>u
