from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 数列
# X に対し、
# f(X)= (
# X を回文にするために変更する必要のある要素の個数の最小値 ) とします。

# 与えられた長さ
# N の数列
# A の全ての 連続 部分列
# X に対する
# f(X) の総和を求めてください。

# 但し、長さ
# m の数列
# X が回文であるとは、全ての
# 1≤i≤m を満たす整数
# i について、
# X の
# i 項目と
# m+1−i 項目が等しいことを指します。

from typing import List, Sequence, Union


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, 'index must be greater than 0'
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)

    def bisectLeft(self, k: int) -> int:
        """返回第一个前缀和大于等于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) < k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def bisectRight(self, k: int) -> int:
        """返回第一个前缀和大于k的位置pos

        1 <= pos <= self.size + 1
        """
        curSum, pos = 0, 0
        for i in range(self.bit, -1, -1):
            nextPos = pos + (1 << i)
            if nextPos <= self.size and curSum + self.tree.get(nextPos, 0) <= k:
                pos = nextPos
                curSum += self.tree.get(pos, 0)
        return pos + 1

    def __repr__(self) -> str:
        preSum = []
        for i in range(self.size):
            preSum.append(self.query(i))
        return str(preSum)

    def __len__(self) -> int:
        return self.size


# 值域树状数组
if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))

    # 子串:关注开头和结尾
    mp = defaultdict(list)
    for i, char in enumerate(nums):
        mp[char].append(i)
    # print(mp)
    # defaultdict(<class 'list'>, {5: [0], 2: [1, 3, 4], 1: [2]})

    # 所有字符都不想等时需要的操作数
    res = 0
    for len_ in range(1, n + 1):
        res += (len_ // 2) * (n - len_ + 1)
    for group in mp.values():
        if len(group) == 1:
            continue

        # 右侧比自己距离左侧小的距离之和 + 右侧距离大于等于自己的位置个数*左侧距离
        bit1, bit2 = BIT1(n + 10), BIT1(n + 10)  # 计数/距离
        distToLeft = [num + 1 for num in group]
        distToRight = [n - num for num in group]
        for i, v in enumerate(distToRight):
            bit1.add(v, 1)
            bit2.add(v, v)

        k = len(group)
        for i in range(k):
            bit1.add(distToRight[i], -1)
            bit2.add(distToRight[i], -distToRight[i])
            lDist = distToLeft[i]
            rightBigger = bit1.queryRange(lDist, n + 1)
            rightSmallSum = bit2.queryRange(0, lDist - 1)
            res -= rightBigger * lDist + rightSmallSum
    print(res)
