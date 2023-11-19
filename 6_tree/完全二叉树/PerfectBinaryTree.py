from typing import List


class PerfectBinaryTree:
    """
    完全二叉树.
    根节点编号为1,左子节点编号为2*x,右子节点编号为2*x+1,父结点编号为x>>1.
    """

    @staticmethod
    def depth(u: int) -> int:
        if u == 0:
            return 0
        return u.bit_length() - 1

    @staticmethod
    def lca(u: int, v: int) -> int:
        if u == v:
            return u
        if u > v:
            u, v = v, u
        depth1, depth2 = PerfectBinaryTree.depth(u), PerfectBinaryTree.depth(v)
        diff = u ^ (v >> (depth2 - depth1))
        if diff == 0:
            return u
        len_ = diff.bit_length()
        return u >> len_

    @staticmethod
    def dist(u: int, v: int) -> int:
        return (
            PerfectBinaryTree.depth(u)
            + PerfectBinaryTree.depth(v)
            - 2 * PerfectBinaryTree.depth(PerfectBinaryTree.lca(u, v))
        )


if __name__ == "__main__":
    # https://leetcode.cn/problems/cycle-length-queries-in-a-tree/description/
    class Solution:
        def cycleLengthQueries(self, n: int, queries: List[List[int]]) -> List[int]:
            res = [0] * len(queries)
            for i, (root1, root2) in enumerate(queries):
                res[i] = PerfectBinaryTree.dist(root1, root2) + 1
            return res
