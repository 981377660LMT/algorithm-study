# ConvexHullTrickLiDeque
# 单调队列维护凸包


from collections import deque
from typing import List, Tuple


Line = Tuple[int, int, int]  # k, b, id
INF = int(4e18)


class ConvexHullTrickDeque:
    __slots__ = "_isMin", "_dq"

    def __init__(self, isMin: bool) -> None:
        self._isMin = isMin
        self._dq = deque()

    def addLine(self, k: int, b: int, id: int = -1) -> None:
        """
        追加一条直线 k*x + b.
        需要保证斜率k是单调递增或者是单调递减的.
        """
        if not self._isMin:
            k, b = -k, -b
        line = (k, b, id)
        if not self._dq:
            self._dq.appendleft(line)
            return
        dp = self._dq
        if dp[0][0] <= k:
            if dp[0][0] == k:
                if dp[0][1] <= b:
                    return
                dp.popleft()
            while len(dp) >= 2 and self._check(line, dp[0], dp[1]):
                dp.popleft()
            dp.appendleft(line)
        else:
            if dp[-1][0] == k:
                if dp[-1][1] <= b:
                    return
                dp.pop()
            while len(dp) >= 2 and self._check(dp[-2], dp[-1], line):
                dp.pop()
            dp.append(line)

    def query(self, x: int) -> Tuple[int, int]:
        """
        O(logn) 查询 k*x + b 的最小(大)值以及对应的直线id.
        如果不存在直线,返回的id为-1.
        """
        if not self._dq:
            res, lineId = INF, -1
            if not self._isMin:
                res = -res
            return res, lineId

        dp = self._dq
        left, right = -1, len(dp) - 1
        while left + 1 < right:
            mid = (left + right) >> 1
            a, _ = self._getY(dp[mid], x)
            b, _ = self._getY(dp[mid + 1], x)
            if a >= b:
                left = mid
            else:
                right = mid
        res, lineId = self._getY(dp[right], x)
        if not self._isMin:
            res = -res
        return res, lineId

    def queryMonotoneInc(self, x: int) -> Tuple[int, int]:
        """
        O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
        需要保证x是单调递增的.
        如果不存在直线,返回的id为-1.
        """
        if not self._dq:
            res, lineId = INF, -1
            if not self._isMin:
                res = -res
            return res, lineId

        dp = self._dq
        while len(dp) >= 2:
            a, _ = self._getY(dp[0], x)
            b, _ = self._getY(dp[1], x)
            if a < b:
                break
            dp.popleft()
        res, lineId = self._getY(dp[0], x)
        if not self._isMin:
            res = -res
        return res, lineId

    def queryMonotoneDec(self, x: int) -> Tuple[int, int]:
        """
        O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
        需要保证x是单调递减的.
        如果不存在直线,返回的id为-1.
        """
        if not self._dq:
            res, lineId = INF, -1
            if not self._isMin:
                res = -res
            return res, lineId

        dp = self._dq
        while len(dp) >= 2:
            a, _ = self._getY(dp[-1], x)
            b, _ = self._getY(dp[-2], x)
            if a < b:
                break
            dp.pop()
        res, lineId = self._getY(dp[-1], x)
        if not self._isMin:
            res = -res
        return res, lineId

    def _check(self, a: Line, b: Line, c: Line) -> bool:
        if b[1] == a[1] or c[1] == b[1]:
            return self._sign(b[0] - a[0]) * self._sign(c[1] - b[1]) >= self._sign(
                c[0] - b[0]
            ) * self._sign(b[1] - a[1])
        return (b[0] - a[0]) * self._sign(c[1] - b[1]) * abs(c[1] - b[1]) >= (
            c[0] - b[0]
        ) * self._sign(b[1] - a[1]) * abs(b[1] - a[1])

    def _getY(self, line: Line, x: int) -> Tuple[int, int]:
        return line[0] * x + line[1], line[2]

    def _sign(self, x: int) -> int:
        if x == 0:
            return 0
        return 1 if x > 0 else -1


if __name__ == "__main__":
    # 3221. 最大数组跳跃得分 II
    # https://leetcode.cn/problems/maximum-array-hopping-score-ii/
    # dp[j]=max(dp[j],dp[i]+(j-i)*nums[j])
    # !dp[j]=max(dp[j],-i*nums[j]+dp[i]+j*nums[j])
    # !dp过程中将直线(-i,dp[i])不断加入到CHT中，查询时查询x=nums[j]时的最大值即可
    class Solution:
        def maxScore(self, nums: List[int]) -> int:
            n = len(nums)
            dp = [0] * n
            cht = ConvexHullTrickDeque(isMin=False)
            cht.addLine(0, 0)
            for j, v in enumerate(nums):
                cur, _ = cht.queryMonotoneInc(v)
                dp[j] = cur + v * j
                cht.addLine(-j, dp[j])
            return dp[n - 1]
