"""树形dp


「力扣嘉年华」的中心广场放置了一个巨型的二叉树形状的装饰树。
每个节点上均有一盏灯和三个开关。
节点值为 0 表示灯处于「关闭」状态,节点值为 1 表示灯处于「开启」状态。

每个节点上的三个开关各自功能如下:
开关 1:切换`当前节点`的灯的状态；
开关 2:切换`以当前节点为根的子树`中,所有节点上的灯的状态;
开关 3:切换`当前节点及其左右子节点`（若存在的话） 上的灯的状态；
给定该装饰的初始状态 root,请返回最少需要操作多少次开关,可以关闭所有节点的灯

枚举当前结点每种操作次数,后序dfs返回每种状态的操作次数
0:全亮 1:全灭 2:当前灯亮,其余全灭 3:当前灯灭,其余全亮
"""

from typing import Optional


# Definition for a binary tree node.
class TreeNode:
    def __init__(self, x):
        self.val = x
        self.left = None
        self.right = None


INF = int(1e20)


class Solution:
    def closeLampInTree(self, root: "TreeNode") -> int:
        """记忆化dfs写法(对象也能哈希 会自动生成哈希值)"""

        def dfs(cur: Optional["TreeNode"], switch2: bool, switch3: bool) -> int:
            ...

        return dfs(root, False, False)

    def closeLampInTree2(self, root: "TreeNode") -> int:
        """树形dp写法(这道题不太推荐,因为状态太多了,用dfs方便一些)

        注意 左右节点一定是同时改变的 想要全灭的话状态必须一样
        """

        ...
