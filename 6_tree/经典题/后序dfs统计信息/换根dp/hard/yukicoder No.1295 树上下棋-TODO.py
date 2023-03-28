# https://yukicoder.me/problems/no/1295
# https://yukicoder.me/submissions/837841

# 给定一个n个点的树
# 所有的顶点分为"访问过的顶点"和"未访问过的顶点"
# 开始时,将棋子放在根节点i上
# 每个回合,有两种移动方式:
# !1.将棋子移动到相邻的"未访问过的顶点"中编号最小的顶点上
# !2.将棋子移动到相邻的"访问过的顶点"中编号最小的顶点上

# 对每个根节点i，问是否可以使得所有的顶点都被访问到.

# TODO 不对

from Rerooting import Rerooting


INF = int(4e18)

if __name__ == "__main__":

    E = int  # 状态

    OK = 0
    TO_MAX = 2
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
            return TO_MAX
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
