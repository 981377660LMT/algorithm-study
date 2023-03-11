# 子矩阵(窗口)中位数最小值
# https://atcoder.jp/contests/abc203/tasks/abc203_d

# 输入 n k (1≤k≤n≤800) 和一个 n*n 的矩阵，元素范围 [0,1e9]。
# !定义 k*k 子矩阵的中位数为子矩阵的第 floor(k*k/2)+1 大的数。
# !输出中位数的最小值。
# 注：「第 x 大」中的 x 从 1 开始。

# 二分答案+二维前缀和加速计算  O(n^2log1e9)

from typing import List
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class M:
    """二维前缀和统计区域内小于等于 upper 的元素个数"""

    def __init__(self, A: List[List[int]], upper: int):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = (
                    (A[r][c] <= upper) + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
                )
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        M.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


def solve(grid: List[List[int]], k: int) -> int:
    def check(mid: int) -> bool:
        """是否存在一个 k*k 的子矩阵，中位数小于等于 mid"""
        preSum = M(grid, mid)
        for r in range(n - k + 1):
            for c in range(n - k + 1):
                smaller = preSum.sumRegion(r, c, r + k - 1, c + k - 1)
                if smaller >= (k * k + 1) // 2:
                    return True
        return False

    left, right = 0, int(1e9 + 10)
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            right = mid - 1
        else:
            left = mid + 1
    return left


if __name__ == "__main__":
    n, k = map(int, input().split())
    grid = [list(map(int, input().split())) for _ in range(n)]
    print(solve(grid, k))
