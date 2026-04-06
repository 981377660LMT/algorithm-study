# 3864. 划分二进制字符串的最小费用
# https://leetcode.cn/problems/minimum-cost-to-partition-a-binary-string/description/


from itertools import accumulate


class Solution:
    def minCost(self, s: str, encCost: int, flatCost: int) -> int:
        n = len(s)
        presum = list(accumulate(map(int, s), initial=0))

        def dfs(start: int, end: int) -> int:
            # 不拆分
            x = presum[end] - presum[start]
            res = (end - start) * x * encCost if x else flatCost
            # 拆分
            if (end - start) % 2 == 0:
                mid = (start + end) // 2
                res = min(res, dfs(start, mid) + dfs(mid, end))
            return res

        return dfs(0, n)
