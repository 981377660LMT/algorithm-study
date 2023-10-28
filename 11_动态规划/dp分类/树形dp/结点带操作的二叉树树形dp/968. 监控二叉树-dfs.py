"""监控二叉树 放照相机

给定一个二叉树，我们在树的节点上安装摄像头。
节点上的每个摄影头都可以监视其父对象、自身及其直接子对象。
计算监控树的所有节点所需的最小摄像头数量。
"""


from functools import lru_cache
from typing import List, Optional, Tuple

INF = int(1e18)


class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


@lru_cache(None)
def gen(curPut: bool, prePut: bool) -> List[Tuple[int, int]]:
    """生成左右儿子所有可能的转移状态"""
    if curPut or prePut:
        return [(0, 0), (1, 0), (0, 1), (1, 1)]
    return [(1, 0), (0, 1), (1, 1)]


class Solution:
    def minCameraCover(self, root: "TreeNode") -> int:
        @lru_cache(None)
        def dfs(root: Optional["TreeNode"], curPut: bool, prePut: bool) -> int:
            """
            Args:
                root (Optional[TreeNode]): 当前结点
                curPut (bool): 当前结点是否放置了摄像头
                prePut (bool): 父结点是否放置了摄像头

            Returns:
                int: 返回当前结点为根的子树,需要放置的摄像头数量
            """
            if root is None:
                return INF if curPut else 0  # !无结点但是当前结点放置了摄像头,返回INF

            res = INF
            for leftPut, rightPut in gen(curPut, prePut):
                res = min(
                    res,
                    dfs(root.left, leftPut, curPut) + dfs(root.right, rightPut, curPut) + curPut,
                )

            return res

        return min(dfs(root, True, False), dfs(root, False, False))
