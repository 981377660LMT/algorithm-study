from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个整数 n 表示一棵 满二叉树 里面节点的数目，节点编号从 1 到 n 。根节点编号为 1 ，树中每个非叶子节点 i 都有两个孩子，分别是左孩子 2 * i 和右孩子 2 * i + 1 。

# 树中每个节点都有一个值，用下标从 0 开始、长度为 n 的整数数组 cost 表示，其中 cost[i] 是第 i + 1 个节点的值。每次操作，你可以将树中 任意 节点的值 增加 1 。你可以执行操作 任意 次。

# 你的目标是让根到每一个 叶子结点 的路径值相等。请你返回 最少 需要执行增加操作多少次。

# 注意：

# 满二叉树 指的是一棵树，它满足树中除了叶子节点外每个节点都恰好有 2 个节点，且所有叶子节点距离根节点距离相同。
# 路径值 指的是路径上所有节点的值之和。


class Solution:
    def minIncrements(self, n: int, cost: List[int]) -> int:
        def dfs(cur: int, dep: int) -> int:
            nonlocal res
            if cur * 2 > n:  # 叶子
                return cost[cur - 1]
            lSum = dfs(cur * 2, dep + 1)
            rSum = dfs(cur * 2 + 1, dep + 1)
            max_ = max(lSum, rSum)
            res += 2 * max_ - lSum - rSum
            return max_ + cost[cur - 1]

        res = 0
        dfs(1, 0)
        return res


# n = 7, cost = [1,5,2,2,3,3,1]
print(Solution().minIncrements(7, [1, 5, 2, 2, 3, 3, 1]))
# 15
# [764,1460,2664,764,2725,4556,5305,8829,5064,5929,7660,6321,4830,7055,3761]
print(
    Solution().minIncrements(
        15, [764, 1460, 2664, 764, 2725, 4556, 5305, 8829, 5064, 5929, 7660, 6321, 4830, 7055, 3761]
    )
)
