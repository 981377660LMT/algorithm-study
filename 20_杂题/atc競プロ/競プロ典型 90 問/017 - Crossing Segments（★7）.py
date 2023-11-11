# 圆周上n个点(1-n)上有m条线 求相交的直线组数
# 3≤N≤3e5
# 2≤M≤3e5

# !1.最初は`全探索`を考えることが大切　すなわちcombinationsを使う O(M^2)
# !2.余事像を考える O(N*M)
# !3.BITで高速化 O(N+MlogN)
# !按顺序遍历左端点 然后看之前的右端点多少个落在当前区间 再更新当前右端点到bit

from collections import defaultdict

import sys


class BIT1:
    """单点修改

    https://github.com/981377660LMT/algorithm-study/blob/master/6_tree/%E6%A0%91%E7%8A%B6%E6%95%B0%E7%BB%84/%E7%BB%8F%E5%85%B8%E9%A2%98/BIT.py
    """

    def __init__(self, n: int):
        self.size = n
        self.tree = defaultdict(int)

    @staticmethod
    def _lowbit(index: int) -> int:
        return index & -index

    def add(self, index: int, delta: int) -> None:
        if index <= 0:
            raise ValueError('index 必须是正整数')
        while index <= self.size:
            self.tree[index] += delta
            index += self._lowbit(index)

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree[index]
            index -= self._lowbit(index)
        return res

    def queryRange(self, left: int, right: int) -> int:
        if left > right:
            return 0
        return self.query(right) - self.query(left - 1)


sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)
n, m = map(int, input().split())
bit = BIT1(n + 10)  # 保存右端点

adjMap = defaultdict(list)
for _ in range(m):
    left, right = map(int, input().split())
    adjMap[left].append(right)

res = 0
for i in range(1, n + 1):
    for j in adjMap[i]:
        res += bit.queryRange(i + 1, j - 1)  # 有之前的右端点落在这个区间里
    for j in adjMap[i]:
        bit.add(j, 1)

print(res)
