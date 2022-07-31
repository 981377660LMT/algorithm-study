# 设一个 n 个节点的二叉树 tree 的中序遍历为（1,2,3,…,n），其中数字 1,2,3,…,n 为节点编号。

# 任一棵子树 subtree（也包含 tree 本身）的加分计算方法如下：
# subtree的左子树的加分 × subtree的右子树的加分 ＋ subtree的根的分数

# 若某个子树为空，规定其加分为 1。
# 叶子的加分就是叶节点本身的分数，不考虑它的空子树。
# 试求一棵符合中序遍历为（1,2,3,…,n）且加分最高的二叉树 tree。
# 要求输出：
# （1）tree的最高加分
# （2）tree的前序遍历


# '左中右'
# 中序遍历的特点：选取其中一个节点，其左边的节点都是其左子树上的节点，其右边的节点都是其右子树上的节点。
# 二叉树节点 向下投影，映射成的数组序列就是 中序遍历序列
# 在中序序列上枚举了每一个点作为根节点


from functools import lru_cache
from typing import Tuple

INF = int(1e20)
n = int(input())
nums = list(map(int, input().split()))


@lru_cache(None)
def dfs1(left: int, right: int) -> Tuple[int, int]:
    """中序序列[left:right+1]的二叉树的价值的最大值以及最大值对应的树的根节点"""
    """ 为了让前序序列的字典序最小，价值一样大情况下，根节点选更靠前的"""
    if left == right:
        return nums[left], left

    res = -INF
    root = -1
    for i in range(left, right + 1):  # !枚举根节点
        # 特殊边界情况
        if i == left:
            curVal = dfs1(left + 1, right)[0] + nums[i]
        elif i == right:
            curVal = dfs1(left, right - 1)[0] + nums[i]
        else:
            curVal = dfs1(left, i - 1)[0] * dfs1(i + 1, right)[0] + nums[i]

        if curVal > res:
            res, root = curVal, i

    return res, root


print(dfs1(0, n - 1)[0])

# 先序遍历打印二叉树， 因为最优划分已经全部算出来了，所以现在一个区间的最优划分方案就对应一个子树的根
res = []


def dfs2(left: int, right: int) -> None:
    if left > right:
        return
    if left == right:
        # 编号从1开始
        res.append(str(left + 1))
        return

    root = dfs1(left, right)[1]
    res.append(str(root + 1))
    dfs2(left, root - 1)
    dfs2(root + 1, right)


dfs2(0, n - 1)
print(" ".join(res))
