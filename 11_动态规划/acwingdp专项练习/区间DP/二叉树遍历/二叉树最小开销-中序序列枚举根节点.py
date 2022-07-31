# 美团 二叉树中序遍历

# 小团有一个由N个节点组成的二叉树，每个节点有一个权值。
# 定义二叉树每条边的开销为其两端节点权值的乘积，
# 二叉树的总开销即每条边的开销之和。小团按照二叉树的中序遍历依次记录下每个节点的权值，
# 即他记录下了N个数，第i个数表示位于中序遍历第i个位置的节点的权值。
# 之后由于某种原因，小团遗忘了二叉树的具体结构。在所有可能的二叉树中，
# 总开销最小的二叉树被称为最优二叉树。
# 现在，小团请小美求出最优二叉树的总开销。
# n<=300
# 复杂度(n ^ 3)

from functools import lru_cache


n = int(input())
nums = [int(i) for i in input().split()]
INF = int(1e20)

# '左中右'
# 中序遍历的特点：选取其中一个节点，其左边的节点都是其左子树上的节点，其右边的节点都是其右子树上的节点。
# 在中序序列上枚举了每一个点作为根节点


@lru_cache(None)
def dfs(left: int, right: int, root: int) -> int:
    if left >= right:
        return 0

    res = INF
    for i in range(left, right):
        leftRes = dfs(left, i, nums[i])
        rightRes = dfs(i + 1, right, nums[i])
        res = min(res, leftRes + rightRes + nums[i] * root)
    return res


print(dfs(0, n, 0))
