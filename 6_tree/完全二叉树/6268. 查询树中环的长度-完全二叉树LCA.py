"""完全二叉树LCA"""

# 给你一个整数 n ，表示你有一棵含有 2n - 1 个节点的 完全二叉树 。
# !根节点的编号是 1 ，树中编号在[1, 2n - 1 - 1] 之间，编号为 val 的节点都有两个子节点，满足：

# 左子节点的编号为 2 * val
# 右子节点的编号为 2 * val + 1
# 给你一个长度为 m 的查询数组 queries ，它是一个二维整数数组，其中 queries[i] = [ai, bi] 。对于每个查询，求出以下问题的解：

# 在节点编号为 ai 和 bi 之间添加一条边。
# 求出图中环的长度。
# 删除节点编号为 ai 和 bi 之间新添加的边。
# 注意：

# 环 是开始和结束于同一节点的一条路径，路径中每条边都只会被访问一次。
# 环的长度是环中边的数目。
# 在树中添加额外的边后，两个点之间可能会有多条边。
# 请你返回一个长度为 m 的数组 answer ，其中 answer[i] 是第 i 个查询的结果。


# !1. 完全二叉树的结点深度可以有bit_length()求出
# !2. 完全二叉树结点的LCA可以由线性上跳/根据深度差上跳求出
from typing import List


class Solution:
    def cycleLengthQueries(self, n: int, queries: List[List[int]]) -> List[int]:
        def getDepth(root: int) -> int:
            return root.bit_length() - 1

        def getLCA(a: int, b: int) -> int:
            """
            https://leetcode.cn/problems/cycle-length-queries-in-a-tree/solutions/1/shinei-by-vclip-ch6n/
            1. 一个结点从最高位的 1开始,往下遍历每一个二进制位,得到的就是从根节点出发走到这个结点的路径选择,
               0代表向左,1代表向右;
            2. 求出两个结点的深度，然后把较深结点向上移动，对齐到另一个节点，让两者的二进制位对应上;
            3. 之后，对两个二进制数求异或，异或值就代表了两者选择路径时走的方向不同的位置，
               异或为0,说明a就是LCA;否则上跳
            """
            if a == b:
                return a
            if a > b:
                a, b = b, a
            depthA, depthB = a.bit_length() - 1, b.bit_length() - 1
            diff = a ^ (b >> (depthB - depthA))
            if diff == 0:
                return a
            len_ = diff.bit_length()
            return a >> len_

        res = [0] * len(queries)
        for i, (root1, root2) in enumerate(queries):
            lca = getLCA(root1, root2)
            res[i] = getDepth(root1) + getDepth(root2) - 2 * getDepth(lca) + 1
        return res


print(Solution().cycleLengthQueries(n=3, queries=[[5, 3], [4, 7], [2, 3]]))
