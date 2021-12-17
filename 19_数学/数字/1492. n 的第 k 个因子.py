from math import sqrt

# O(sqrt(n)) solution

# 两两配对，如果最后一对相等(完全平方数)，则pop一个
class Solution:
    def kthFactor(self, n: int, k: int) -> int:
        small, big = [], []
        for f in range(1, int(sqrt(n)) + 1):
            if n % f == 0:
                small += [f]
                big += [n // f]

        if small[-1] == big[-1]:
            big.pop()

        factors = small + big[::-1]
        return factors[k - 1] if k - 1 < len(factors) else -1


print(Solution().kthFactor(n=12, k=3))
# 输出：3
# 解释：因子列表包括 [1, 2, 3, 4, 6, 12]，第 3 个因子是 3 。
