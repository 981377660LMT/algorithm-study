# edges[i] = [typei, ui, vi] 表示节点 ui 和 vi 之间存在类型为 typei 的双向边
# 类型 1：只能由 Alice 遍历。
# 类型 2：只能由 Bob 遍历。
# 类型 3：Alice 和 Bob 都可以遍历。
# 请你在保证图仍能够被 Alice和 Bob `完全遍历`的前提下，找出可以删除的最大边数。
# 如果从任何节点开始，Alice 和 Bob 都可以到达所有其他节点，则认为图是可以完全遍历的。


# 总结：
# 我们删除最多数目的边，这等价于保留最少数目的边
# 换句话说，我们可以从一个仅包含 n 个节点（而没有边）的无向图开始，
# 逐步添加边，使得满足上述的要求。

# 那么我们应该按照什么策略来添加边呢？ 公共边先加
# 在处理完了所有的「公共边」之后，我们需要处理他们各自的独占边
# 将当前的并查集复制一份，一份交给 Alice，一份交给 Bob。
# 随后 Alice 不断地向并查集中添加「Alice 独占边」，
# Bob 不断地向并查集中添加「Bob 独占边」。
# 在处理完了所有的独占边之后，如果这两个并查集都只包含一个连通分量，
# 那么就说明 Alice 和 Bob 都可以遍历整个无向图。

# 在使用并查集进行合并的过程中，
# 我们每遇到一次失败的合并操作（即需要合并的两个点属于同一个连通分量），
# 那么就说明当前这条边可以被删除，将答案增加 1。


# 总结:
# 先处理公共边

from typing import List


class Solution:
    def maxNumEdgesToRemove(self, n: int, edges: List[List[int]]) -> int:
        ufa, ufb = UF(n), UF(n)
        # 可以删除的边数(已经联通就可以删除了)
        res = 0

        # 节点编号改为从 0 开始
        for edge in edges:
            edge[1] -= 1
            edge[2] -= 1

        # 公共边
        for t, u, v in edges:
            if t == 3:
                if not ufa.union(u, v):
                    res += 1
                else:
                    ufb.union(u, v)

        # 独占边
        for t, u, v in edges:
            if t == 1:
                # Alice 独占边
                if not ufa.union(u, v):
                    res += 1
            elif t == 2:
                # Bob 独占边
                if not ufb.union(u, v):
                    res += 1

        # print(ufa.count, ufb.count)
        # 是否全部联通
        if ufa.count != 1 or ufb.count != 1:
            return -1

        return res


class UF:
    def __init__(self, n):
        self.parent = list(range(n))
        self.count = n

    def union(self, x, y):
        rx, ry = self.find(x), self.find(y)
        if rx == ry:
            return False
        low, high = sorted([rx, ry])
        self.parent[high] = low
        self.count -= 1
        return True

    def find(self, i):
        if i != self.parent[i]:
            self.parent[i] = self.find(self.parent[i])
        return self.parent[i]

    def isConnected(self, x, y):
        rx, ry = self.find(x), self.find(y)
        return rx == ry


print(
    Solution().maxNumEdgesToRemove(
        n=4, edges=[[3, 1, 2], [3, 2, 3], [1, 1, 3], [1, 2, 4], [1, 1, 2], [2, 3, 4]]
    )
)


# 输出：2
# 解释：如果删除 [1,1,2] 和 [1,1,3] 这两条边，
# Alice 和 Bob 仍然可以完全遍历这个图。
# 再删除任何其他的边都无法保证图可以完全遍历。所以可以删除的最大边数是 2 。
