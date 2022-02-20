# 用bfsOrder来判断是否有子树,dfsOrder来深度遍历

from typing import List

#
#
# @param k int整型 表示完全k叉树的叉数k
# @param a int整型一维数组 表示这棵完全k叉树的Dfs遍历序列的结点编号
# @return long长整型
#
# 将所有根结点与它的叶子结点异或,并将结果相加


class Solution:
    def tree6(self, k: int, a: List[int]) -> int:
        def dfs(root: int) -> None:
            if self.dfsIndex >= n:
                return
            parent = a[self.dfsIndex]
            for i in range(1, k + 1):
                child = root * k + i
                if child < n:
                    self.dfsIndex += 1
                    self.res += parent ^ a[self.dfsIndex]
                    dfs(child)
                else:
                    break

        n = len(a)
        self.dfsIndex = 0
        self.res = 0
        dfs(0)

        return self.res


# k叉树父节点:parent=(child-1)//k
# k叉树子节点:child=root*k+i
