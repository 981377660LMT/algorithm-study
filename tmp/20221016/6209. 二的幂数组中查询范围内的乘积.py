from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 你需要找到一个下标从 0 开始的数组 powers ，它包含 最少 数目的 2 的幂，且它们的和为 n 。
# powers 数组是 非递减 顺序的。根据前面描述，构造 powers 数组的方法是唯一的。
# 同时给你一个下标从 0 开始的二维整数数组 queries ，其中 queries[i] = [lefti, righti] ，其中 queries[i] 表示请你求出满足 lefti <= j <= righti 的所有 powers[j] 的乘积。
# 请你返回一个数组 answers ，长度与 queries 的长度相同，其中 answers[i]是第 i 个查询的答案。由于查询的结果可能非常大，请你将每个 answers[i] 都对 109 + 7 取余 。


class Solution:
    def productQueries(self, n: int, queries: List[List[int]]) -> List[int]:
        """前缀积"""
        pow2 = []
        for bit in range(32):
            if n & (1 << bit):
                pow2.append((1 << bit) % MOD)
        preMul = [1] + list(accumulate(pow2, lambda x, y: x * y % MOD))
        res = []
        for left, right in queries:
            res.append(preMul[right + 1] * pow(preMul[left], MOD - 2, MOD) % MOD)
        return res

    def productQueries2(self, n: int, queries: List[List[int]]) -> List[int]:
        """用前缀和来处理"""
        pow2 = []
        for bit in range(32):
            if n & (1 << bit):
                pow2.append(bit)
        preSum = [0] + list(accumulate(pow2))
        res = []
        for left, right in queries:
            sum_ = preSum[right + 1] - preSum[left]
            res.append(pow(2, sum_, MOD))
        return res


print(Solution().productQueries(n=15, queries=[[0, 1], [2, 2], [0, 3]]))
print(Solution().productQueries2(n=15, queries=[[0, 1], [2, 2], [0, 3]]))
