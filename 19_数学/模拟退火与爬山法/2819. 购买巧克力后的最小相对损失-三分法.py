# https://leetcode.cn/problems/minimum-relative-loss-after-buying-chocolates/

# 现给定一个整数数组 prices，表示巧克力的价格；以及一个二维整数数组 queries，其中 queries[i] = [threshold_i, count_i]。
# Alice 和 Bob 去买巧克力，Alice 提出了一种付款方式，而 Bob 同意了。
# 对于每个 queries[i] ，它的条件如下：
# !如果一块巧克力的价格 小于等于 threshold_i，那么 Bob 为它付款。
# !否则，Bob 为其中 threshold_i 部分付款，而 Alice 为 剩余 部分付款。
# Bob 想要选择 恰好 count_i 块巧克力，使得他的 相对损失最小 。
# 更具体地说，如果总共 Alice 付款了 ai，Bob 付款了 bi，那么 Bob 希望最小化 bi - ai。
# 返回一个整数数组 ans，其中 ans[i] 是 Bob 在 queries[i] 中可能的 最小相对损失 。


# !从两个有序数组中找出最小的 m 个数的和 => 三分法

from itertools import accumulate
from typing import List
from 三分法求凸函数极值 import minimize


class Solution:
    def minimumRelativeLosses(self, prices: List[int], queries: List[List[int]]) -> List[int]:
        def fun(preLen: int) -> int:
            """从prices的前缀中选择preLen个数时,损失的最小值."""
            sufLen = count - preLen
            return preSum[preLen] + (2 * threshold * sufLen - sufSum[sufLen])

        prices.sort()
        preSum = [0] + list(accumulate(prices))
        sufSum = [0] + list(accumulate(prices[::-1]))
        res = [0] * len(queries)
        for i, (threshold, count) in enumerate(queries):
            res[i] = minimize(fun, 0, count)
        return res


print(Solution().minimumRelativeLosses(prices=[1, 9, 22, 10, 19], queries=[[18, 4], [5, 2]]))
