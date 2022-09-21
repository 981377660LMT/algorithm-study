# 给定两个长度为n的数组A[i], B[i] (0<=A[i], B[i]<=1e9,n <=2e5)。
# !问有多少对(i,j)满足1<=i,j <=n,A[i] >= A[j],B[i] <= B[j].
# !注意(i,j)间没有大小的关系，所以(1,2)!=(2,1)
# !二维偏序 注意等号要取到 所以要一次处理完相同的数

from collections import defaultdict
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class BIT1:
    """单点修改"""

    def __init__(self, n: int):
        self.size = n
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        index += 1
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        index += 1
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    mp = defaultdict(list)
    for a, b in zip(nums1, nums2):
        a, b = a + 1, b + 1
        mp[a].append(b)

    bit = BIT1(int(1e9 + 5))
    res = 0
    for a in sorted(mp):  # 一维排序，处理完相同的数
        for b in mp[a]:
            bit.add(b, 1)
        for b in mp[a]:
            res += bit.queryRange(b, int(1e9 + 5))
    print(res)
