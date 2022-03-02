from collections import defaultdict
from typing import List

# 1 <= pairs.length <= 10^5
# 1 <= xi < yi <= 500
# https://leetcode-cn.com/problems/number-of-ways-to-reconstruct-a-tree/comments/1383349
class Solution:
    def checkWays(self, pairs: List[List[int]]) -> int:
        # 「所谓构筑一棵树，其实就是在连通图里，为每一个节点寻找父节点」。
        degree, nxt = defaultdict(int), defaultdict(set)
        for cur, next in pairs:
            degree[cur] += 1
            degree[next] += 1
            nxt[cur] |= {cur, next}
            nxt[next] |= {cur, next}

        pairs = [sorted(pair, key=lambda x: degree[x]) for pair in pairs]
        # 根是出现次数最多的
        if max(degree.values()) != len(degree) - 1 or not all(
            nxt[cur] <= nxt[next] for cur, next in pairs
        ):
            return 0

        # 多种：存在等价的结点
        return 2 if any(degree[u] == degree[v] for u, v in pairs) else 1


# 真没看懂
