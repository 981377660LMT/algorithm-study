# 有依赖的背包问题-树上背包
# 有 N 个物品和一个容量是 V 的背包。
# 物品之间具有依赖关系，且依赖关系组成一棵树的形状。如果选择一个物品，则必须选择它的父节点。
# 每件物品的编号是 i，体积是 vi，价值是 wi，依赖的父节点编号是 pi。物品的下标范围是 1…N。
# 求解将哪些物品装入背包，可使物品总体积不超过背包容量，且总价值最大。
# 输出最大价值。

"""
用树形dp的思想来求解
dp(i, j)表示在i节点为根的子树中选择节点，总开销不超过j的所有选法中，最优选法的价值的最大和
在i子树上选择，根节点是必选的，剩下的容量需要分配给其子树，把每一个
子树看成一个分组，每个分组中选择这个分组需要的容量，这些子树的总容量加起来不能超过i子树的剩下容量
这一步就转换为一个分组背包问题，物品就是选择的分配给子树的容量，价值就是这个子树在该容量约束下的最大子节点价值和
分组背包完成后，就可以得到一种最佳的剩余容量分配方案让i的所有子树的价值和最大，进而dp(i, j)也就求出来了
"""
from collections import defaultdict
from functools import lru_cache

N, V = map(int, input().split())
adjMap = defaultdict(list)
nodes = [(0, 0)]  # 虚拟节点体积为0，价值为0
for i in range(1, N + 1):
    cost, score, parent = map(int, input().split())
    if parent == -1:
        parent = 0
    adjMap[parent].append(i)
    nodes.append((cost, score))


@lru_cache(None)
def dfs(root: int, cap: int) -> int:
    """根节点为root的子树, 并且体积不超过cap的约束条件下, 能够获得的最大价值"""
    # 根节点的体积大于V的时候, 连更节点都无法选中, 很明显, 价值为0
    if nodes[root][0] > cap:
        return 0

    # res: 根节点的价值
    rootValue = nodes[root][1]
    # 去除掉根节点后的体积
    cap -= nodes[root][0]

    # 分组背包 看每个组选多少个商品, 可以使得每个组的总体积不超过cap且价格最大
    dp = [0] * (cap + 1)
    for i in range(len(adjMap[root])):
        for j in range(cap, -1, -1):
            cost = nodes[adjMap[root][i]][0]
            for select in range(cost, j + 1):
                dp[j] = max(dp[j], dfs(adjMap[root][i], select) + dp[j - select])
    return rootValue + dp[cap]


print(dfs(0, V))  # 这里的0是虚拟节点, 它的体积为0, 价值为0
