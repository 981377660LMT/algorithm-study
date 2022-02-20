#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
# @param k int整型 表示完全k叉树的叉数k
# @param a int整型一维数组 表示这棵完全k叉树的Bfs遍历序列的结点编号
# @return long长整型
#
from typing import List

# 将所有根结点与它的叶子结点异或,并将结果相加
# 遍历叶子结点就ok,从1开始一直到最后一个,根节点为i-1/k


class Solution:
    def tree2(self, k: int, a: List[int]) -> int:
        res = 0
        for leaf in range(1, len(a)):
            res += a[leaf] ^ a[(leaf - 1) // k]
        return res


# k叉树父节点:parent=(child-1)//k
# k叉树子节点:child=root*k+i
