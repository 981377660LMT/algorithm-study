"""
返回一个大小为 n 的数组 ans ,
其中 ans[i]是离节点 i 最近的祖先节点且满足 nums[i] 和 nums[ans[i]] 是 互质的 ,
如果不存在这样的祖先节点,ans[i] 为 -1 。
1 <= nums[i] <= 50
1 <= n <= 1e5


离线查询、小值域搜答案
"""

from math import gcd
from typing import List


class Solution:
    def getCoprimes(self, nums: List[int], edges: List[List[int]]) -> List[int]:
        def dfs(cur: int, pre: int, dep: int) -> None:
            resDep, resId = -1, -1
            for i in range(1, 51):
                if not stack[i] or gcd(i, nums[cur]) != 1:
                    continue
                candDep, candId = stack[i][-1]
                if candDep > resDep:
                    resDep, resId = candDep, candId
            res[cur] = resId

            for next in adjList[cur]:
                if next == pre:
                    continue
                stack[nums[cur]].append((dep, cur))
                dfs(next, cur, dep + 1)
                stack[nums[cur]].pop()

        n = len(nums)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        stack = [[] for _ in range(51)]
        res = [-1] * n
        dfs(0, -1, 0)
        return res


print(Solution().getCoprimes(nums=[2, 3, 3, 2], edges=[[0, 1], [1, 2], [1, 3]]))
