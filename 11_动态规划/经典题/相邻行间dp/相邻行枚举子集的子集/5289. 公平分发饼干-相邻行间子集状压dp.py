from typing import List

INF = int(1e20)


class Solution:
    def distributeCookies(self, cookies: List[int], k: int) -> int:
        """求单个孩子在分发过程中能够获得饼干的最大总数的最小值"""
        n = len(cookies)
        subsum = [0] * (1 << n)
        for i in range(n):
            for preState in range(1 << i):
                subsum[preState | (1 << i)] = subsum[preState] + cookies[i]

        dp = [INF] * (1 << n)
        for state in range(1 << n):
            dp[state] = subsum[state]

        for i in range(1, k):
            ndp = [INF] * (1 << n)
            # 行间枚举子集的子集
            for state in range(1 << n):
                g1, g2 = state, 0
                while g1:
                    ndp[state] = min(ndp[state], max(subsum[g1], dp[g2]))
                    g1 = (g1 - 1) & state
                    g2 = state ^ g1
            dp = ndp

        return dp[-1]
