# https://www.acwing.com/video/786/
# 给你一大串数字（编号为 1 到 N，大小可不一定哦！），
# 在你看过一遍之后，它便消失在你面前，随后问题就出现了，
# 给你 M 个询问，每次询问就给你两个数字 A,B，
# 要求你瞬间就说出属于 A 到 B 这段区间内的最大数。
from math import ceil, log, log2
from typing import Any, Generic, List, TypeVar

T = TypeVar('T', int, float)


class SparseTable(Generic[T]):
    def __init__(self, nums: List[T]):
        n, upper = len(nums), ceil(log2(len(nums))) + 1
        self._nums = nums
        # dp[i][j]表示从i开始长度为1<<j的区间中的最值
        # dp[i][j]=max(dp[i][j-1],dp[i+1<<(j-1)][j-1]) 前一半后一半
        self._dp1: List[List[Any]] = [[0] * upper for _ in range(n)]
        self._dp2: List[List[Any]] = [[0] * upper for _ in range(n)]
        for i, num in enumerate(nums):
            self._dp1[i][0] = num
            self._dp2[i][0] = num
        for j in range(1, upper):
            for i in range(n):
                if i + (1 << (j - 1)) >= n:
                    break
                self._dp1[i][j] = max(self._dp1[i][j - 1], self._dp1[i + (1 << (j - 1))][j - 1])
                self._dp2[i][j] = min(self._dp2[i][j - 1], self._dp2[i + (1 << (j - 1))][j - 1])

    def query(self, left: int, right: int, *, ismax=True) -> T:
        """[left,right]区间的最大值"""
        assert 0 <= left <= right < len(self._nums)

        # 找到最大的k使得1<<k<=length 1<<(k+1)>length
        # dp[left][k]和dp[right-(1<<k)+1][k]可以覆盖[left,right]
        k = int(log(right - left + 1) / log(2))
        if ismax:
            return max(self._dp1[left][k], self._dp1[right - (1 << k) + 1][k])
        else:
            return min(self._dp2[left][k], self._dp2[right - (1 << k) + 1][k])


n = int(input())
nums = list(map(int, input().split()))
st = SparseTable(nums)

m = int(input())
for _ in range(m):
    a, b = map(int, input().split())
    print(st.query(a - 1, b - 1))

