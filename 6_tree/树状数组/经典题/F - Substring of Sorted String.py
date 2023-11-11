#  给定一个长为n的只由小写字母组成的字符串s
#  有q次操作:
#  1 i c 将s[i]变为c (1<=i<=n)
#  2 left right 如果将s中的字符按照字典序排序,记为T
#               s的第left个字母到第right个字母是否为T的子串  (1<=left<=right<=n)

# https://www.cnblogs.com/SkyRainWind/p/17060088.html
# 维护单调性和字母个数 => 开26+1个树状数组
# !1.如何判断子串内单调不减 => 用一个标志位表示某个位置是否下降了
# !2.如何判断是子串 => 对大于开头字母和小于结尾字母的所有字母，区间内的个数必须等于[1,n]间的个数

from typing import List, Tuple, Union


def subStringOfSortedString(s: str, operations: List[Tuple[int, int, int]]) -> List[bool]:
    n = len(s)
    ords = [ord(c) for c in s]
    down = BITArray(n)  # 这个位置是否出现下降 0/1数组
    indexes = [BITArray(n) for _ in range(26)]  # 每个字母出现的位置
    for i, c in enumerate(s):
        indexes[ord(c) - 97].add(i + 1, 1)
        if i > 0 and ord(c) < ord(s[i - 1]):
            down.add(i + 1, 1)

    def add(index: int, ord_: int, delta: int) -> None:
        indexes[ord_ - 97].add(index + 1, delta)
        if index > 0 and ords[index - 1] > ord_:  # 下降了
            down.add(index + 1, delta)
        if index + 1 < n and ords[index + 1] < ord_:  # 下降了
            down.add((index + 1) + 1, delta)

    res = []
    for i, (op, a, b) in enumerate(operations):
        if op == 1:
            index, ord_ = a - 1, b
            add(index, ords[index], -1)
            ords[index] = ord_
            add(index, ord_, 1)
        else:
            left, right = a, b
            hasDown = down.queryRange(left + 1, right) != 0
            if hasDown:
                res.append(False)
                continue

            first, last = ords[left - 1] - 97, ords[right - 1] - 97
            ok = True
            for i in range(first + 1, last):
                if indexes[i].queryRange(left, right) != indexes[i].queryRange(1, n):
                    ok = False
                    break
            res.append(ok)

    return res


from typing import List, Sequence, Union


class BITArray:
    """Point Add Range Sum, 1-indexed."""

    @staticmethod
    def _build(sequence: Sequence[int]) -> List[int]:
        tree = [0] * (len(sequence) + 1)
        for i in range(1, len(tree)):
            tree[i] += sequence[i - 1]
            parent = i + (i & -i)
            if parent < len(tree):
                tree[parent] += tree[i]
        return tree

    __slots__ = ("_n", "_tree")

    def __init__(self, lenOrSequence: Union[int, Sequence[int]]):
        if isinstance(lenOrSequence, int):
            self._n = lenOrSequence
            self._tree = [0] * (lenOrSequence + 1)
        else:
            self._n = len(lenOrSequence)
            self._tree = self._build(lenOrSequence)

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, f'add index must be greater than 0, but got {index}'
        while index <= self._n:
            self._tree[index] += delta
            index += index & -index

    def query(self, right: int) -> int:
        """Query sum of [1, right]."""
        if right > self._n:
            right = self._n
        res = 0
        while right > 0:
            res += self._tree[right]
            right -= right & -right
        return res

    def queryRange(self, left: int, right: int) -> int:
        """Query sum of [left, right]."""
        return self.query(right) - self.query(left - 1)

    def __len__(self) -> int:
        return self._n

    def __repr__(self) -> str:
        nums = []
        for i in range(1, self._n + 1):
            nums.append(self.queryRange(i, i))
        return f"BITArray({nums})"


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    s = input()
    q = int(input())
    operations = []
    for _ in range(q):
        op, *args = input().split()
        if op == "1":
            operations.append((int(op), int(args[0]), ord(args[1])))
        else:
            operations.append((int(op), int(args[0]), int(args[1])))

    res = subStringOfSortedString(s, operations)
    for v in res:
        print("Yes" if v else "No")
