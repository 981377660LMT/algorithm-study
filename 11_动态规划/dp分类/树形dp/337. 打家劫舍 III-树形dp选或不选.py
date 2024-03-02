# https://leetcode.cn/problems/house-robber-iii/description/
# 337. 打家劫舍 III-树形dp选或不选
# 如果 两个直接相连的房子在同一天晚上被打劫 ，房屋将自动报警。
# 给定二叉树的 root 。返回 在不触动警报的情况下 ，小偷能够盗取的最高金额 。


from typing import Optional, Tuple

INF = int(1e20)


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution:
    def rob(self, root: Optional[TreeNode]) -> int:
        def dfs(cur: Optional["TreeNode"]) -> Tuple[int, int]:
            """不选/选"""
            if not cur:
                return 0, 0
            res0, res1 = 0, cur.val
            nexts = [cur.left, cur.right]
            nexts = [x for x in nexts if x]
            for next_ in nexts:
                a, b = dfs(next_)
                res0, res1 = res0 + max(a, b), res1 + a
            return res0, res1

        return max(dfs(root))
