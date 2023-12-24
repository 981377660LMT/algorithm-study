from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一棵 n 个节点的 无向 树，节点编号为 0 到 n - 1 ，树的根节点在节点 0 处。同时给你一个长度为 n - 1 的二维整数数组 edges ，其中 edges[i] = [ai, bi] 表示树中节点 ai 和 bi 之间有一条边。

# 给你一个长度为 n 下标从 0 开始的整数数组 cost ，其中 cost[i] 是第 i 个节点的 开销 。

# 你需要在树中每个节点都放置金币，在节点 i 处的金币数目计算方法如下：


# 如果节点 i 对应的子树中的节点数目小于 3 ，那么放 1 个金币。
# 否则，计算节点 i 对应的子树内 3 个不同节点的开销乘积的 最大值 ，并在节点 i 处放置对应数目的金币。如果最大乘积是 负数 ，那么放置 0 个金币。
# 请你返回一个长度为 n 的数组 coin ，coin[i]是节点 i 处的金币数目。
class Solution:
    def placedCoins(self, edges: List[List[int]], cost: List[int]) -> List[int]:
        ...
