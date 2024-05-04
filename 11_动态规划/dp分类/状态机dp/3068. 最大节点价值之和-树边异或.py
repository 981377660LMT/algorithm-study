# 3068. 最大节点价值之和-树边异或
# https://leetcode.cn/problems/find-the-maximum-sum-of-node-values/
# Alice 想 最大化 树中所有节点价值之和。为了实现这一目标，Alice 可以执行以下操作 任意 次（包括 0 次）：
# 选择连接节点 u 和 v 的边 [u, v] ，并将它们的值更新为：
# nums[u] = nums[u] XOR k
# nums[v] = nums[v] XOR k
# !请你返回 Alice 通过执行以上操作 任意次 后，可以得到所有节点 价值之和 的 最大值 。
#
# !等价于：选择 nums 中的偶数个元素，把这些数都异或 k，数组的最大元素和是多少
# !dp[i][0/1]表示前i个数中选择了偶数/奇数个数时的最大和 -> 奇偶dp

from typing import List

INF = int(1e18)


class Solution:
    def maximumValueSum(self, nums: List[int], k: int, edges: List[List[int]]) -> int:
        dp0, dp1 = 0, -INF
        for v in nums:
            ndp0, ndp1 = max(dp0 + v, dp1 + (v ^ k)), max(dp1 + v, dp0 + (v ^ k))
            dp0, dp1 = ndp0, ndp1
        return dp0
