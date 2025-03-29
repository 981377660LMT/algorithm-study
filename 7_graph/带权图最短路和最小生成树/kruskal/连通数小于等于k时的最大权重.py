# 给定一个无向图，该图包含n个节点和若干条边。每条边有以下四个属性：
# • from：边的起点节点
# • to：边的终点节点
# • value：边的权重
# • isRed：是否为红色边（如果 isRed = true，则可以选择删除此边）
#
# 你的任务是：
# 对于k在1-N中，求出当连通数小于等于k时的最大权重
#
# 输入
# • 一个整数 n（1sns1e5），表示图的节点数。
# • 一个整数m（1sms2e5），表示边的数量。
# • m条边，每条边包含四个值：
# • from, to, value, isRed
# • isRed 是0（不可删除）或1（可删除）
#
# 输出
# !n行，每行一个整数，第i行表示当图被划分成i个连通分量时，可以获得的最大权重和。
# 示例
# 示例1
# 输入：
# 4 4
# 1251
# 2331
# 3440
# 1421
# 输出：
# 5
# 8 （5+3）
# 10 (5+3+2)
# 10 (5+3+2)


class DSU:
    def __init__(self, n):
        self.parent = list(range(n))
        self.rank = [0] * n

    def find(self, x):
        if self.parent[x] != x:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x, y):
        rx, ry = self.find(x), self.find(y)
        if rx == ry:
            return False
        if self.rank[rx] < self.rank[ry]:
            rx, ry = ry, rx
        self.parent[ry] = rx
        if self.rank[rx] == self.rank[ry]:
            self.rank[rx] += 1
        return True


def main():
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w, is_red = map(int, input().split())
        edges.append((u, v, w, is_red))

    # 1. 用非红边 (is_red == 0) 构造固定骨架，计算固定连通分量数 c0
    dsu_non_red = DSU(n)
    for u, v, w, is_red in edges:
        if is_red == 0:
            dsu_non_red.union(u, v)

    comp_set = set()
    for i in range(n):
        comp_set.add(dsu_non_red.find(i))
    c0 = len(comp_set)

    # 为辅助图中固定连通分量重新编号（0-indexed）
    comp_index = {}
    idx = 0
    for r in comp_set:
        comp_index[r] = idx
        idx += 1

    # 计算所有红边总权值（无论是否连接同一分量，都可以删除）
    total_red = 0
    for u, v, w, is_red in edges:
        if is_red == 1:
            total_red += w

    # 2. 构造辅助图：只考虑红边连接不同固定连通分量的情况
    cand = []
    for u, v, w, is_red in edges:
        if is_red == 1:
            ru = dsu_non_red.find(u)
            rv = dsu_non_red.find(v)
            if ru != rv:
                cu = comp_index[ru]
                cv = comp_index[rv]
                cand.append((w, cu, cv))

    # 在辅助图上使用 Kruskal 算法求 MST
    cand.sort()  # 按权值升序排序
    dsu_comp = DSU(c0)
    T_MST = 0
    mst_edges = []  # 保存 MST 中红边的权值
    for w, u, v in cand:
        if dsu_comp.union(u, v):
            T_MST += w
            mst_edges.append(w)

    # 对 MST 中的边按权值降序排序，便于“剪枝”
    mst_edges.sort(reverse=True)
    msize = len(mst_edges)
    prefix = [0] * (msize + 1)
    for i in range(1, msize + 1):
        prefix[i] = prefix[i - 1] + mst_edges[i - 1]

    # 3. 输出答案
    # 对于 k = 1..n，若 k >= c0，则答案为 total_red，
    # 否则必须保留 (c0 - k) 条红边，保留红边的成本为 T_MST - prefix[k - 1]
    for k in range(1, n + 1):
        if k >= c0:
            ans = total_red
        else:
            cost = T_MST - prefix[k - 1]
            ans = total_red - cost
        print(ans)


if __name__ == "__main__":
    main()
