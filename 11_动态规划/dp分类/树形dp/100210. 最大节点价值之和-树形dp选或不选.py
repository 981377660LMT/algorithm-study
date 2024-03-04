# 100210. 最大节点价值之和-树形dp选或不选
# https://leetcode.cn/problems/find-the-maximum-sum-of-node-values/description/
# 给定一个无向带权树和一个正整数k.可以选择任意一条边，将两个端点异或上k.
# 最大化所有节点的价值之和
#
# !解法1：在每个节点枚举该点是否需要父节点做异或 k 操作
#
# !解法2：树上异或 => 联想到差分 => 联想到距离


from typing import List, Tuple

INF = int(1e20)


class Solution:
    def maximumValueSum(self, nums: List[int], k: int, edges: List[List[int]]) -> int:
        n = len(nums)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        def dfs(cur: int, pre: int) -> Tuple[int, int]:
            """该点"不与父节点一起异或/与父节点一起异或"的最大价值"""
            res0, res1 = nums[cur], nums[cur] ^ k
            for next_ in adjList[cur]:
                if next_ == pre:
                    continue
                a, b = dfs(next_, cur)
                res0, res1 = max(res0 + a, res1 + b), max(res0 + b, res1 + a)
            return res0, res1

        res = dfs(0, -1)
        return res[0]

    def maximumValueSum2(self, nums: List[int], k: int, edges: List[List[int]]) -> int:
        """注意到树上异或性质,等价于可以将偶数个结点异或上k,最大化所有节点的价值之和."""
        dp0, dp1 = 0, -INF  # 前i个元素中，有偶数个/奇数个元素异或上k的最大价值
        for num in nums:
            dp0, dp1 = max(dp0 + num, dp1 + (num ^ k)), max(dp0 + (num ^ k), dp1 + num)
        return dp0
