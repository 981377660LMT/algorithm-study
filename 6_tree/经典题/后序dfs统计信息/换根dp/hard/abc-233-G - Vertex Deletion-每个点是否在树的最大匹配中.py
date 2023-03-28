# abc-233-G - Vertex Deletion-每个点是否在树的最大匹配中
# https://atcoder.jp/contests/abc223/tasks/abc223_g
# 给定一棵树
# 对每个结点i为根,删除根连接的所有边后,
# !使得剩下的树的最大匹配和原树最大匹配相等
# 求这样的根的个数

# !解:即不参与二分图的最大匹配
# https://yukicoder.me/problems/2085
# 二分图博弈
# Alice和Bob在树上博弈
# 先手放一个棋子,后手在相邻的结点放一个棋子
# 交替放棋子,直到不能放棋子的时候,输
# !问先手是否必胜 => 如果起点不在二分图的最大匹配中,先手必胜


from Rerooting import Rerooting


if __name__ == "__main__":

    E = int  # 当前节点是否构成子树的最大匹配, 0: 不参与, 1: 参与

    def e(root: int) -> E:
        return 0

    def op(childRes1: E, childRes2: E) -> E:
        return childRes1 | childRes2

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        return fromRes ^ 1  # 孩子参与匹配则父亲不参与, 反之成立

    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    R = Rerooting(n)
    for u, v in edges:
        R.addEdge(u, v)

    dp = R.rerooting(e=e, op=op, composition=composition, root=0)

    print(dp.count(0))  # 不在最大匹配中的点的个数
