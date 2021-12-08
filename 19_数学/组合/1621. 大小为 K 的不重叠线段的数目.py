from math import comb

# 请你找到 恰好 k 个不重叠 线段且每个线段至少覆盖两个点的方案数
# 这 k 个线段不需要全部覆盖全部 n 个点，且它们的端点 可以 重合。
class Solution:
    def numberOfSets(self, n: int, k: int) -> int:
        return comb(n + k - 1, 2 * k) % (10 ** 9 + 7)  # 如果线段端点可以重合
        return comb(n, 2 * k) % (10 ** 9 + 7)  # 如果线段端点不能重合


print(Solution().numberOfSets(4, 2))

