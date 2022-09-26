from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 我们认为当销售员同一天推销的产品数目大于 8 个的时候，那么这一天就是「成功销售的一天」。
# 所谓「销售出色区间」，意味在这段时间内，「成功销售的天数」是严格 大于「未成功销售的天数」。
# 请你返回「销售出色区间」的最大长度。


class Solution:
    def longestESR(self, sales: List[int]) -> int:
        pre = dict({0: -1})
        res = cursum = 0

        for i, h in enumerate(sales):
            cursum += 1 if h > 8 else -1
            if cursum > 0:
                res = i + 1
            if cursum - 1 in pre:
                res = max(res, i - pre[cursum - 1])
            pre.setdefault(cursum, i)
        return res
