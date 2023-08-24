# 2 <= n <= 2e4
# 1 <= queries.length <= 20

# 第 j 个查询的答案是满足如下条件的点对 (a, b) 的数目：
# a < b
# cnt 是与 a 或者 b 相连的边的数目，且 cnt `严格`大于 queries[j] 。

# !请注意，图中可能会有 重复边 。
# !无向图本质不同的度数只有根号m个(有若干数的和为n,则不同数的个数最多为根号n个).

from itertools import accumulate
from typing import List
from collections import Counter, defaultdict


# https://leetcode.cn/problems/count-pairs-of-nodes/solutions/2400682/ji-bai-100cong-shuang-zhi-zhen-dao-zhong-yhze/


class Solution:
    def countPairs(self, n: int, edges: List[List[int]], queries: List[int]) -> List[int]:
        """O(n+m+q)
        !无向图本质不同的度数只有根号m个.因此可以二重循环枚举度数.
        """
        deg = [0] * n
        edgeCounter = defaultdict(int)  # 统计重复边
        for u, v in edges:
            u, v = u - 1, v - 1
            deg[u] += 1
            deg[v] += 1
            if u > v:
                u, v = v, u
            edgeCounter[(u, v)] += 1

        degFreq = Counter(deg)  # !最多只有根号m个不同的度数.因此可以二重循环枚举度数.
        pairDegSum = [0] * (max(deg) * 2 + 2)  # [0, 0, 0, 1, 1, 2, 1, 1, 0, 0]
        for deg1, freq1 in degFreq.items():
            for deg2, freq2 in degFreq.items():
                if deg1 == deg2:
                    pairDegSum[deg1 * 2] += freq1 * (freq1 - 1) // 2
                elif deg1 < deg2:
                    pairDegSum[deg1 + deg2] += freq1 * freq2

        # 减去重边的影响
        for (u, v), count in edgeCounter.items():
            sum_ = deg[u] + deg[v]
            pairDegSum[sum_] -= 1
            pairDegSum[sum_ - count] += 1

        sufSum = ([0] + list(accumulate(pairDegSum[::-1])))[::-1]  # pairDegSum的后缀和
        res = [0] * len(queries)
        for qi, threshold in enumerate(queries):
            res[qi] = sufSum[threshold + 1] if threshold + 1 < len(sufSum) else 0
        return res

    def countPairs2(self, n: int, edges: List[List[int]], queries: List[int]) -> List[int]:
        """!排序+双指针 O(nlogn+q(n+m))"""
        deg, res = [0] * n, [0] * len(queries)
        edgeCounter = Counter()  # 统计重复边
        for u, v in edges:
            u, v = u - 1, v - 1
            deg[u] += 1
            deg[v] += 1
            if u > v:
                u, v = v, u
            edgeCounter[(u, v)] += 1

        # 排序的两数之和
        sortedDeg = sorted(deg)
        for qi, threshold in enumerate(queries):
            cur = 0
            left, right = 0, n - 1
            while left < right:
                if sortedDeg[left] + sortedDeg[right] > threshold:
                    cur += right - left
                    right -= 1
                else:
                    left += 1

            for (u, v), count in edgeCounter.items():
                # !重复边导致的假阳性：度数重复算了count 次
                if threshold < deg[u] + deg[v] and threshold >= deg[u] + deg[v] - count:
                    cur -= 1

            res[qi] = cur

        return res


print(Solution().countPairs(n=4, edges=[[1, 2], [2, 4], [1, 3], [2, 3], [2, 1]], queries=[2, 3]))
# 输出：[6,5]
# 解释：每个点对中，与至少一个点相连的边的数目如上图所示。
