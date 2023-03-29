# https://yukicoder.me/problems/no/1295
# https://yukicoder.me/submissions/837841

# 类似反转颜色的题:
# https://atcoder.jp/contests/arc097/tasks/arc097_d
# 给定一个n个点的树
# 所有的顶点分为"访问过的顶点"和"未访问过的顶点"
# 开始时,将棋子放在根节点i上
# 每个回合,有两种移动方式:
# !1.将棋子移动到相邻的"未访问过的顶点"中编号最小的顶点上
# !2.将棋子移动到相邻的"访问过的顶点"中编号最小的顶点上
# 对每个根节点i，问是否可以使得所有的顶点都被访问到.


# TODO
# 首先证明，存在方式当且仅当存在某一个点为根时，满足：

# 1. 除了一条最终不返回的路径以外，其余点均满足“父亲是自己的最小邻居”。这样可以遍历完对应的子树后，才有机会回到父节点再遍历其它兄弟。
# 2. 最终不返回的路径上的每个点，必须是父节点的“最大孩子”或“最小邻居”。如果是最大孩子，在遍历父节点的子节点时，它是最后一个进入的。
# 如果是最小邻居，在遍历父节点的子节点时，先进入它，然后立即退出，遍历完其它孩子后，再次进入它。
# 最大仅要求孩子中最大(重儿子)，最小要求包含父节点在内最小

# 然后解释状态：

# * DP[i]=0，表示该节点为根所有点均满足“父亲是自己的最小邻居”，也就是可以返回的路径。
# * DP[i]=2，表示该节点存在不返回而遍历的路径，是父节点的最大邻居或最小邻居。
# * DP[i]=3,  表示该节点可能存在不返回而遍历的路径，是父节点的第二大邻居。只有父节点的父节点是其最大邻居时，才存在合法的路径。
# * DP[i]>=4，表示彻底无解。

# 然后解释代码：
# d表示当前DP值，p表示父节点，e.to表示当前节点

# f_ee = return l+r;  表示如果有两个孩子无法返回，就无解。

# d>=4 return 4; 如果子树无解，那么无解。
# d==3 && p!=ma[e.to] return 4; 如果存在某个后代是e.to的第二大邻居，仅当p=ma[e.to]时，存在合法而不返回的路径。
# e.to == ma[p] || e.to=mi[p] return 2;  存在不返回而遍历的路径
# e.to == ma2[p] return 3;  可能存在不返回而遍历的路径，还需要看父节点确认（参见d==3那条）
# return 4;    子节点不能返回，但又不是父节点的最大孩子或最小邻居，直接无解。

from Rerooting import Rerooting


INF = int(4e18)

if __name__ == "__main__":

    E = int  # 状态

    OK = 0
    CANT_BACK = 2
    TO_MAX_2 = 3
    BAD = 4

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 + childRes2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        from_ = cur if direction == 0 else parent
        to = parent if direction == 0 else cur
        if fromRes >= BAD or (fromRes == TO_MAX_2 and from_ != max_[to]):
            return BAD
        if fromRes == OK and from_ == min_[to]:
            return OK
        if to == max_[from_] or to == min_[from_]:
            return CANT_BACK
        if to == max2[from_]:
            return TO_MAX_2
        return BAD

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))
    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    tree = R.adjList
    max_, max2, min_ = [0] * n, [0] * n, [0] * n  # 每个顶点邻接的最大值, 次大值(不存在为-1), 最小值
    for i in range(n):
        nexts = sorted(tree[i])
        max_[i] = nexts[-1]
        max2[i] = nexts[-2] if len(nexts) > 1 else -1
        min_[i] = nexts[0]

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)
    for state in dp:
        print("Yes" if state <= 2 else "No")
