# #  找到区间[L,R]恰好为K个B的幂次方之和的数的个数
# # 1≤X≤Y≤231−1,
# # 1≤K≤20,
# # 2≤B≤10
# from functools import lru_cache


from functools import lru_cache

from typing import List


left, right = map(int, input().split())
k = int(input())
radix = int(input())

# 预处理组合数 C(n,k)=C(n-1,k)+C(n-1,k-1)
comb = [[0] * 36 for _ in range(36)]
for i in range(36):
    comb[i][0] = 1
    for j in range(1, i + 1):
        comb[i][j] = comb[i - 1][j - 1] + comb[i - 1][j]


def cal(upper: int) -> int:
    @lru_cache(None)
    def dfs(pos: int, count: int, isLimit: bool) -> int:
        """当前在第pos位，取了radix进制上的count个1，isLimit表示是否贴合上界"""
        if pos == 0:
            return int(count == k)
        # 如果不贴合上界，即剩下i位可以从00…00~99…99随便填，用组合数即可直接计算方案数，从i位中选取j个1
        if not isLimit:
            return comb[pos][k - count]

        res = 0
        up = nums[pos - 1] if isLimit else radix - 1
        # 枚举该位填0/1
        for cur in range(min(up, 1) + 1):
            res += dfs(pos - 1, count + int(cur == 1), (isLimit and cur == up))
        return res

    nums = []
    while upper:
        div, mod = divmod(upper, radix)
        nums.append(mod)
        upper = div
    return dfs(len(nums), 0, True)


print(cal(right) - cal(left - 1))
