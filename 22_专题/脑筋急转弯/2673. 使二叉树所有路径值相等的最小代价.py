# 给你一个整数 n 表示一棵 满二叉树 里面节点的数目，节点编号从 1 到 n 。
# 根节点编号为 1 ，树中每个非叶子节点 i 都有两个孩子，分别是左孩子 2 * i 和右孩子 2 * i + 1 。

# 树中每个节点都有一个值，用下标从 0 开始、长度为 n 的整数数组 cost 表示，
# 其中 cost[i] 是第 i + 1 个节点的值。每次操作，你可以将树中 任意 节点的值 增加 1 。你可以执行操作 任意 次。

# !你的目标是让根到每一个 叶子结点 的路径值相等。请你返回 最少 需要执行增加操作多少次。


from typing import List


class Solution:
    def minIncrements(self, n: int, cost: List[int]) -> int:
        def dfs(cur: int) -> int:
            nonlocal res
            if cur * 2 > n:  # 叶子
                return cost[cur - 1]
            lSum = dfs(cur * 2)
            rSum = dfs(cur * 2 + 1)
            max_ = max(lSum, rSum)
            res += 2 * max_ - lSum - rSum
            return max_ + cost[cur - 1]

        res = 0
        dfs(1)
        return res


# n = 7, cost = [1,5,2,2,3,3,1]
print(Solution().minIncrements(7, [1, 5, 2, 2, 3, 3, 1]))
