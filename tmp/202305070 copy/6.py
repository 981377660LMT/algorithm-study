from ast import Tuple, mod
from itertools import permutations
from typing import List, Sequence

# 探险家小扣终于来到了万灵之树前，挑战最后的谜题。
# 已知小扣拥有足够数量的链接节点和 n 颗幻境宝石，gem[i] 表示第 i 颗宝石的数值。现在小扣需要使用这些链接节点和宝石组合成一颗二叉树，其组装规则为：

# 链接节点将作为二叉树中的非叶子节点，且每个链接节点必须拥有 2 个子节点；
# 幻境宝石将作为二叉树中的叶子节点，所有的幻境宝石都必须被使用。
# 能量首先进入根节点，而后将按如下规则进行移动和记录：

# 若能量首次到达该节点时：
# 记录数字 1；
# 若该节点为叶节点，将额外记录该叶节点的数值；
# 若存在未到达的子节点，则选取未到达的一个子节点（优先选取左子节点）进入；
# 若无子节点或所有子节点均到达过，此时记录 9，并回到当前节点的父节点（若存在）。
# 如果最终记下的数依序连接成一个整数 num，满足

# nummod p=target，则视为解开谜题。
# 请问有多少种二叉树的组装方案，可以使得最终记录下的数字可以解开谜题

# 注意：


# 两棵结构不同的二叉树，作为不同的组装方案
# 两棵结构相同的二叉树且存在某个相同位置处的宝石编号不同，也作为不同的组装方案
# 可能存在数值相同的两颗宝石


# 多少种方案使得二叉树的前序遍历mod p=target
# 9!*4!*2!


class Solution:
    def treeOfInfiniteSouls(self, gem: List[int], p: int, target: int) -> int:
        def dfs(cur: int) -> None:
            path.append(1)
            if 0 <= cur < n:
                path.append(gem[cur])
                path.append(9)
                return
            if leftChild[cur] != -1:
                dfs(leftChild[cur])
            if rightChild[cur] != -1:
                dfs(rightChild[cur])
            path.append(9)

        def build(curs: Sequence[int], curId: int) -> None:
            if len(curs) <= 1:
                nonlocal res, path
                path = []
                dfs(curs[0])
                print(path)
                res += int("".join(map(str, path))) % p == target
                return
            for perm in permutations(curs):
                if len(curs) & 1:
                    nextLevel = []
                    for i in range(0, len(perm) - 1, 2):
                        leftChild[curId + i] = perm[i]
                        rightChild[curId + i] = perm[i + 1]
                        nextLevel.append(curId + i)
                    last = perm[-1]
                    nextLevel.append(last)
                    # 剩下一个
                    build(nextLevel, nextLevel[-2] + 1)
                    for i in range(0, len(perm) - 1, 2):
                        leftChild[curId + i] = -1
                        rightChild[curId + i] = -1
                else:
                    nexts = []
                    for i in range(0, len(perm), 2):
                        leftChild[curId + i] = perm[i]
                        rightChild[curId + i] = perm[i + 1]
                        nexts.append(curId + i)
                    build(nexts, nexts[-1] + 1)
                    # todo
                    for i in range(0, len(perm), 2):
                        leftChild[curId + i] = -1
                        rightChild[curId + i] = -1

        res = 0
        path = []
        n = len(gem)
        ids = list(range(len(gem)))
        leftChild, rightChild = [-1] * 100, [-1] * 100
        build(ids, len(gem))
        return res


# # gem = [2,3]
# # p = 100000007
# # target = 11391299
# print(Solution().treeOfInfiniteSouls([2, 3], 100000007, 11391299))
# # gem = [3,21,3]
# # p = 7
# # target = 5
# print(Solution().treeOfInfiniteSouls([3, 21, 3], 7, 5))
# [32,81,75,43]
# 14591
# 5395
print(Solution().treeOfInfiniteSouls([32, 81, 75, 43], 14591, 5395))
# [78,11,75,8]
# 83196847
# 1493875

print(Solution().treeOfInfiniteSouls([78, 11, 75, 8], 83196847, 1493875))
