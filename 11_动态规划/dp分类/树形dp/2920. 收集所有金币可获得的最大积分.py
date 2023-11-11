# https://leetcode.cn/problems/maximum-points-after-collecting-coins-from-all-nodes/


# 节点 0 处现有一棵由 n 个节点组成的无向树，节点编号从 0 到 n - 1 。给你一个长度为 n - 1 的二维 整数 数组 edges ，其中 edges[i] = [ai, bi] 表示在树上的节点 ai 和 bi 之间存在一条边。另给你一个下标从 0 开始、长度为 n 的数组 coins 和一个整数 k ，其中 coins[i] 表示节点 i 处的金币数量。

# 从根节点开始，你必须收集所有金币。要想收集节点上的金币，必须先收集该节点的祖先节点上的金币。

# 节点 i 上的金币可以用下述方法之一进行收集：

# 收集所有金币，得到共计 coins[i] - k 点积分。如果 coins[i] - k 是负数，你将会失去 abs(coins[i] - k) 点积分。
# 收集所有金币，得到共计 floor(coins[i] / 2) 点积分。如果采用这种方法，节点 i 子树中所有节点 j 的金币数 coins[j] 将会减少至 floor(coins[j] / 2) 。
# 返回收集 所有 树节点的金币之后可以获得的最大积分。

# !记忆化有问题是指：从所有结点开始DFS的情况，菊花图会卡


from functools import lru_cache
from typing import List


def max2(a, b):
    return a if a > b else b


def min2(a, b):
    return a if a < b else b


class Solution:
    def maximumPoints(self, edges: List[List[int]], coins: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(cur: int, pre: int, log: int) -> int:
            curCoin = coins[cur] >> log

            res1 = curCoin >> 1
            res2 = curCoin - k
            for next_ in adjList[cur]:
                if next_ == pre:
                    continue
                res1 += dfs(next_, cur, min2(log + 1, maxLog))
                res2 += dfs(next_, cur, log)
            return max2(res1, res2)

        n = len(coins)
        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)
        max_ = max(coins)
        maxLog = max_.bit_length()

        res = dfs(0, -1, 0)
        dfs.cache_clear()
        return res
