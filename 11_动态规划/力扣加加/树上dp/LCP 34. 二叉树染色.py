from typing import List, Optional


class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


# https://leetcode-cn.com/problems/er-cha-shu-ran-se-UGC/solution/c-python3-dong-tai-gui-hua-dp-shu-xing-d-fgcm/
# 二叉树上每个蓝色相连部分的结点个数不能超过 k 个，求所有染成蓝色的结点价值总和最大是多少？
# k<=10
class Solution:
    def maxValue(self, root: 'TreeNode', k: int) -> int:
        def dfs(root: Optional[TreeNode]) -> List[int]:
            if not root:
                return [0] * (k + 1)

            left = dfs(root.left)
            right = dfs(root.right)

            # dp[i]以root为根结点的树，染色后，连着i个蓝色结点的最大值
            dp = [0] * (k + 1)

            # 不是蓝色
            for i in range(k + 1):
                for j in range(k + 1):
                    dp[0] = max(dp[0], left[i] + right[j])

            # 是蓝色 i+j < k
            for i in range(k):
                for j in range(k - i):
                    dp[i + j + 1] = max(dp[i + j + 1], left[i] + right[j] + root.val)

            return dp

        return max(dfs(root))


# root = [5,2,3,4], k = 2
# 结点 5、3、4 染成蓝色，获得最大的价值 5+3+4=12
