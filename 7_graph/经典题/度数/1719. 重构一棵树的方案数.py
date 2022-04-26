from typing import List
from collections import defaultdict

# 1 <= pairs.length <= 10^5
# 1 <= xi < yi <= 500
# https://leetcode-cn.com/problems/number-of-ways-to-reconstruct-a-tree/comments/1383349

# 同1932. 合并多棵二叉搜索树一样
# 我们可以先考虑`构造好的树有什么性质`，然后再反推解题。

# 时间复杂度 O(N*M)
class Solution:
    def checkWays(self, pairs: List[List[int]]) -> int:
        # 每一对pair确定了一条链
        degree, adjMap = defaultdict(int), defaultdict(set)
        for cur, next in pairs:
            degree[cur] += 1
            degree[next] += 1
            adjMap[cur] |= {cur, next}
            adjMap[next] |= {cur, next}

        pairs = [sorted(pair, key=degree.get) for pair in pairs]  # point free 写法

        # 根的度数为n-1 越靠近根的结点度数越大 配对数也要越大
        if max(degree.values()) != len(degree) - 1 or not all(
            adjMap[pre] <= adjMap[cur] for pre, cur in pairs
        ):
            return 0

        # 多种：存在等价的结点
        return 2 if any(degree[u] == degree[v] for u, v in pairs) else 1


# 真没看懂
print(Solution().checkWays([[1, 2], [2, 3], [3, 4], [4, 5], [5, 6]]))
