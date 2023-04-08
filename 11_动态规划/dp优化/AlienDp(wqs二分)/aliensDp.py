from typing import Callable, List, Tuple


def aliensDp(k: int, getDp: Callable[[int], Tuple[int, int]]) -> int:
    left, right = 1, int(1e18)  # 二分罚款
    penalty = 0
    while left <= right:
        mid = (left + right) >> 1
        cand = getDp(mid)
        if cand[1] >= k:
            penalty = mid
            left = mid + 1
        else:
            right = mid - 1
    res = getDp(penalty)
    return res[0] + penalty * k


if __name__ == "__main__":

    def max(a: int, b: int) -> int:
        return a if a > b else b

    class Solution:
        def maxProfit(self, k: int, prices: List[int]) -> int:
            """
            给定一个整数数组 prices ，
            它的第 i 个元素 prices[i] 是一支给定的股票在第 i 天的价格。
            设计一个算法来计算你所能获取的最大利润。你最多可以完成 k 笔交易。
            """
            n = len(prices)

            def getDp(penalty: int) -> Tuple[int, int]:
                dp0, dp1 = [0, 0], [-prices[0], 0]  # dp1: 不持有股票, dp2: 持有股票
                for i in range(1, n):
                    ndp0, ndp1 = dp0, dp1
                    cand = dp1[0] + prices[i] - penalty
                    if cand >= dp0[0]:
                        ndp0 = [cand, max(dp0[1], dp1[1] + 1)]  # !注意让使用次数最大
                    cand = dp0[0] - prices[i]
                    if cand >= dp1[0]:
                        ndp1 = [cand, max(dp0[1], dp1[1])]  # !注意让使用次数最大
                    dp0, dp1 = ndp0, ndp1
                return tuple(dp0)

            return aliensDp(k, getDp)
