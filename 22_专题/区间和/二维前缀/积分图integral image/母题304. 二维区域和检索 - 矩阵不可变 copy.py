from typing import List


class NumMatrix:
    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])
        # pre[r][c] == sum(A[:r, :c])
        pre = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                pre[r + 1][c + 1] = A[r][c] + pre[r][c + 1] + pre[r + 1][c] - pre[r][c]
        self.pre = pre

    def sumRegion(self, r0: int, c0: int, r1: int, c1: int) -> int:
        """ return sum(A[r0:r1+1, c0:c1+1]) """
        # assert 0<=r0<=r1<m, 0<=c0<=c1<n
        return (
            self.pre[r1 + 1][c1 + 1]
            - self.pre[r1 + 1][c0]
            - self.pre[r0][c1 + 1]
            + self.pre[r0][c0]
        )

